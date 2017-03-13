package main

import (
	"flag"
	"net/http"
	"github.com/gorilla/mux"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	"github.com/reactivesystemsarchitecture/eas/ingest"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"fmt"
	"github.com/reactivesystemsarchitecture/eas/ingest/cassandra"
)

var (
	addr = flag.String("addr", ":8080", "The address to bind to")
)

func main() {
	flag.Parse()

	/*
	var envelope p.Envelope
	var session sp.Session
	var sensorData sp.SensorData
	var sensor sp.Sensor
	sensor.Location = sp.SensorLocation_Wrist
	sensor.DataTypes = []sp.SensorDataType{sp.SensorDataType_Acceleration}
	sensorData.Sensors = []*sp.Sensor{&sensor}
	sensorData.Values = make([]float32, 400, 400)
	session.SessionId = "random-string"
	session.SensorData = &sensorData

	a, _ := ptypes.MarshalAny(&session)
	envelope.Payload = a

	eb, _ := proto.Marshal(&envelope)
	log.Println("Marshalled ", eb)

	var envelope2 p.Envelope
	var session2 sp.Session
	proto.Unmarshal(eb, &envelope2)
	ptypes.UnmarshalAny(envelope2.Payload, &session2)

	log.Println("E2 ", envelope2, session2)
	*/

	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	if handler, err := cassandra.NewCluster(); err == nil {
		r.Handle("/session", PostSessionHandler(handler)).
			Methods("POST").
			Headers("Content-Type", "application/x-protobuf")

		// Our application will run on port 3000. Here we declare the port and pass in our router.
		log.Fatal(http.ListenAndServe(*addr, r))
	} else {
		log.Fatal(err)
	}

}

func PostSessionHandler(envelopeHandler ingest.EnvelopeHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var envelope p.Envelope
			if umerr := proto.Unmarshal(body, &envelope); umerr == nil {
				if herr := envelopeHandler.Validate(&envelope); herr != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", herr.Error())))
				} else if herr := envelopeHandler.Handle(&envelope); herr != nil {
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

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
})
