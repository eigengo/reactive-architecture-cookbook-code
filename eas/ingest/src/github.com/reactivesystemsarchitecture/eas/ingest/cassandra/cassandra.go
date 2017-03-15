package cassandra

import (
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	"github.com/golang/protobuf/ptypes"
	ps "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/reactivesystemsarchitecture/eas/ingest"
	"errors"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type sessionEnvelopeHandler struct {
	session *gocql.Session
}

func validateKeyspace(session *gocql.Session) error {
	if md, err := session.KeyspaceMetadata("eas"); err == nil {
		if _, ok := md.Tables["ingested_sessions"]; ok {
			// TODO: Check schema version
			return nil
		}

		return errors.New("Missing or bad schema for 'ingested_sessions'")
	} else {
		return err
	}
}

func (c *sessionEnvelopeHandler) Handle(envelope *p.Envelope) error {
	var session ps.Session
	if err := ptypes.UnmarshalAny(envelope.Payload, &session); err != nil {
		return err
	}
	if uuid, err := uuid.Parse(session.SessionId); err != nil {
		return err
	} else {
		return c.session.Query("insert into ingested_sessions(id) values (?)", uuid.String()).Exec()
	}
}

func NewPersistSessionEnvelopeHandler(hosts ...string) (ingest.EnvelopeHandler, error) {
	if len(hosts) == 0 {
		return nil, errors.New("No hosts supplied")
	}
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Keyspace = "eas"
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
