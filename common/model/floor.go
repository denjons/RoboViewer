package model

import "errors"

// Floor represents a floor that can be cleaned
type Floor struct {
	Name        *string
	grid        *[]int
	width       int
	coveredArea int
}

// Marker handles basic operations which can be made on a Floor
type Marker interface {
	Mark(p *Point, r *Robot) error
	MarkPoint(point *Point) error
	GetCoveredAreaInPercent() float32
	Size() int
}

// NewFloor creates a new Floor
func NewFloor(name *string, width int, height int) (*Floor, error) {

	if name == nil {
		return nil, errors.New("name cannot be nil")
	}

	if *name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if width <= 0 {
		return nil, errors.New("width must be positive")
	}
	if height <= 0 {
		return nil, errors.New("height must be positive")
	}
	grid := make([]int, width*height)

	return &Floor{name, &grid, width, 0}, nil
}

// Mark the position of the robot on this floor
func (floor *Floor) Mark(point *Point, robot *Robot) error {
	shape := *robot.GetSahpe()
	for i := range shape {
		p := shape[i].Translate(point)
		err := floor.MarkPoint(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// MarkPoint markes a ponit on the floor as covered
func (floor *Floor) MarkPoint(point *Point) error {
	pos := floor.width*point.Y + point.X

	if pos >= len(*floor.grid) {
		return errors.New("point is outside of the floor grid")
	}

	if (*floor.grid)[pos] == 0 {
		(*floor.grid)[pos] = 1
		floor.coveredArea++
	}

	return nil
}

// GetCoveredAreaInPercent returns how much of the floor area that had been covered by a robot
func (floor *Floor) GetCoveredAreaInPercent() float32 {
	return 100.0 * (float32(floor.coveredArea) / float32(len(*floor.grid)))
}

//Size of the floor grid
func (floor *Floor) Size() int {
	return len(*floor.grid)
}
