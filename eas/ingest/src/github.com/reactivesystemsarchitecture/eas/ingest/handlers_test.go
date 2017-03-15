package ingest

import (
	"testing"
	"github.com/reactivesystemsarchitecture/eas/protocol"
	"errors"
	"github.com/stretchr/testify/assert"
)

var f = EnvelopeHandlerFunc(func(_ *protocol.Envelope) error {
	return errors.New("bab")
})

func TestEnvelopeHandlerFunc_Handle(t *testing.T) {
	var e protocol.Envelope
	assert.EqualError(t, f.Handle(&e), "bab")
}
