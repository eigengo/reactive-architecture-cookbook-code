package ingest

import (
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	"github.com/golang/protobuf/ptypes"
	ps "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"log"
	"errors"
)

type CassandraSessionEnvelopeHandler struct {

}

func (c *CassandraSessionEnvelopeHandler) Handle(envelope *p.Envelope) error {
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

func (c *CassandraSessionEnvelopeHandler) Validate(envelope *p.Envelope) error {
	var session ps.Session
	if envelope.Payload == nil {
		return errors.New("envelope.Payload == nil")
	}
	if !ptypes.Is(envelope.Payload, &session) {
		return errors.New("Payload is not Session")
	}

	return nil
}

func NewCassandraSessionEnvelopeHandler() EnvelopeHandler {
	return &CassandraSessionEnvelopeHandler{}
}
