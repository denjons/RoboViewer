package kafka

// PositionUpdateEvent represents position update from a robot
type PositionUpdateEvent struct {
	Sequence  int64
	RobotID   string
	SessionID string
	Position  Position
}

// Position represents a robots position
type Position struct {
	X int32
	Y int32
}

// SessionUpdateEvent represents a session update from a robot
type SessionUpdateEvent struct {
	SessionID    string
	SessionState SessionState
}

// SessionState represents the state of session
type SessionState string

const (
	// STARTED indicates that a session is started
	STARTED SessionState = "STARTED"
	// FINISHED indicates that a session is finished
	FINISHED SessionState = "FINISHED"
)
