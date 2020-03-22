package model_test

import (
	"testing"

	pb "github.com/denjons/RoboViewer/common/grpc/positionreport"
	model "github.com/denjons/RoboViewer/common/model"
	assert "github.com/stretchr/testify/assert"
)

func TestConvertToPositionUpdateEvent(t *testing.T) {
	positionUpdate, err := pb.CreateTestPositionUpdate(1, 2, 3)

	if err != nil {
		t.Fatalf("Error creating PositionUpdate: %v", err)
	}

	positionUpdateEvent := model.ConvertToPositionUpdateEvent(positionUpdate)

	assert.NotNil(t, positionUpdateEvent.Position)

	x := positionUpdateEvent.Position.X
	assert.Equal(t, int32(1), x, "Position x should be same")

	y := positionUpdateEvent.Position.Y
	assert.Equal(t, int32(2), y, "Position y should be same")

	sequence := positionUpdateEvent.Sequence
	assert.Equal(t, int64(3), sequence, "Sequence count should be same")

	sessionId := positionUpdateEvent.SessionID
	assert.Equal(t, positionUpdate.SessionId.Value, sessionId, "SessionId should be same")

	robotId := positionUpdateEvent.RobotID
	assert.Equal(t, positionUpdate.RobotId.Value, robotId, "SessionId should be same")
}
