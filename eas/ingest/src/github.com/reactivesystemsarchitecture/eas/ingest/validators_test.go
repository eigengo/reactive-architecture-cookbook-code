package ingest

import (
	"testing"
	"github.com/stretchr/testify/assert"
	p "github.com/reactivesystemsarchitecture/eas/protocol"
	sp "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/proto"
)

func assertFailWithCode(t *testing.T, payload proto.Message, code int) {
	ee := &ValidationError{code: code}
	var e p.Envelope
	if payload != nil {
		a, _ := ptypes.MarshalAny(payload)
		e.Payload = a
	}

	if err, ok := SessionEnvelopeValidator.Validate(&e).(*ValidationError); ok {
		assert.Equal(t, err, ee)
	} else {
		assert.Fail(t, "Bad error")
	}
}

func TestValidationError_Error(t *testing.T) {
	assert.Equal(t, (&ValidationError{code: 123}).Error(), "{\"validation_error\":123}")
}

func TestSessionEnvelopeValidator_NoSession(t *testing.T) {
	assertFailWithCode(t, nil, ValidationErrorNilPayload)
}

func TestSessionEnvelopeValidator_NotSession(t *testing.T) {
	var l sp.Label
	assertFailWithCode(t, &l, ValidationErrorPayloadNotSession)
}

func TestSessionEnvelopeValidator_EmptySessionId(t *testing.T) {
	var s sp.Session
	assertFailWithCode(t, &s, ValidationErrorEmptySessionId)
}

func TestSessionEnvelopeValidator_EmptySessionSensorData(t *testing.T) {
	var s sp.Session
	s.SessionId = "abc"
	assertFailWithCode(t, &s, ValidationErrorEmptySensorData)
}

func TestSessionEnvelopeValidator_EmptySessionSensorDataValues(t *testing.T) {
	var s sp.Session
	var sd sp.SensorData
	s.SensorData = &sd
	s.SessionId = "abc"
	assertFailWithCode(t, &s, ValidationErrorEmptySensorDataValues)
}

func TestSessionEnvelopeValidator_EmptySessionSensors(t *testing.T) {
	var s sp.Session
	var sd sp.SensorData
	sd.Values = []float32{0}
	sd.Sensors = []*sp.Sensor{}
	s.SensorData = &sd
	s.SessionId = "abc"
	assertFailWithCode(t, &s, ValidationErrorEmptySensors)
}

func TestSessionEnvelopeValidator_MisalignedSensorValues(t *testing.T) {
	var s sp.Session
	var sd sp.SensorData
	var sn sp.Sensor
	sn.DataTypes = []sp.SensorDataType{sp.SensorDataType_Acceleration}
	sd.Values = []float32{0, 1}
	sd.Sensors = []*sp.Sensor{&sn}
	s.SensorData = &sd
	s.SessionId = "abc"
	assertFailWithCode(t, &s, ValidationErrorMisalignedSensorValues)
}

func TestSessionEnvelopeValidator_Valid(t *testing.T) {
	var s sp.Session
	var sd sp.SensorData
	var sn1, sn2 sp.Sensor
	sn1.DataTypes = []sp.SensorDataType{sp.SensorDataType_Acceleration, sp.SensorDataType_HeartRate}
	sn2.DataTypes = []sp.SensorDataType{sp.SensorDataType_Rotation}
	sd.Values = []float32{0, 1, 3, 4, 10, 11, 13}
	sd.Sensors = []*sp.Sensor{&sn1, &sn2}
	s.SensorData = &sd
	s.SessionId = "abc"

	var e p.Envelope
	a, _ := ptypes.MarshalAny(&s)
	e.Payload = a

	assert.Nil(t, SessionEnvelopeValidator.Validate(&e))
}
