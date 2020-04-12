package client

// RobotDTO is a dto for adding a robot to the service
type RobotDTO struct {
	Name   string
	Width  int
	Height int
}

// FloorDTO is a dto dor adding a floor the service
type FloorDTO struct {
	Name  string
	Width int
	Grid  []int
}

//RobotToFloorMapDTO is a DTO for mapping floor to a robot to a floor
type RobotToFloorMapDTO struct {
	RobotID RobotID
	FloorID FloorID
}

// RobotID represents a robot ID outside the service
type RobotID struct {
	Value string
}

// FloorID represents a floor ID outside the service
type FloorID struct {
	Value string
}
