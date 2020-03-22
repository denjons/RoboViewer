package model

import (
	pb "github.com/denjons/RoboViewer/common/grpc/positionreport"
)

// ConvertToPositionUpdateEvent from PositionUpdate to PositionUpdateEvent
func ConvertToPositionUpdateEvent(positionUpdate *pb.PositionUpdate) *PositionUpdateEvent {
	return &PositionUpdateEvent{
		Sequence:  positionUpdate.SequenceNumber.Value,
		RobotID:   positionUpdate.RobotId.Value,
		SessionID: positionUpdate.SessionId.Value,
		Position:  Position{X: positionUpdate.Position.X, Y: positionUpdate.Position.Y},
	}
}
