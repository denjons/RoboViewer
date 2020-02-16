package producer

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer puts an event to a Kafka topic
type Producer interface {
	Put(events chan Event) error
}

// Event is generic type for events publidhed on a Kafka topic
type Event interface {
	ToJSON() []byte
}

// JSONProducer for the Kafka broker
type JSONProducer struct {
	Topic  string
	Broker string
}

// Put events on a Kafla topic through the given channel
func (jp *JSONProducer) Put(eventChannel chan Event) error {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": jp.Broker})

	if err != nil {
		return err
	}

	for value := range eventChannel {
		doneChan := make(chan bool)
		deferCloseProducer(producer, doneChan)
		producer.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &jp.Topic, Partition: kafka.PartitionAny}, Value: value.ToJSON()}
		_ = <-doneChan
	}

	producer.Close()

	return nil
}

func deferCloseProducer(producer *kafka.Producer, doneChan chan bool) {
	go func() {
		defer close(doneChan)
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				log.Printf("Ignored event: %s\n", ev)
			}
		}
	}()
}
