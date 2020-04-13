package service

import (
	"errors"
	"log"

	kafka "github.com/denjons/RoboViewer/robot_gateway/client/kafka"
	db "github.com/denjons/RoboViewer/robot_service/database"
)

// SessionEventService handles all session events
type SessionEventService struct {
	sessionCache *db.SessionCache
}

// NewSessionEventService creates a new SessionCache
func NewSessionEventService(sessionCache *db.SessionCache) (*SessionEventService, error) {
	if sessionCache == nil {
		return nil, errors.New("SessionCache cannot be nil")
	}
	return &SessionEventService{sessionCache: sessionCache}, nil
}

// HandlePositionUpdateEvent handles incoming position update evenets
func (s *SessionEventService) HandlePositionUpdateEvent(event *kafka.PositionUpdateEvent) {
	sessionID := &db.SessionID{Value: event.SessionID}
	robotID := &db.RobotID{Value: event.RobotID}
	point := &db.Point{Sequence: event.Sequence, X: event.Position.X, Y: event.Position.Y}
	err := s.sessionCache.UpdateSession(sessionID, robotID, point)

	if err != nil {
		log.Fatalf("Error handling position event %v, message %v", err, *event)
	}

}
