package cassandra

import (
	"testing"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	sp "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

func TestNewPersistSessionEnvelopeHandler(t *testing.T) {
	var e p.Envelope
	var s sp.Session
	ms, _ := ptypes.MarshalAny(&s)
	e.Payload = ms

	h, err := NewPersistSessionEnvelopeHandler("localhost:9092")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	if h.Handle(&e) != nil {
		t.Fail()
	}
}
