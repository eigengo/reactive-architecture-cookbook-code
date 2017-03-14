package ingest

import (
	"bytes"
	"testing"
	"strings"
	"net/http"
	"errors"
	"io"
	"net/http/httptest"
	"github.com/golang/protobuf/proto"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
)

type testEnvelopeHandler struct {
	validationError error
	handleError     error
	handledEnvelope *p.Envelope
}

type brokenReader struct {}

func (b *brokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Bantha poodoo!")
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

func TestPostSessionHandlerBrokenBody(t *testing.T) {
	res := postSession(&brokenReader{}, newEnvelopeHandler())
	if res.StatusCode != http.StatusInternalServerError {
		t.Fail()
	}
}
