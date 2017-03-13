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

func postSession(body io.Reader, handler *testEnvelopeHandler) *http.Response {
	rw := httptest.NewRecorder()
	rq := httptest.NewRequest("POST", "/session", body)
	PostSessionHandler(*handler).ServeHTTP(rw, rq)
	return rw.Result()
}

func newEnvelopeReader() io.Reader {
	var e p.Envelope
	b, _ := proto.Marshal(&e)
	return bytes.NewReader(b)
}

func TestPostSessionHandlerNoProtobuf(t *testing.T) {
	res := postSession(strings.NewReader("(not-protobuf"), newEnvelopeHandler())
	if res.StatusCode != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostSessionHandlerOK(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeHandler())
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}
}

func TestPostSessionHandlerWithFailedValidation(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeHandler().withValidationError(someError))
	if res.StatusCode != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostSessionHandlerWithFailedHandle(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeHandler().withHandleError(someError))
	if res.StatusCode != http.StatusInternalServerError {
		t.Fail()
	}
}
