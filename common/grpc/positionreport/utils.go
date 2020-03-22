package positionreport

import uuid "github.com/satori/go.uuid"

// CreateTestPositionUpdate creates a simple PositionUpdate for testing
func CreateTestPositionUpdate(x int32, y int32, sequence int64) (*PositionUpdate, error) {
	sessionUUID := uuid.NewV1()

	robotUUID := uuid.NewV1()

	sessionID := SessionId{
		Value: sessionUUID.String(),
	}

	robotID := RobotId{
		Value: robotUUID.String(),
	}

	sequenceNumber := SequenceNumber{Value: sequence}

	position := Position{
		X: x,
		Y: y,
	}

	return &PositionUpdate{
		SessionId:      &sessionID,
		RobotId:        &robotID,
		SequenceNumber: &sequenceNumber,
		Position:       &position,
	}, nil
}
