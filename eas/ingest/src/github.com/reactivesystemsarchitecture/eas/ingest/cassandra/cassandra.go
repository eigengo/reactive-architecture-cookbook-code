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
	var expectedTables = []string{"session_sensor_values", "session_labels", "session_sensor_data_sensors"}

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
	if err := ptypes.UnmarshalAny(envelope.Payload, &session); err != nil {
		return err
	}
	if session_id, err := uuid.Parse(session.SessionId); err != nil {
		return err
	} else {
		sliceSize := 16384
		sliceCount := len(session.SensorData.Values) / sliceSize
		for i := 0; i < sliceCount; i++ {
			sliceStart := i * sliceSize
			sliceEnd := sliceStart + sliceSize
			if sliceEnd > len(session.SensorData.Values) {
				sliceEnd = len(session.SensorData.Values)
			}
			slice := session.SensorData.Values[sliceStart:sliceEnd]
			if err := c.session.Query("insert into session_sensor_values(session_id, sequence, sensor_values) values (?, ?, ?)",
				session_id.String(), i, slice).Exec(); err != nil {
				return err
			}

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
	clusterConfig.Consistency = gocql.Quorum
	clusterConfig.Timeout = 30 * time.Second

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
