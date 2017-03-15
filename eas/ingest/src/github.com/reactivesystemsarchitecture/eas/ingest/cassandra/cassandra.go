package cassandra

import (
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	"github.com/golang/protobuf/ptypes"
	ps "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/reactivesystemsarchitecture/eas/ingest"
	"log"
	"errors"

	"github.com/gocql/gocql"
)

type sessionEnvelopeHandler struct {
	session *gocql.Session
}

func (c *sessionEnvelopeHandler) Handle(envelope *p.Envelope) error {
	var session ps.Session
	if err := ptypes.UnmarshalAny(envelope.Payload, &session); err != nil {
		return err
	}


	log.Println("Persisting", session)

	return nil
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
		return &sessionEnvelopeHandler{
			session: session,
		}, nil
	}
}
