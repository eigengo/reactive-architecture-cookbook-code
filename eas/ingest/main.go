package main

import (
	"flag"
	"net/http"
	"github.com/gorilla/mux"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	h "github.com/reactivesystemsarchitecture/eas"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
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

	r.Handle("/session", PostSessionHandler(h.EnvelopeHandlerFunc(func(e *p.Envelope) error {
		return nil
	}))).Methods("POST").Headers("Content-Type", "application/x-protobuf")

	// Our application will run on port 3000. Here we declare the port and pass in our router.
	http.ListenAndServe(*addr, r)
}

func PostSessionHandler(envelopeHandler h.EnvelopeHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Here we are converting the slice of products to json
		emptyBody := []byte("{}")
		w.Header().Set("Content-Type", "application/json")

		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var envelope p.Envelope
			if umerr := proto.Unmarshal(body, &envelope); umerr == nil {
				if herr := envelopeHandler.Handle(&envelope); herr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(":("))
				} else {
					w.Write(emptyBody)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(":("))
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(":("))
		}
	})
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
})
