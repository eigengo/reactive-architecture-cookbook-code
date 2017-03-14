package main

import (
	"flag"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/reactivesystemsarchitecture/eas/ingest"
	"log"
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

	if handler, err := cassandra.NewEnvelopeHandler(":aa"); err == nil {
		r.Handle("/session", ingest.PostSessionHandler(handler)).
			Methods("POST").
			Headers("Content-Type", "application/x-protobuf")

		// Our application will run on port 3000. Here we declare the port and pass in our router.
		log.Fatal(http.ListenAndServe(*addr, r))
	} else {
		log.Fatal(err)
	}

}
