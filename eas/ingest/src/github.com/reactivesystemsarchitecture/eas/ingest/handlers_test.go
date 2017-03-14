package ingest

import (
	"testing"
	"github.com/reactivesystemsarchitecture/eas/protocol"
	"errors"
)

var f = EnvelopeHandlerFunc(func(_ *protocol.Envelope) error {
	return errors.New("bab")
})

func TestEnvelopeHandlerFunc_Handle(t *testing.T) {
	var e protocol.Envelope
	if f.Handle(&e).Error() != "bab" {
		t.Fail()
	}
}
