package eventhandler

import (
	"encoding/json"
	"log"

	"github.com/denjons/RoboViewer/common/kafka/producer"
	ps "github.com/denjons/RoboViewer/common/model/position"
)

// EventHandler handles incoming events
type EventHandler interface {
	HandleEvent(event ps.PositionUpdateEvent)
	Close()
}

// KafkaEventHandler handles events and forwards them to the correct destination
type KafkaEventHandler struct {
	Producer        *producer.KafkaProducer
	producerChannel chan *[]byte
}

// NewEventHandler creates a new event eventhandler
func NewEventHandler(eventProducer *producer.KafkaProducer) *KafkaEventHandler {
	handler := &KafkaEventHandler{}
	handler.producerChannel = make(chan *[]byte)
	handler.Producer = eventProducer
	go func() {
		err := handler.Producer.Listen(handler.producerChannel)
		if err != nil {
			log.Fatalf("Could not start Kafka producer: %v", err)
		}
	}()
	return handler
}

// HandleEvent handles the event and forwards it to it's next destination
func (handler *KafkaEventHandler) HandleEvent(event *ps.PositionUpdateEvent) error {
	data, err := json.Marshal(*event)

	if err != nil {
		return err
	}

	handler.producerChannel <- &data

	return nil
}

// Close any connections or channels which the event handler may be using
func (handler *KafkaEventHandler) Close() {
	close(handler.producerChannel)
}
