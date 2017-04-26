package main

import (
	"flag"
	"io/ioutil"
	"github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"os"
	"bufio"
	"regexp"
	"strings"
	"fmt"
	"strconv"
	"path"
	"errors"
	"time"
	"math/rand"
	"reflect"
	"github.com/google/uuid"
	"math"
	"log"
	"net/http"
	"github.com/golang/protobuf/proto"
	"bytes"
)

// An alias for function that decides whether to accept a label given its name
type acceptLabel func(name string) bool

// Returns acceptLabel function that accepts all labels
func acceptLabelAny() acceptLabel {
	return func(_ string) bool {
		return true
	}
}

// A small block of sensor data that can be accumulated to form a Session
type labelledSensorData struct {
	label string
	data  *v1m0.SensorData
}

// Returns a slice of the values matching the sensors
func (d *labelledSensorData) sensors(sensors []*v1m0.Sensor) (values []float32, err error) {
	// We could make this more flexible (d.data.Sensors to be a _subset_ of sensors),
	// but this is enough for now.
	if !reflect.DeepEqual(d.data.Sensors, sensors) {
		return nil, errors.New("Mismatched sensors")
	}

	return d.data.Values, nil
}

// Returns the duration of this block of data
func (d *labelledSensorData) duration() time.Duration {
	return sensorDataDuration(d.data.Values, d.data.Sensors)
}

type sessionBuilder struct {
	// the required sensors
	sensors []*v1m0.Sensor
	// the amount of rest / exercise.
	restFactor float64

	// -- working set --

	// the accumulated labels
	labels []*v1m0.Label
	// the accumulated data
	values []float32
}

// Returns the new session builder
func newSessionBuilder(sensors []*v1m0.Sensor, restFactor float64) *sessionBuilder {
	return &sessionBuilder{
		sensors:    sensors,
		restFactor: restFactor,
	}
}

func sensorDataValuesWidth(sensors []*v1m0.Sensor) (result int) {
	for _, s := range sensors {
		for _, t := range s.DataTypes {
			switch t {
			case v1m0.SensorDataType_Acceleration:
				result += 3
			case v1m0.SensorDataType_Rotation:
				result += 3
			case v1m0.SensorDataType_HeartRate:
				result += 1
			}
		}
	}

	return result
}

// Computes the duration of the values from the sensors
func sensorDataDuration(values []float32, sensors []*v1m0.Sensor) time.Duration {
	w := sensorDataValuesWidth(sensors)
	rows := len(values) / w
	return time.Duration(rows/50.0) * time.Second
}

// Returns empty values for the given duration of data from sensors
func emptyValues(sensors []*v1m0.Sensor, duration time.Duration) []float32 {
	w := sensorDataValuesWidth(sensors)
	return make([]float32, int(float64(w)*duration.Seconds()*50.0))
}

// Appends the data in `data` together with the rest (empty values) as specified
// by the restFactor.
func (b *sessionBuilder) appendData(data *labelledSensorData) error {
	if exerciseValues, err := data.sensors(b.sensors); err != nil {
		return err
	} else {
		ed := data.duration()
		rd := time.Duration((b.restFactor*(1+rand.Float64()/10))*ed.Seconds()) * time.Second
		d := sensorDataDuration(b.values, b.sensors)

		log.Printf("Appending data of duration %.fs\n", ed.Seconds())
		log.Printf("Appending rest of duration %.fs\n", rd.Seconds())

		restValues := emptyValues(b.sensors, rd)
		label := v1m0.Label{
			StartTime: d.Seconds(),
			Duration:  ed.Seconds(),
			Label:     data.label,
		}

		b.values = append(b.values, exerciseValues...)
		b.values = append(b.values, restValues...)
		b.labels = append(b.labels, &label)

		return nil
	}
}

// build the session
func (b *sessionBuilder) build() *v1m0.Session {
	sensorData := v1m0.SensorData{
		Values:  b.values,
		Sensors: b.sensors,
	}
	return &v1m0.Session{
		SessionId:       uuid.New().String(),
		UserLabels:      b.labels,
		AutomaticLabels: b.labels,
		SensorData:      &sensorData,
	}
}

var sensorsRegexp *regexp.Regexp = regexp.MustCompile("\\W+((\\w+)->\\[([^]]+)])+")

func readSensors(r *bufio.Scanner) (sensors []*v1m0.Sensor, err error) {
	if r.Scan() {
		line := r.Text()
		matches := sensorsRegexp.FindAllStringSubmatch(line, math.MaxInt32)
		if len(matches) == 0 {
			return nil, fmt.Errorf("'%s' is not a sensor definition", line)
		}
		for _, groups := range matches {
			var s v1m0.Sensor

			if l, ok := v1m0.SensorLocation_value[groups[2]]; ok {
				s.Location = v1m0.SensorLocation(l)
			} else {
				return nil, fmt.Errorf("Bad location %s", groups[2])
			}
			for _, dataType := range strings.Split(groups[3], ",") {
				if dt, ok := v1m0.SensorDataType_value[dataType]; ok {
					s.DataTypes = append(s.DataTypes, v1m0.SensorDataType(dt))
				} else {
					return nil, fmt.Errorf("Bad data type %s", dataType)
				}
			}

			sensors = append(sensors, &s)
		}
		return sensors, nil
	} else {
		return nil, errors.New("Could not scan")
	}

}

func readSensorValues(r *bufio.Scanner) (values []float32, err error) {
	for r.Scan() {
		line := r.Text()
		for _, value := range strings.Split(line, ",") {
			if f, err := strconv.ParseFloat(value, 32); err == nil {
				values = append(values, float32(f))
			} else {
				return nil, err
			}
		}
	}
	return values, nil
}

func readDataFilesIn(dirname string, acceptLabel acceptLabel) (sd []labelledSensorData, err error) {
	if entries, err := ioutil.ReadDir(dirname); err != nil {
		return nil, err
	} else {
		for _, entry := range entries {
			if entry.IsDir() && entry.Name()[0] != '.' {
				// recurse into directory
				if sr, err := readDataFilesIn(path.Join(dirname, entry.Name()), acceptLabel); err == nil {
					sd = append(sd, sr...)
				} else {
					return nil, err
				}
			} else if path.Ext(entry.Name()) == ".csv" {
				log.Printf("Reading %s\n", entry.Name())
				label := path.Base(dirname)
				if !acceptLabel(label) {
					log.Printf("Skipping label %s\n", label)
					continue
				}

				if f, err := os.Open(path.Join(dirname, entry.Name())); err == nil {
					fileScanner := bufio.NewScanner(f)
					sensors, serr := readSensors(fileScanner)
					values, verr := readSensorValues(fileScanner)
					if serr != nil {
						return nil, serr
					}
					if verr != nil {
						return nil, verr
					}

					sd = append(sd, labelledSensorData{
						data: &v1m0.SensorData{
							Values:  values,
							Sensors: sensors,
						},
						label: label,
					})

					log.Printf("Read label '%s' for %s", label, sensors)
				}
			}

		}

		return sd, nil
	}
}

func newSession(dirname string, acceptLabel acceptLabel) (*v1m0.Session, error) {
	sensors := []*v1m0.Sensor{
		{Location: v1m0.SensorLocation_Wrist, DataTypes: []v1m0.SensorDataType{v1m0.SensorDataType_Acceleration}},
	}
	builder := newSessionBuilder(sensors, 2.0)
	if sds, err := readDataFilesIn(dirname, acceptLabel); err == nil {
		for _, sd := range sds {
			if err := builder.appendData(&sd); err != nil {
				return nil, err
			}
		}
	} else {
		return nil, err
	}

	return builder.build(), nil
}

func newRequest(session *v1m0.Session, url string) *http.Request {
	body, _ := proto.Marshal(session)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Transfer-Encoding", "octet-stream")
	req.Header.Set("Content-Type", "application/x-protobuf")
	return req
}

func main() {
	var dataDir string
	flag.StringVar(&dataDir, "dir", "../data", "The directory containing the labelled data")

	if session, err := newSession("/Users/janmachacek/OReilly/reactive-architecture-cookbook-code/eas/it/data/labelled", acceptLabelAny()); err == nil {
		d := sensorDataDuration(session.SensorData.Values, session.SensorData.Sensors)
		rq := newRequest(session, "http://localhost:8080")
		log.Printf("Session %s (%s labels, duration %.fs)", session.SessionId, session.UserLabels, d.Seconds())
		log.Printf("Request %s", rq)
	} else {
		log.Fatalf("%s", err)
	}
}
