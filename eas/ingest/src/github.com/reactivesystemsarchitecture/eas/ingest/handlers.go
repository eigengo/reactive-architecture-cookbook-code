package ingest

import "github.com/reactivesystemsarchitecture/eas/protocol"

// Defines a function that can be used to handle the envelope
type EnvelopeHandlerFunc func(envelope *protocol.Envelope) error

// A Handler interface that matches the `EnvelopeHandlerFunc`
type EnvelopeHandler interface {
	Handle(envelope *protocol.Envelope) error
}

type EnvelopeValidatorFunc func(envelope *protocol.Envelope) error

type EnvelopeValidator interface {
	Validate(envelope *protocol.Envelope) error
}

type EnvelopeProcessor interface {
	EnvelopeHandler
	EnvelopeValidator
}

type envelopeProcessor struct {
	validator EnvelopeValidator
	handler EnvelopeHandler
}

func (f envelopeProcessor) Handle(envelope *protocol.Envelope) error {
	return f.handler.Handle(envelope)
}

func (f envelopeProcessor) Validate(envelope *protocol.Envelope) error {
	return f.validator.Validate(envelope)
}

func (f EnvelopeValidatorFunc) Validate(envelope *protocol.Envelope) error {
	return f(envelope)
}

// Handle calls f(session).
func (f EnvelopeHandlerFunc) Handle(envelope *protocol.Envelope) error {
	return f(envelope)
}

func NewEnvelopeProcessor(validator EnvelopeValidator, handler EnvelopeHandler) EnvelopeProcessor {
	return envelopeProcessor {
		validator: validator,
		handler: handler,
	}
}
