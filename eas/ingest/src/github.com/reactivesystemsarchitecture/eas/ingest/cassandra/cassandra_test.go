package cassandra

import (
	"testing"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	sp "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
	"log"
	"github.com/golang/protobuf/proto"
	"compress/gzip"
	"bytes"
	"math/rand"
)

func newEnvelopeWithSession() p.Envelope {
	var s sp.Session
	var sd sp.SensorData
	var sn1, sn2 sp.Sensor
	count := 50 * 3600
	sn1.DataTypes = []sp.SensorDataType{sp.SensorDataType_Acceleration, sp.SensorDataType_Rotation, sp.SensorDataType_HeartRate}
	sn2.DataTypes = []sp.SensorDataType{sp.SensorDataType_Acceleration}

	sampleCount := count * 10
	sd.Values = make([]float32, sampleCount, sampleCount)
	for i := 0; i < sampleCount; i++ {
		sd.Values[i] = rand.Float32()
	}

	sd.Sensors = []*sp.Sensor{&sn1, &sn2}
	s.SensorData = &sd
	sessionId, _ := uuid.NewRandom()
	s.SessionId = sessionId.String()

	var e p.Envelope
	a, _ := ptypes.MarshalAny(&s)
	e.Payload = a

	return e
}

func TestNewPersistSessionEnvelopeHandler(t *testing.T) {
	e := newEnvelopeWithSession()
	h, err := NewPersistSessionEnvelopeHandler("localhost")
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.NoError(t, h.Handle(&e))
}

func BenchmarkNewPersistSessionEnvelopeHandler10(b *testing.B) {
	log.Println(b.N)

	b.StopTimer()
	h, err := NewPersistSessionEnvelopeHandler("localhost")
	if err != nil {
		assert.FailNow(b, err.Error())
	}

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		e := newEnvelopeWithSession()
		assert.NoError(b, h.Handle(&e))
	}
}
