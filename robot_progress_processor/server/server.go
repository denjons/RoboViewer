package main

import (
	"flag"
	"log"
	"strings"

	c "github.com/denjons/RoboViewer/common/kafka/consumer"
)

var (
	kafkaHost    = flag.String("kafka.host", "localhost", "broker host")
	kafkaTopic   = flag.String("kafka.topic", "position-events", "Topic for postion events")
	kafkaGroupID = flag.String("kafka.groupid", "default", "Kafka group id")
)

func main() {
	log.Print("Starting robot progress processor")
	flag.Parse()

	topics := strings.Split(*kafkaTopic, ",")
	kafkaConsumer, err := c.NewKafkaConsumer(kafkaHost, kafkaGroupID, &topics)

	if err != nil {
		log.Fatalf("Error %v", err)
	}

	channel := make(chan []byte)

	go printMessages(channel)

	kafkaConsumer.Start(channel)

}

func printMessages(channel chan []byte) {
	for msg := range channel {
		log.Printf("Received message: %v", string(msg))
	}
}
