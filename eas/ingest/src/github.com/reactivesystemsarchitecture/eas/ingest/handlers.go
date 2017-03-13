package ingest

import "github.com/reactivesystemsarchitecture/eas/protocol"

// Defines a function that can be used to handle the envelope
type EnvelopeHandlerFunc func(envelope *protocol.Envelope) error

// A Handler interface that matches the `EnvelopeHandlerFunc`
type EnvelopeHandler interface {
	Validate(envelope *protocol.Envelope) error
	Handle(envelope *protocol.Envelope) error
}

// Handle calls f(session).
func (f EnvelopeHandlerFunc) Handle(envelope *protocol.Envelope) error {
	return f(envelope)
}

// HandlerFunc takes all input
func (f EnvelopeHandlerFunc) Validate(envelope *protocol.Envelope) error {
	return nil
}