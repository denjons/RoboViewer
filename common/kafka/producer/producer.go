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

// KafkaProducer for the Kafka broker
type KafkaProducer struct {
	Topic    string
	Broker   string
	Producer *kafka.Producer
}

// NewKafkaProducer creates a new KafkaProducer
func NewKafkaProducer(broker, topic string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		return nil, err
	}

	p := &KafkaProducer{}
	p.Producer = producer
	p.Broker = broker
	p.Topic = topic

	return p, nil
}

// Listen to events from given channel and put them on Kafka topic
func (kp *KafkaProducer) Listen(eventChannel chan *[]byte) error {

	for value := range eventChannel {
		doneChan := make(chan bool)
		deferCloseProducer(kp.Producer, doneChan)
		kp.Producer.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &kp.Topic, Partition: kafka.PartitionAny}, Value: *value}
		_ = <-doneChan
	}

	kp.Producer.Close()

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
