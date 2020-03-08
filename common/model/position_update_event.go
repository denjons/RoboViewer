package model

// PositionUpdateEvent is an internal represeantion of an position update from a robot
type PositionUpdateEvent struct {
	Sequence  int64
	RobotID   int32
	SessionID int32
	Position  *Position
}

// Position is an internal representaion of a robots position
type Position struct {
	X int32
	Y int32
}
