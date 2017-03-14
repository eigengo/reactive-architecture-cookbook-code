package ingest

import (
	"github.com/reactivesystemsarchitecture/eas/protocol"
	ps "github.com/reactivesystemsarchitecture/eas/protocol/session/v1m0"

	"errors"
	"github.com/golang/protobuf/ptypes"
)

var SessionEnvelopeValidator EnvelopeValidator = EnvelopeValidatorFunc(func(envelope *protocol.Envelope) error {
	var session ps.Session
	if envelope.Payload == nil {
		return errors.New("envelope.Payload == nil")
	}
	if !ptypes.Is(envelope.Payload, &session) {
		return errors.New("Payload is not Session")
	}
	ptypes.UnmarshalAny(envelope.Payload, &session)

	if len(session.SessionId) == 0 {
		return errors.New("Missing session id")
	}
	if session.SensorData == nil {
		return errors.New("Missing sensor data")
	}

	if len(session.SensorData.Values) == 0 {
		return errors.New("Missing sensor data values")
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

	if len(session.SensorData.Values) % sampleSize != 0 {
		return errors.New("Sensor data values does not contain enough values for the sensor data")
	}

	return nil
})
