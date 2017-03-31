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
	"compress/gzip"
	"io/ioutil"
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

func (t *testEnvelopeProcessor) Handle(envelope *p.Envelope) error {
	t.handledEnvelope = envelope
	return t.handleError
}

func (t *testEnvelopeProcessor) Validate(envelope *p.Envelope) error {
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

func postSession(body io.Reader, transferEncoding string, processor EnvelopeProcessor) *http.Response {
	rw := httptest.NewRecorder()
	rq := httptest.NewRequest("POST", "/session", body)
	rq.Header.Set("Transfer-Encoding", transferEncoding)
	PostSessionHandler(processor).ServeHTTP(rw, rq)
	return rw.Result()
}

func newEnvelopeReader() io.Reader {
	var e p.Envelope
	e.CorrelationId = "const-correlation-id"
	b, _ := proto.Marshal(&e)
	return bytes.NewReader(b)
}

func gzipped(reader io.Reader) io.Reader {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	readerBytes, _ := ioutil.ReadAll(reader)
	w.Write(readerBytes)
	w.Close()

	return bytes.NewReader(b.Bytes())
}

func TestPostSessionHandlerNoProtobuf(t *testing.T) {
	res := postSession(strings.NewReader("(not-protobuf"), "octet-stream", newEnvelopeProcessor())
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestPostSessionHandlerOKRaw(t *testing.T) {
	processor := newEnvelopeProcessor()
	res := postSession(newEnvelopeReader(), "octet-stream", processor)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "const-correlation-id", processor.handledEnvelope.CorrelationId)
}

func TestPostSessionHandlerOKGZipped(t *testing.T) {
	processor := newEnvelopeProcessor()
	res := postSession(gzipped(newEnvelopeReader()), "gzip", processor)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "const-correlation-id", processor.handledEnvelope.CorrelationId)
}

func TestPostSessionHandlerWithFailedValidation(t *testing.T) {
	res := postSession(newEnvelopeReader(), "octet-stream", newEnvelopeProcessor().withValidationError(someError))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestPostSessionHandlerWithFailedHandle(t *testing.T) {
	res := postSession(newEnvelopeReader(), "octet-stream", newEnvelopeProcessor().withHandleError(someError))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestPostSessionHandlerBrokenBody(t *testing.T) {
	res := postSession(&brokenReader{}, "octet-stream", newEnvelopeProcessor())
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
