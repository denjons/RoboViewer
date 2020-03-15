package consumer

import (
	"log"
	"os"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Consumer consumer messages from configured topic
type Consumer interface {
	Start(stopChannel, consumerChannel chan *[]byte)
}

// KafkaConsumer implementation
type KafkaConsumer struct {
	Topics   *[]string
	Broker   string
	consumer *kafka.Consumer
}

// NewKafkaConsumer creates a new consumer
func NewKafkaConsumer(broker *string, group *string, topics *[]string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		// Enable generation of PartitionEOF when the
		// end of a partition is reached.
		"enable.partition.eof": true,
		"auto.offset.reset":    "earliest"})

	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{topics, *broker, c}, nil
}

// Start consuming on configured topics
func (kc *KafkaConsumer) Start(consumerChannel chan *[]byte) {

	err := kc.consumer.SubscribeTopics(*kc.Topics, nil)

	log.Fatalf("Error starting consumer %v", err)

	sigchan := make(chan os.Signal, 1)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating", sig)
			run = false

		case ev := <-kc.consumer.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Printf("Assigned partition %v", e)
				kc.consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				log.Printf("Revoked partitions %v", e)
				kc.consumer.Unassign()
			case *kafka.Message:
				log.Printf(" Message on %s", e.TopicPartition)
				consumerChannel <- &e.Value
			case kafka.PartitionEOF:
				log.Printf("Oartition EOF: %v", e)
			case kafka.Error:
				log.Printf("Kafka Consumer Error: %v", e)
			}
		}
	}

	log.Printf("Closing consumer")
	kc.consumer.Close()
}
