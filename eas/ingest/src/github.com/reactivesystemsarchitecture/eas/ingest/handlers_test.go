package ingest

import (
	"testing"
	"github.com/reactivesystemsarchitecture/eas/protocol"
	"errors"
	"github.com/stretchr/testify/assert"
)

func h(err error) EnvelopeHandlerFunc {
	return EnvelopeHandlerFunc(func(_ *protocol.Envelope) error {
		return err
	})
}

func v(err error) EnvelopeValidatorFunc {
	return EnvelopeValidatorFunc(func(_ *protocol.Envelope) error {
		return err
	})
}

func TestEnvelopeHandlerFunc_Handle(t *testing.T) {
	var e protocol.Envelope
	err := errors.New("kushadfkjasdf")
	assert.EqualError(t, h(err).Handle(&e), err.Error())
}

func TestEnvelopeValidatorFunc_Validate(t *testing.T) {
	var e protocol.Envelope
	err := errors.New("nanananananananana Batman!")
	assert.EqualError(t, v(err).Validate(&e), err.Error())
}

func TestNewEnvelopeProcessor(t *testing.T) {
	var e protocol.Envelope
	p := NewEnvelopeProcessor(v(nil), h(nil))
	assert.Nil(t, p.Handle(&e))
	assert.Nil(t, p.Validate(&e))
}