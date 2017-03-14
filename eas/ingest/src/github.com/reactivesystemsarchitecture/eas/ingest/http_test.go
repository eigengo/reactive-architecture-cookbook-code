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

type testEnvelopeProcessor struct {
	validationError error
	handleError     error
	handledEnvelope *p.Envelope
}

type brokenReader struct {}

func (b *brokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Bantha poodoo!")
}

func (t testEnvelopeProcessor) Handle(envelope *p.Envelope) error {
	t.handledEnvelope = envelope
	return t.handleError
}

func (t testEnvelopeProcessor) Validate(envelope *p.Envelope) error {
	return t.validationError
}

func (t *testEnvelopeProcessor) withValidationError(validationError error) *testEnvelopeProcessor {
	t.validationError = validationError
	return t
}

func (t *testEnvelopeProcessor) withHandleError(handleError error) *testEnvelopeProcessor {
	t.handleError = handleError
	return t
}

func newEnvelopeProcessor() *testEnvelopeProcessor {
	return &testEnvelopeProcessor{}
}

var someError = errors.New("some error")

func postSession(body io.Reader, processor *testEnvelopeProcessor) *http.Response {
	rw := httptest.NewRecorder()
	rq := httptest.NewRequest("POST", "/session", body)
	PostSessionHandler(*processor).ServeHTTP(rw, rq)
	return rw.Result()
}

func newEnvelopeReader() io.Reader {
	var e p.Envelope
	b, _ := proto.Marshal(&e)
	return bytes.NewReader(b)
}

func TestPostSessionHandlerNoProtobuf(t *testing.T) {
	res := postSession(strings.NewReader("(not-protobuf"), newEnvelopeProcessor())
	if res.StatusCode != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostSessionHandlerOK(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeProcessor())
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}
}

func TestPostSessionHandlerWithFailedValidation(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeProcessor().withValidationError(someError))
	if res.StatusCode != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostSessionHandlerWithFailedHandle(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeProcessor().withHandleError(someError))
	if res.StatusCode != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestPostSessionHandlerBrokenBody(t *testing.T) {
	res := postSession(&brokenReader{}, newEnvelopeProcessor())
	if res.StatusCode != http.StatusInternalServerError {
		t.Fail()
	}
}
