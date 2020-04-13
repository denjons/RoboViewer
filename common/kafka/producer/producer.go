package producer

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer puts an event to a Kafka topic
type Producer interface {
	Listen(events chan *[]byte) error
	Start() error
}

// KafkaMessage a message to be sent a Kafka topic
type KafkaMessage struct {
	Topic   string
	Payload []byte
}

// KafkaProducer for the Kafka broker
type KafkaProducer struct {
	Broker   string
	Producer *kafka.Producer
}

// NewKafkaProducer creates a new KafkaProducer
func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		return nil, err
	}

	return &KafkaProducer{Broker: broker, Producer: producer}, nil
}

// Publish events to a kafka topic
func (kp *KafkaProducer) Publish(kafkaMessage *KafkaMessage) {
	doneChan := make(chan bool)
	deferCloseProducer(kp.Producer, doneChan)
	kp.Producer.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &kafkaMessage.Topic, Partition: kafka.PartitionAny}, Value: kafkaMessage.Payload}
	_ = <-doneChan
}

// Close the producer
func (kp *KafkaProducer) Close() {
	kp.Producer.Close()
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
