package main

import (
	"testing"
	"errors"
	"net/http/httptest"
	"strings"
	"net/http"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	"github.com/golang/protobuf/proto"
	"bytes"
	"io"
)

type testEnvelopeHandler struct {
	validationError error
	handleError     error
	handledEnvelope *p.Envelope
}

func (t testEnvelopeHandler) Handle(envelope *p.Envelope) error {
	t.handledEnvelope = envelope
	return t.handleError
}

func (t testEnvelopeHandler) Validate(envelope *p.Envelope) error {
	return t.validationError
}

func (t *testEnvelopeHandler) withValidationError(validationError error) *testEnvelopeHandler {
	t.validationError = validationError
	return t
}

func (t *testEnvelopeHandler) withHandleError(handleError error) *testEnvelopeHandler {
	t.handleError = handleError
	return t
}

func newEnvelopeHandler() *testEnvelopeHandler {
	return &testEnvelopeHandler{}
}

var someError = errors.New("some error")

func someEnvelope() io.Reader {
	var e p.Envelope
	b, _ := proto.Marshal(&e)
	return bytes.NewReader(b)
}

func TestPostSessionHandlerNoProtobuf(t *testing.T) {
	rw := httptest.NewRecorder()
	rq := httptest.NewRequest("POST", "/session", strings.NewReader("(not-protobuf)"))
	h := newEnvelopeHandler()
	PostSessionHandler(*h).ServeHTTP(rw, rq)

	if rw.Result().StatusCode != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostSessionHandlerOK(t *testing.T) {
	rw := httptest.NewRecorder()
	rq := httptest.NewRequest("POST", "/session", someEnvelope())
	PostSessionHandler(*newEnvelopeHandler()).ServeHTTP(rw, rq)

	if rw.Result().StatusCode != http.StatusOK {
		t.Fail()
	}
}
