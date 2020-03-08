package model_test

import (
	"testing"

	model "github.com/denjons/RoboViewer/common/model"
)

func TestNewRectangularRobot(t *testing.T) {
	name := "rect robot"
	robot, err := model.NewRectangularRobot(&name, 4, 4)

	if err != nil {
		t.Errorf("NewRectangularRobot() got error %v", err)
	}

	shape := *robot.GetSahpe()

	for i := range shape {
		if shape[i].X < -2 || shape[i].X > 2 {
			t.Errorf("NewRectangularRobot() Wrong x at %v. Got %v", i, shape[i].X)
		}
		if shape[i].Y < -2 || shape[i].Y > 2 {
			t.Errorf("NewRectangularRobot() Wrong y at %v. Got %v", i, shape[i].Y)
		}
	}

}

func TestNewRectangularRobotNilName(t *testing.T) {
	_, err := model.NewRectangularRobot(nil, 4, 4)

	evalErrorMessage(err, "Robot name cannot be nil", t)
}

func TestNewRectangularRobotEmptyName(t *testing.T) {
	name := ""
	_, err := model.NewRectangularRobot(&name, 4, 4)

	evalErrorMessage(err, "Robot name cannot be empty", t)
}

func TestNewRectangularRobotZeroWidth(t *testing.T) {
	name := "robot"
	_, err := model.NewRectangularRobot(&name, 0, 4)

	evalErrorMessage(err, "Robot width must be positive", t)
}

func TestNewRectangularRobotZeroHeight(t *testing.T) {
	name := "robot"
	_, err := model.NewRectangularRobot(&name, 4, 0)

	evalErrorMessage(err, "Robot height must be positive", t)
}
