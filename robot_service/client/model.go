package client

// RobotDTO is a dto for adding a robot to the service
type RobotDTO struct {
	ID    string
	Name  string
	Width int
	grid  []int
}

// FloorDTO is a dto dor adding a floor the service
type FloorDTO struct {
	Name  string
	Width int
	grid  []int
}
