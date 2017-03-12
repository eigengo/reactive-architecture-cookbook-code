package handlers

import "github.com/reactivesystemsarchitecture/eas/protocol"

type EnvelopeHandlerFunc func(envelope *protocol.Envelope) error

type EnvelopeHandler interface {
	Handle(envelope *protocol.Envelope) error
}

// Handle calls f(session).
func (f EnvelopeHandlerFunc) Handle(envelope *protocol.Envelope) error {
	return f(envelope)
}
