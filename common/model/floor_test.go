package model_test

import (
	"testing"

	model "github.com/denjons/RoboViewer/common/model"
)

func TestNewFloorNilName(t *testing.T) {
	_, err := model.NewFloor(nil, 10, 10)

	evalErrorMessage(err, "name cannot be nil", t)
}

func TestNewFloorEmptyName(t *testing.T) {
	name := ""

	_, err := model.NewFloor(&name, 10, 10)

	evalErrorMessage(err, "name cannot be empty", t)
}

func TestNewFloorZeroWidth(t *testing.T) {
	name := "living room"

	_, err := model.NewFloor(&name, 0, 10)

	evalErrorMessage(err, "width must be positive", t)
}

func TestNewFloorZeroHeight(t *testing.T) {
	name := "living room"

	_, err := model.NewFloor(&name, 10, 0)

	evalErrorMessage(err, "height must be positive", t)
}

func TestNewFloor(t *testing.T) {
	name := "living room"
	width := 40
	height := 30
	size := width * height
	floor, err := model.NewFloor(&name, width, height)

	if err != nil {
		t.Errorf("NewFloor() got error %v", err)
	}

	if floor == nil {
		t.Error("NewFloor() returned nil")
	}

	if floor.Size() != size {
		t.Errorf("Size() is %v, want %v", floor.Size(), size)
	}

	num := floor.GetCoveredAreaInPercent()

	if num != 0.0 {
		t.Errorf("GetCoveredAreaInPercent() is %v, want %v", num, 0.0)
	}
}

func TestMarkPoint(t *testing.T) {

	floor := createFloor(40, 30, t)

	point := &model.Point{X: 1, Y: 1}

	err := floor.MarkPoint(point)

	if err != nil {
		t.Errorf("MarkPoint() got error %v", err)
	}
}

func TestGetCoveredAreaInPercentZero(t *testing.T) {
	floor := createFloor(40, 30, t)

	zeroCoveredArea := floor.GetCoveredAreaInPercent()

	if zeroCoveredArea > 0.0 {
		t.Errorf("GetCoveredAreaInPercent() vanted 0.0 got %v", zeroCoveredArea)
	}
}

func TestGetCoveredAreaInPercentHalf(t *testing.T) {
	floor := createFloor(10, 10, t)

	for i := 0; i < 5; i++ {
		for j := 0; j < 10; j++ {
			err := floor.MarkPoint(&model.Point{i, j})
			if err != nil {
				t.Errorf("MarkPoint() got error %v", err)
			}
		}
	}

	coveredArea := floor.GetCoveredAreaInPercent()

	if coveredArea != 50.0 {
		t.Errorf("GetCoveredAreaInPercent() wanted 50.0 got %v", coveredArea)
	}

}

func TestMark(t *testing.T) {
	robot := createRobot(4, 4, t)
	floor := createFloor(8, 8, t)
	err := floor.Mark(&model.Point{2, 2}, robot)

	if err != nil {
		t.Errorf("Mark() got error %v", err)
	}

	coveredArea := floor.GetCoveredAreaInPercent()

	if coveredArea != 25.0 {
		t.Errorf("GetCoveredAreaInPercent() wanted 25 but got %v", coveredArea)
	}
}

func createRobot(width int, height int, t *testing.T) *model.Robot {
	name := "square robot"
	robot, err := model.NewRectangularRobot(&name, 4, 4)

	if err != nil {
		t.Errorf("NewSqueareRobot() got error %v", err)
	}

	return robot
}

func createFloor(width, height int, t *testing.T) *model.Floor {
	name := "living room"
	floor, err := model.NewFloor(&name, width, height)

	if err != nil {
		t.Errorf("NewFloor() got error %v", err)
	}

	return floor
}

func evalErrorMessage(err error, expected string, t *testing.T) {
	if err == nil {
		t.Errorf("wanted error but got nil")
	}
	if err.Error() != expected {
		t.Errorf("wanted error '%v' but got %v", expected, err.Error())
	}
}
