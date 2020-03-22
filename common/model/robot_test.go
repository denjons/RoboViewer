package model_test

import (
	"testing"

	model "github.com/denjons/RoboViewer/common/model"
)

var (
	robot_ID   = "id"
	robot_name = "test_robot"
)

func TestNewRectangularRobot(t *testing.T) {
	robot, err := model.NewRectangularRobot(robot_ID, robot_name, 4, 4)

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

func TestNewRectangularRobotEmptyName(t *testing.T) {
	_, err := model.NewRectangularRobot(robot_ID, "", 4, 4)

	evalErrorMessage(err, "Robot name cannot be empty", t)
}

func TestNewRectangularRobotEmptyID(t *testing.T) {
	_, err := model.NewRectangularRobot("", robot_name, 4, 4)

	evalErrorMessage(err, "Robot ID cannot be empty", t)
}

func TestNewRectangularRobotZeroWidth(t *testing.T) {
	_, err := model.NewRectangularRobot(robot_ID, robot_name, 0, 4)

	evalErrorMessage(err, "Robot width must be positive", t)
}

func TestNewRectangularRobotZeroHeight(t *testing.T) {
	_, err := model.NewRectangularRobot(robot_ID, robot_name, 4, 0)

	evalErrorMessage(err, "Robot height must be positive", t)
}
