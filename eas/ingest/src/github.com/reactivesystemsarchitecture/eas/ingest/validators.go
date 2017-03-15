package ingest

import (
	"github.com/reactivesystemsarchitecture/eas/protocol"
	ps "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"

	"github.com/golang/protobuf/ptypes"
	"fmt"
)

const (
	ValidationErrorNilPayload             int = iota
	ValidationErrorPayloadNotSession
	ValidationErrorEmptySessionId
	ValidationErrorEmptySensorData
	ValidationErrorEmptySensorDataValues
	ValidationErrorMisalignedSensorValues

)

type ValidationError struct {
	code int
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Validation failed %d", v.code)
}

var SessionEnvelopeValidator EnvelopeValidator = EnvelopeValidatorFunc(func(envelope *protocol.Envelope) error {
	var session ps.Session
	if envelope.Payload == nil {
		return &ValidationError{code: ValidationErrorNilPayload}
	}
	if !ptypes.Is(envelope.Payload, &session) {
		return &ValidationError{code: ValidationErrorPayloadNotSession}
	}
	ptypes.UnmarshalAny(envelope.Payload, &session)

	if len(session.SessionId) == 0 {
		return &ValidationError{code: ValidationErrorEmptySessionId}
	}
	if session.SensorData == nil {
		return &ValidationError{code: ValidationErrorEmptySensorData}
	}

	if len(session.SensorData.Values) == 0 {
		return &ValidationError{code: ValidationErrorEmptySensorDataValues}
	}

	sampleSize := 0
	for _, sensor := range session.SensorData.Sensors {
		for _, dataType := range sensor.DataTypes {
			switch dataType {
			case ps.SensorDataType_Acceleration:
			case ps.SensorDataType_Rotation:
				sampleSize += 3
			case ps.SensorDataType_HeartRate:
				sampleSize += 1
			}
		}
	}

	if len(session.SensorData.Values)%sampleSize != 0 {
		return &ValidationError{code: ValidationErrorMisalignedSensorValues}
	}

	return nil
})
