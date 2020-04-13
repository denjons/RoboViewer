package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// RobotRepositoryWriter handles databse write actions on Robots and floors
type RobotRepositoryWriter interface {
	CreateRobot(name string, width int, height int) (*RobotID, error)
	CreateFloor(floor *Floor) (*FloorID, error)
	MapRobotToFloor(robotID *RobotID, floorID *FloorID) error
}

// SessionRepository handles databse actions on Floors
type SessionRepository interface {
	CreateSession(sessionID *SessionID, robot *Robot) error
	GetSession(sessionID *SessionID) (*Session, error)
	InsertPosition(session *Session, point *Point) error
	GetRobot(robotID *RobotID) (*Robot, error)
}

// Repository represents a repository
type Repository struct {
	ctx         context.Context
	stopChannel context.CancelFunc
	db          *sql.DB
}

// NewRepository return a new Repository that is connected
func NewRepository(parent context.Context, DSN string) (*Repository, error) {
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(parent)
	return &Repository{db: db, ctx: ctx, stopChannel: cancel}, nil
}

// Stop stops all connections
func (r *Repository) Stop() {
	r.db.Close()
	r.stopChannel()
}

// CreateSession creates a new session
func (r *Repository) CreateSession(sessionID *SessionID, robot *Robot, floor *Floor) error {
	stmt, err := r.db.Prepare("INSERT INTO sessions(uuid, robotid, floorid) VALUES( $1, $2, $3)")

	if err != nil {
		return err
	}

	_, execErr := stmt.Exec(sessionID.Value, robot.id, floor.id)

	return execErr
}

/*
GetSession gets a session from the database
If an error occurs error is returned.
*/
func (r *Repository) GetSession(sessionID *SessionID) (*Session, error) {
	var ID, robotID, floorID int
	var robotWidth, robotHeight, floorWidth int
	var sessionUUID, robotUUID, floorUUID, sessionCreated, robotCreated, floorCreated, robotName, floorName string
	var floorGrid []int

	err := r.db.QueryRowContext(r.ctx, `SELECT s.id as session_id, s.uuid as session_uuid, s.created as session_created, 
	r.id as robot_id, r.uuid as robot_uuid, r.created as robot_created, r.name as robot_name, r.width as robot_width, r.height as robot_height,
	f.id as floor_id, f.uuid as floor_uuid, f.created as floor_created, f.name as floor_name, f.width as floor_width, f.grid as floor_grid
	FROM sessions s 
	LEFT OUTER JOIN robots r ON r.id = s.robotid 
  LEFT OUTER JOIN floors f ON f.id = s.floorid 
	WHERE s.uuid=$1`, sessionID.Value).Scan(
		&ID, &sessionUUID, &sessionCreated,
		&robotID, &robotUUID, &robotCreated, &robotName, &robotWidth, &robotHeight,
		&floorID, &floorUUID, &floorCreated, &floorName, &floorWidth, &floorGrid)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		robot := &Robot{id: robotID, ID: RobotID{robotUUID}, Created: robotCreated, Name: robotName, Width: robotWidth}
		floor := &Floor{id: floorID, ID: FloorID{floorUUID}, Created: floorCreated, Name: floorName, Width: floorWidth, Grid: floorGrid}
		return &Session{id: ID, ID: *sessionID, Created: sessionCreated, Robot: robot, Floor: floor}, nil
	}
}

// GetRobot gets a robot from the database
func (r *Repository) GetRobot(robotID *RobotID) (*Robot, error) {
	var ID, floorID int
	var robotWidth, robotHeight, floorWidth int
	var robotUUID, floorUUID, robotCreated, floorCreated, robotName, floorName string
	var floorGrid []int

	err := r.db.QueryRowContext(r.ctx, `SELECT 
	r.id as robot_id, r.uuid as robot_uuid, r.created as robot_created, r.name as robot_name, r.width as robot_width, r.height as robot_height,
	f.id as floor_id, f.uuid as floor_uuid, f.created as floor_created, f.name as floor_name, f.width as floor_width, f.grid as floor_grid
	FROM robots r 
  LEFT OUTER JOIN floors f ON f.id = r.floorid
	WHERE r.uuid=$1`, robotID.Value).Scan(
		&ID, &robotUUID, &robotCreated, &robotName, &robotWidth, &robotHeight,
		&floorID, &floorUUID, &floorCreated, &floorName, &floorWidth, &floorGrid)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		floor := &Floor{id: floorID, ID: FloorID{floorUUID}, Created: floorCreated, Name: floorName, Width: floorWidth, Grid: floorGrid}
		return &Robot{id: ID, ID: RobotID{robotUUID}, Created: robotCreated, Name: robotName, Width: robotWidth, Floor: floor}, nil
	}
}

// InsertPosition insert a posiiton update into the database
func (r *Repository) InsertPosition(session *Session, point *Point) error {
	stmt, err := r.db.Prepare("INSERT INTO points(sessionid, sequence, x, y, ) VALUES($1, $2, $3, $4)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, execErr := stmt.Exec(session.ID.Value, point.Sequence, point.X, point.Y)

	return execErr
}

// CreateRobot creates a new robot
func (r *Repository) CreateRobot(name string, width int, height int) (*RobotID, error) {

	id, err := uuid.NewUUID()
	log.Printf("Creating robot with ID: %v", id)

	if err != nil {
		return nil, err
	}

	stmt, err := r.db.Prepare("INSERT INTO robots(uuid, name, width, height) VALUES($1, $2, $3, $4)")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

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

	stmt, err := r.db.Prepare("INSERT INTO floors(uuid, name, width, grid) VALUES($1, $2, $3, $4)")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	array := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(*grid)), ", "), "[]")

	_, execErr := stmt.Exec(id.String(), name, width, "{"+array+"}")

	if execErr != nil {
		return nil, execErr
	}

	return &FloorID{Value: id.String()}, nil
}

// MapRobotToFloor maps a robot to a floor
func (r *Repository) MapRobotToFloor(robotID *RobotID, floorID *FloorID) error {
	stmt, err := r.db.Prepare("UPDATE Robots SET floorId = (SELECT id FROM floors where uuid = $1) WHERE uuid = $2")

	if err != nil {
		return err
	}

	_, execErr := stmt.Exec(robotID.Value, floorID.Value)

	return execErr
}
