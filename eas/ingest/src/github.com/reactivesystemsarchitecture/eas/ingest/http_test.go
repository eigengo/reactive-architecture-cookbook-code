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
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestPostSessionHandlerOK(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeProcessor())
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestPostSessionHandlerWithFailedValidation(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeProcessor().withValidationError(someError))
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestPostSessionHandlerWithFailedHandle(t *testing.T) {
	res := postSession(newEnvelopeReader(), newEnvelopeProcessor().withHandleError(someError))
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
}

func TestPostSessionHandlerBrokenBody(t *testing.T) {
	res := postSession(&brokenReader{}, newEnvelopeProcessor())
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
}
