package cassandra

import (
	"testing"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	sp "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
)

func newEnvelopeWithSession() p.Envelope {
	var s sp.Session
	var sd sp.SensorData
	var sn1, sn2 sp.Sensor
	sn1.DataTypes = []sp.SensorDataType{sp.SensorDataType_Acceleration, sp.SensorDataType_HeartRate}
	sn2.DataTypes = []sp.SensorDataType{sp.SensorDataType_Rotation}
	sd.Values = []float32{0, 1, 3, 4, 10, 11, 13}
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

	if err := h.Handle(&e); err != nil {
		assert.FailNow(t, err.Error())
	}
}
