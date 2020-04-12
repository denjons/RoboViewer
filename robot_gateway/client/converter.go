package client

import (
	pb "github.com/denjons/RoboViewer/robot_gateway/client/grpc/positionreport"
	rgKafka "github.com/denjons/RoboViewer/robot_gateway/client/kafka"
)

// ConvertToPositionUpdateEvent from PositionUpdate to PositionUpdateEvent
func ConvertToPositionUpdateEvent(positionUpdate *pb.PositionUpdate) *rgKafka.PositionUpdateEvent {
	return &rgKafka.PositionUpdateEvent{
		Sequence:  positionUpdate.SequenceNumber.Value,
		RobotID:   positionUpdate.RobotId.Value,
		SessionID: positionUpdate.SessionId.Value,
		Position:  rgKafka.Position{X: positionUpdate.Position.X, Y: positionUpdate.Position.Y},
	}
}
