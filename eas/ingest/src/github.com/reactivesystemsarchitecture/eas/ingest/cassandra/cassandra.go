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

type cassandraEnvelopeHandler struct {
	clusterConfig *gocql.ClusterConfig
}

func (c *cassandraEnvelopeHandler) Handle(envelope *p.Envelope) error {

	var session ps.Session
	if err := c.Validate(envelope); err != nil {
		return err
	}
	if err := ptypes.UnmarshalAny(envelope.Payload, &session); err != nil {
		return err
	}

	log.Println("Persisting", session)

	return nil
}

func (c *cassandraEnvelopeHandler) Validate(envelope *p.Envelope) error {
	var session ps.Session
	if envelope.Payload == nil {
		return errors.New("envelope.Payload == nil")
	}
	if !ptypes.Is(envelope.Payload, &session) {
		return errors.New("Payload is not Session")
	}

	return nil
}

func NewCluster(hosts ...string) (ingest.EnvelopeHandler, error) {
	if len(hosts) == 0 {
		return nil, errors.New("No hosts supplied")
	}
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Keyspace = "eas"
	clusterConfig.Consistency = gocql.Quorum

	return &cassandraEnvelopeHandler{
		clusterConfig: clusterConfig,
	}, nil
}
