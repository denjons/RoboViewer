package database

// Robot represents internal model for robot_service
type Robot struct {
	id      int
	ID      RobotID
	Name    string
	Created string
	Width   int
	Height  int
	floor   *Floor
}

// RobotID is a uuid for a Robot
type RobotID struct {
	Value string
}

// Floor represents internal model for robot_service
type Floor struct {
	id      int
	ID      FloorID
	Name    string
	Created string
	Grid    []int
	Width   int
}

// FloorID is a uuid for a Floor
type FloorID struct {
	Value string
}

// Session represents an ongoing cleaning Session
type Session struct {
	id      int
	ID      SessionID
	Created string
	Robot   *Robot
	Floor   *Floor
}

// SessionID is a uuid for a Session
type SessionID struct {
	Value string
}

// Point represents an X,Y position on the floor grid
type Point struct {
	X        int
	Y        int
	Sequence int64
}
