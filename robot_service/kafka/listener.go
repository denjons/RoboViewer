package kafka

import (
	"encoding/json"
	"errors"

	c "github.com/denjons/RoboViewer/common/kafka/consumer"
	client "github.com/denjons/RoboViewer/robot_gateway/client/kafka"
)

// Listener listens on kafka topic
type Listener struct {
	KafkaConsumer *c.KafkaConsumer
}

// NewListener creates a new listener
func NewListener(kafkaConsumer *c.KafkaConsumer) (*Listener, error) {
	if kafkaConsumer == nil {
		return nil, errors.New("kafkaConsumer is nil")
	}
	return &Listener{kafkaConsumer}, nil
}

// ListenForPositionUpdateEvents listens to events from the position update topic
func (l *Listener) ListenForPositionUpdateEvents(handler func(p *client.PositionUpdateEvent)) error {
	channel := make(chan []byte)

	go handleMessage(channel, handler)

	go l.KafkaConsumer.Start(channel)

	return nil
}

func handleMessage(channel chan []byte, handler func(p *client.PositionUpdateEvent)) {
	for msg := range channel {
		positionUpdateEvent := client.PositionUpdateEvent{}
		json.Unmarshal(msg, &positionUpdateEvent)
		handler(&positionUpdateEvent)
	}
}
