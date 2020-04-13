package eventhandler

import (
	"encoding/json"

	kafka "github.com/denjons/RoboViewer/common/kafka/producer"
	client "github.com/denjons/RoboViewer/robot_gateway/client"
	grpcClient "github.com/denjons/RoboViewer/robot_gateway/client/grpc/positionreport"
	kafkaClient "github.com/denjons/RoboViewer/robot_gateway/client/kafka"
)

// EventHandler handles incoming events
type EventHandler interface {
	HandleEvent(event kafkaClient.PositionUpdateEvent)
	Close()
}

// KafkaEventHandler handles events and forwards them to the correct destination
type KafkaEventHandler struct {
	producer *kafka.KafkaProducer
}

// NewEventHandler creates a new event eventhandler
func NewEventHandler(kafkaProducer *kafka.KafkaProducer) *KafkaEventHandler {
	return &KafkaEventHandler{producer: kafkaProducer}
}

// HandlePositionUpdate handles the event and forwards it to it's next destination
func (handler *KafkaEventHandler) HandlePositionUpdate(positionUpdate *grpcClient.PositionUpdate) error {

	positionUpdateEvent := client.ConvertToPositionUpdateEvent(positionUpdate)

	data, err := json.Marshal(*positionUpdateEvent)

	if err != nil {
		return err
	}

	handler.producer.Publish(&kafka.KafkaMessage{Topic: "position-events", Payload: data})

	return nil
}
