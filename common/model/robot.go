package model

import "errors"

//RobotShape handles operations on the shape of the robot
type RobotShape interface {
	GetShape() *[]Point
}

// Robot represents a rebot vacuum cleaner
type Robot struct {
	ID    string
	Name  string
	shape []Point
}

// NewRectangularRobot return s new robot in a square shape
func NewRectangularRobot(ID *string, name *string, width int, height int) (*Robot, error) {

	if name == nil {
		return nil, errors.New("Robot name cannot be nil")
	}

	if *name == "" {
		return nil, errors.New("Robot name cannot be empty")
	}

	if width <= 0 {
		return nil, errors.New("Robot width must be positive")
	}

	if height <= 0 {
		return nil, errors.New("Robot height must be positive")
	}

	if ID == nil {
		return nil, errors.New("Robot ID cannot be nil")
	}

	if *ID == "" {
		return nil, errors.New("Robot ID cannot be empty")
	}

	points := make([]Point, width*height)
	position := 0
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			points[position] = Point{i - width/2, j - height/2}
			position++
		}
	}

	return &Robot{*ID, *name, points}, nil
}

// GetSahpe returns a copy of the robots shape
func (robot *Robot) GetSahpe() *[]Point {
	points := robot.shape
	return &points
}
