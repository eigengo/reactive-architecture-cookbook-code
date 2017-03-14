package ingest

import (
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"net/http"
	"fmt"

	p "github.com/reactivesystemsarchitecture/eas/protocol"
)

func PostSessionHandler(envelopeProcessor EnvelopeProcessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var envelope p.Envelope
			if umerr := proto.Unmarshal(body, &envelope); umerr == nil {
				if herr := envelopeProcessor.Validate(&envelope); herr != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", herr.Error())))
				} else if herr := envelopeProcessor.Handle(&envelope); herr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", herr.Error())))
				} else {
					w.Write([]byte("{}"))
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", umerr.Error())))
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", err.Error())))
		}
	})
}
/*
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
})
*/