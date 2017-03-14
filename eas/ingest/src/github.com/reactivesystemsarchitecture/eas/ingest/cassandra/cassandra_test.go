package cassandra

import (
	"testing"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	sp "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/golang/protobuf/ptypes"
)

func TestNewPersistSessionEnvelopeHandler(t *testing.T) {
	var e p.Envelope
	var s sp.Session
	ms, _ := ptypes.MarshalAny(&s)
	e.Payload = ms

	h, err := NewPersistSessionEnvelopeHandler("localhost:9092")
	if err != nil {
		t.Fail()
	}
	if h.Handle(&e) != nil {
		t.Fail()
	}
}
