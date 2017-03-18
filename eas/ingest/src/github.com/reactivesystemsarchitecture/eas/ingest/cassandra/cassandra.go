package cassandra

import (
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	"github.com/golang/protobuf/ptypes"
	ps "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/reactivesystemsarchitecture/eas/ingest"
	"errors"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"fmt"
	"time"
)

type sessionEnvelopeHandler struct {
	session *gocql.Session
}

func validateKeyspace(session *gocql.Session) error {
	var expectedTables = []string{"session_sensor_values", "session_labels", "session_sensors"}

	if md, err := session.KeyspaceMetadata("eas"); err == nil {
		for _, table := range expectedTables {
			if _, ok := md.Tables[table]; !ok {
				// TODO: Check schema version
				return fmt.Errorf("Missing or bad schema for '%s'", table)
			}

		}

		return nil
	} else {
		return err
	}
}

func (c *sessionEnvelopeHandler) Handle(envelope *p.Envelope) error {
	var session ps.Session
	var err error = nil
	if err := ptypes.UnmarshalAny(envelope.Payload, &session); err != nil {
		return err
	}
	if _, err := uuid.Parse(session.SessionId); err != nil {
		return err
	}

	sliceSize := 32768
	sliceCount := len(session.SensorData.Values) / sliceSize
	insertSensorValues := c.session.Query("insert into session_sensor_values(session_id, sequence, sensor_values, valid) values (?, ?, ?, false)")
	insertSensorValues.RetryPolicy(&gocql.SimpleRetryPolicy{NumRetries: 3})
	insertSensorValues.Consistency(gocql.One)

	insertLabels := c.session.Query("insert into session_labels(session_id, sequence, start_time, duration, label, automatic, valid) values (?, ?, ?, ?, ?, ?, false)")
	insertLabels.RetryPolicy(&gocql.SimpleRetryPolicy{NumRetries: 3})
	insertLabels.Consistency(gocql.One)

	insertSensors := c.session.Query("insert into session_sensors(session_id, sequence, location, data_types, valid) values (?, ?, ?, ?, false)")
	insertSensors.RetryPolicy(&gocql.SimpleRetryPolicy{NumRetries: 3})
	insertSensors.Consistency(gocql.One)

	for i := 0; i < sliceCount; i++ {
		sliceStart := i * sliceSize
		sliceEnd := sliceStart + sliceSize
		if sliceEnd > len(session.SensorData.Values) {
			sliceEnd = len(session.SensorData.Values)
		}
		slice := session.SensorData.Values[sliceStart:sliceEnd]

		if err = insertSensorValues.Bind(session.SessionId, i, slice).Exec(); err != nil {
			return err
		}
	}

	if err = doInsertLabels(session.SessionId, session.AutomaticLabels, true, insertLabels); err != nil {
		return err
	}
	if err = doInsertLabels(session.SessionId, session.UserLabels, false, insertLabels); err != nil {
		return err
	}
	for i, sensor := range session.SensorData.Sensors {
		if err = insertSensors.Bind(session.SessionId, i, sensor.Location, sensor.DataTypes).Exec(); err != nil {
			return err
		}
	}

	// now set valid = true
	if err = c.session.Query("update session_sensor_values set valid = true where session_id = ? and sequence in ?",
		session.SessionId, makeRangeTo(sliceCount)).Exec(); err != nil {
		return err
	}
	if err = c.session.Query("update session_labels set valid = true where session_id = ? and automatic = true and sequence in ?",
		session.SessionId, makeRangeTo(len(session.AutomaticLabels))).Exec(); err != nil {
		return err
	}
	if err = c.session.Query("update session_labels set valid = true where session_id = ? and automatic = false and sequence in ?",
		session.SessionId, makeRangeTo(len(session.UserLabels))).Exec(); err != nil {
		return err
	}
	if err = c.session.Query("update session_sensors set valid = true where session_id = ? and sequence in ?",
		session.SessionId, makeRangeTo(len(session.SensorData.Sensors))).Exec(); err != nil {
		return err
	}

	return nil

}

func makeRangeTo(max int) []int {
	result := make([]int, max, max)
	for i := 0; i < max; i++ {
		result[i] = i
	}
	return result
}

func doInsertLabels(sessionId string, labels []*ps.Label, automatic bool, insertLabels *gocql.Query) error {
	for i, label := range labels {
		err := insertLabels.Bind(sessionId, i, label.StartTime, label.Duration, label.Label, automatic).Exec()

		if err != nil {
			break
		}
	}
	return nil
}

func NewPersistSessionEnvelopeHandler(hosts ...string) (ingest.EnvelopeHandler, error) {
	if len(hosts) == 0 {
		return nil, errors.New("No hosts supplied")
	}
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Keyspace = "eas"
	clusterConfig.Timeout = 3 * time.Second
	clusterConfig.Compressor = gocql.SnappyCompressor{}
	clusterConfig.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	clusterConfig.Consistency = gocql.Quorum

	if session, err := clusterConfig.CreateSession(); err != nil {
		return nil, err
	} else {
		if err := validateKeyspace(session); err == nil {
			return &sessionEnvelopeHandler{
				session: session,
			}, nil
		} else {
			return nil, err
		}
	}
}
