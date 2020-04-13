package database

import (
	"errors"
	"fmt"
	"log"
)

// SessionCache handles high level database operations on Robot Sessions.
type SessionCache struct {
	repository Repository
	cache      map[string]Session
}

// NewSessionCache return a new SessionCache
func NewSessionCache(r *Repository) (*SessionCache, error) {
	if r == nil {
		return nil, errors.New("repository cannot be nil")
	}
	return &SessionCache{repository: *r, cache: make(map[string]Session)}, nil
}

/*
UpdateSession updates a session with a new position of a robot.
If the session does not exist from before a new one will be created.
*/
func (sc *SessionCache) UpdateSession(sessionID *SessionID, robotID *RobotID, point *Point) error {
	session, err := sc.getSession(sessionID, robotID)
	if err != nil {
		return err
	}
	return sc.repository.InsertPosition(session, point)
}

func (sc *SessionCache) getSession(sessionID *SessionID, robotID *RobotID) (*Session, error) {
	session, err := sc.getExistingSession(sessionID)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	log.Printf("Could not find existing session for ID %v", sessionID.Value)

	createErr := sc.createSessionForRobot(sessionID, robotID)
	if createErr != nil {
		return nil, createErr
	}

	log.Printf("Created new session with ID %v and robot %v", sessionID.Value, robotID.Value)

	createdSession, err := sc.getExistingSession(sessionID)
	if err != nil {
		return nil, err
	}
	if createdSession == nil {
		return nil, fmt.Errorf("Could not get created session with id '%v'", sessionID.Value)
	}
	return createdSession, nil

}

func (sc *SessionCache) getExistingSession(sessionID *SessionID) (*Session, error) {
	cachedSession, exists := sc.cache[sessionID.Value]
	if exists {
		return &cachedSession, nil
	}

	session, err := sc.repository.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	sc.cache[sessionID.Value] = *session

	return session, nil
}

func (sc *SessionCache) createSessionForRobot(sessionID *SessionID, robotID *RobotID) error {
	robot, err := sc.repository.GetRobot(robotID)
	if err != nil {
		return err
	}
	if robot == nil {
		return fmt.Errorf("Robot with ID '%v' does not exists", robotID.Value)
	}
	return sc.repository.CreateSession(sessionID, robot, robot.Floor)
}
