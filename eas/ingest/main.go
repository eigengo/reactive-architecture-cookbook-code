package main

import (
	"flag"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/reactivesystemsarchitecture/eas/ingest"
	"log"
	"github.com/reactivesystemsarchitecture/eas/ingest/cassandra"
	"strings"
	"compress/gzip"
	"io"
	"github.com/stretchr/testify/require"
)

type stringsValue []string

func (*stringsValue) String() string {
	return ""
}

func (s *stringsValue) Set(value string) error {
	*s =  append(*s, value)
	return nil
}

func main() {
	var cassandraHosts stringsValue
	var addr string
	flag.Var(&cassandraHosts, "cassandra-hosts", "The cassandra hosts")
	flag.StringVar(&addr, "addr", ":8080", "The address to bind to")

	if len(cassandraHosts) == 0 {
		cassandraHosts = stringsValue{"localhost"}
	}

	flag.Parse()

	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	if handler, err := cassandra.NewPersistSessionEnvelopeHandler(cassandraHosts...); err == nil {
		processor := ingest.NewEnvelopeProcessor(ingest.SessionEnvelopeValidator, handler)

		r.Handle("/session", ingest.PostSessionHandler(processor)).
			Methods("POST").
			Headers("Content-Type", "application/x-protobuf")

		// Our application will run on the given addr.
		log.Fatal(http.ListenAndServe(addr, r))
	} else {
		log.Fatal(err)
	}

}
