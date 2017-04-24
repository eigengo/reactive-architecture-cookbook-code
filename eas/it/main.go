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
)

type acceptLabel func (string) bool

func acceptLabelAny() acceptLabel {
	return func(x string) bool {
		return true
	}
}

type sensorData struct {
	label string
	data *v1m0.SensorData
}

var sensorsRegexp *regexp.Regexp = regexp.MustCompile(`\W+((\w+)->\[([^]]+)])+`)
func readSensors(r *bufio.Scanner) (sensors []*v1m0.Sensor, err error) {
	if r.Scan() {
		line := r.Text()
		for _, groups := range sensorsRegexp.FindAllStringSubmatch(line, 0) {
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

func readDataFilesIn(dirname string, acceptLabel acceptLabel) (sd []sensorData, err error) {
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
				label := path.Base(dirname)
				if !acceptLabel(label) { continue }

				if f, err := os.Open(path.Join(dirname, entry.Name())); err == nil {
					fileScanner := bufio.NewScanner(f)
					sensors, serr := readSensors(fileScanner)
					values, verr := readSensorValues(fileScanner)
					if serr != nil { return nil, serr }
					if verr != nil { return nil, verr }

					sd = append(sd, sensorData{
						data: &v1m0.SensorData{
							Values:  values,
							Sensors: sensors,
						},
						label: label,
					})
				}
			}

		}

		return sd, nil
	}
}

func newSession(dirname string, acceptLabel acceptLabel) (*v1m0.Session, error) {


	return nil, errors.New("f")
}

//func (s *session)toRequest(url string) http.Request {
//	// http.NewRequest("POST", url, body)
//}



func main() {
	var dataDir string
	flag.StringVar(&dataDir, "dir", "../data", "The directory containing the labelled data")

	if sds, err := readDataFilesIn("/Users/janmachacek/OReilly/reactive-architecture-cookbook-code/eas/it/data/labelled", acceptLabelAny()); err == nil {
		for _, sd := range sds {
			fmt.Println(sd.label)
			fmt.Println(sd.data.Values)
		}
	} else {
		fmt.Println(err)
	}

	line := "Wrist->[Acceleration]"
	fmt.Println(sensorsRegexp.FindString(line))
	fmt.Println(sensorsRegexp.FindAllString(line, 0))

}
