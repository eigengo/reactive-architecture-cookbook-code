package ingest

import (
	"testing"
	"github.com/stretchr/testify/assert"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	sp "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/proto"
)

func assertFailWithCode(t *testing.T, payload proto.Message, code int) {
	ee := &ValidationError{code: code}
	var e p.Envelope
	if payload != nil {
		a, _ := ptypes.MarshalAny(payload)
		e.Payload = a
	}

	if err, ok := SessionEnvelopeValidator.Validate(&e).(*ValidationError); ok {
		assert.Equal(t, err, ee)
	} else {
		assert.Fail(t, "Bad error")
	}
}

func TestSessionEnvelopeValidator_NoSession(t *testing.T) {
	assertFailWithCode(t, nil, ValidationErrorNilPayload)
}

func TestSessionEnvelopeValidator_NotSession(t *testing.T) {
	var l sp.Label
	assertFailWithCode(t, &l, ValidationErrorPayloadNotSession)
}

