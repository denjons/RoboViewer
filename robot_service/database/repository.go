package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// RobotRepository handles databse actions on Robots
type RobotRepository interface {
	CreateRobot(name string, width int, height int) (*RobotID, error)
	CreateFloor(floor *Floor) (*FloorID, error)
	MapRobotToFloor(robotID *RobotID, floorID *FloorID) error
}

// SessionRepository handles databse actions on Floors
type SessionRepository interface {
	CreateSession(sessionID *SessionID)
	GetSession(sessionID *SessionID)
	InsertPosition(session Session)
}

// Repository represents a repository
type Repository struct {
	ctx context.Context
	db  *sql.DB
}

// CreateRobot creates a new robot
func (r *Repository) CreateRobot(name string, width int, height int) (*RobotID, error) {

	id, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	stmt, err := r.db.Prepare("INSERT INTO robots(uuid, name, width, height) VALUES( ?, ?, ?, ? )")

	if err != nil {
		return nil, err
	}

	_, execErr := stmt.Exec(id.String(), name, width, height)

	if execErr != nil {
		return nil, execErr
	}

	return &RobotID{Value: id.String()}, nil
}

// CreateFloor creates a new floor
func (r *Repository) CreateFloor(name string, width int, grid *[]int) (*FloorID, error) {

	id, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	stmt, err := r.db.Prepare("INSERT INTO floors(uuid, name, width, grid) VALUES( ?, ?, ?, ? )")

	if err != nil {
		return nil, err
	}

	array := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(&grid)), ", "), "[]")

	_, execErr := stmt.Exec(id.String(), name, width, "{"+array+"}")

	if execErr != nil {
		return nil, execErr
	}

	return &FloorID{Value: id.String()}, nil
}

// MapRobotToFloor maps a robot to a floor
func (r *Repository) MapRobotToFloor(robot *Robot, floor *Floor) error {
	stmt, err := r.db.Prepare("UPDATE Robots SET floorId = ? WHERE id ?")

	if err != nil {
		return err
	}

	_, execErr := stmt.Exec(floor.id, robot.id)

	return execErr
}

func (r *Repository) executeInTransaction(query string, args ...interface{}) (*sql.Result, error) {
	tx, err := r.db.BeginTx(r.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		return nil, err
	}

	result, execErr := tx.Exec(query, args)

	if execErr != nil {
		_ = tx.Rollback()
		return nil, execErr
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &result, nil

}
