package service

import (
	"errors"
	"fmt"
	"log"

	client "github.com/denjons/RoboViewer/robot_service/client"
	db "github.com/denjons/RoboViewer/robot_service/database"
)

// RobotService handles events related to Robots and floors
type RobotService struct {
	repository *db.Repository
}

//NewRobotService creates a new RobotService
func NewRobotService(repository *db.Repository) (*RobotService, error) {
	if repository == nil {
		return nil, errors.New("repository cannot be nil")
	}
	return &RobotService{repository: repository}, nil
}

// CreateRobot a new Robot in the database
func (rs *RobotService) CreateRobot(robotDto *client.RobotDTO) (*db.RobotID, error) {
	err := rs.ValidateRobotDTO(robotDto)
	if err != nil {
		return nil, err
	}
	robotID, err := rs.repository.CreateRobot(robotDto.Name, robotDto.Width, robotDto.Height)
	if err != nil {
		return nil, err
	}

	log.Printf("Created robot %v", robotID.Value)

	return robotID, nil
}

//CreateFloor creates a new Floor in the database
func (rs *RobotService) CreateFloor(floorDTO *client.FloorDTO) (*db.FloorID, error) {
	err := rs.ValidateFloorDTO(floorDTO)
	if err != nil {
		return nil, err
	}

	floorID, err := rs.repository.CreateFloor(floorDTO.Name, floorDTO.Width, &floorDTO.Grid)
	if err != nil {
		return nil, err
	}

	log.Printf("Created floor %v", floorID.Value)

	return floorID, nil
}

// MapRobotToFloor maps a robot to a floor in the database
func (rs *RobotService) MapRobotToFloor(rfMap *client.RobotToFloorMapDTO) error {

	err := rs.ValidateID(rfMap.RobotID.Value)
	if err != nil {
		return err
	}

	floorIDerr := rs.ValidateID(rfMap.FloorID.Value)
	if floorIDerr != nil {
		return floorIDerr
	}

	log.Printf("Mapping robot %v to floor %v", rfMap.RobotID, rfMap.FloorID)

	return rs.repository.MapRobotToFloor(&db.RobotID{Value: rfMap.RobotID.Value}, &db.FloorID{Value: rfMap.FloorID.Value})
}

//ValidateID validates that an ID is of the correct size
func (rs *RobotService) ValidateID(id string) error {
	length := len(id)
	if length != 36 {
		return fmt.Errorf("ID length must be 36, but is: %v", length)
	}
	return nil
}

// ValidateFloorDTO validates that all fields are correct
func (rs *RobotService) ValidateFloorDTO(floorDTO *client.FloorDTO) error {

	if floorDTO == nil {
		return errors.New("FloorDTO cannot be nil")
	}

	nameLength := len(floorDTO.Name)
	if nameLength <= 0 || nameLength > 100 {
		return fmt.Errorf("FloorDTO Name length x must be 0 < x < 100, but is: %v", nameLength)
	}

	if floorDTO.Width <= 0 {
		return fmt.Errorf("FloorDTO Width x must be x > 0, but is: %v", floorDTO.Width)
	}

	gridLength := len(floorDTO.Grid)
	if gridLength == 0 {
		return errors.New("FloorDTO Grid must be over 0 cells")
	}

	return nil
}

// ValidateRobotDTO validates that all fields are correct
func (rs *RobotService) ValidateRobotDTO(robotDto *client.RobotDTO) error {
	if robotDto == nil {
		return errors.New("RobotDto cannot be nil")
	}

	nameLength := len(robotDto.Name)
	if nameLength <= 0 || nameLength > 100 {
		return fmt.Errorf("RobotDto Name length x must be 0 < x < 100, but is: %v", nameLength)
	}

	if robotDto.Width <= 0 {
		return fmt.Errorf("RobotDto Width x must be x > 0, but is: %v", robotDto.Width)
	}

	if robotDto.Height <= 0 {
		return fmt.Errorf("RobotDto Height x must be x > 0, but is: %v", robotDto.Height)
	}

	return nil
}
