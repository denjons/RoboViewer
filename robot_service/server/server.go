package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	c "github.com/denjons/RoboViewer/common/kafka/consumer"
	common "github.com/denjons/RoboViewer/common/model"
	"github.com/denjons/RoboViewer/robot_service/client/model"
	"github.com/denjons/RoboViewer/robot_service/database"
	kafka "github.com/denjons/RoboViewer/robot_service/kafka"
)

var (
	serverPort   = flag.Int("server.port", 8080, "Server bind port.")
	serverHost   = flag.String("server.host", "localhost", "server bind host.")
	dbURL        = flag.String("db.url", "postgres://robot_service:robot_service@localhost:5432/robot_service?sslmode=disable", "database url")
	dbMigration  = flag.String("db.migration", "file://robot_service/database/migrations", "db migration file location")
	kafkaHost    = flag.String("kafka.host", "localhost", "broker host")
	kafkaTopic   = flag.String("kafka.topic", "position-events", "Topic for postion events")
	kafkaGroupID = flag.String("kafka.groupid", "default", "Kafka group id")
)

func main() {

	startKafkaListener()

	log.Printf("Running database migrations")
	if err := database.MigrateDatabase(*dbMigration, *dbURL); err != nil {
		log.Fatalf("Could not migrate database schema: %v", err)
	}

	log.Printf("Creating listeners")
	http.HandleFunc("/floor/create", func(w http.ResponseWriter, r *http.Request) {
		floorDTO, err := parseFloor(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not parse json body: %v", err), 400)
		}
		fmt.Printf("Got floor: %v", floorDTO)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/robot/create", func(w http.ResponseWriter, r *http.Request) {
		robotDTO, err := parseRobot(r)

		if err != nil {
			http.Error(w, fmt.Sprintf("Could not parse json body: %v", err), 400)
		}
		fmt.Printf("Got robot: %v", robotDTO)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/floor/set/robot", func(w http.ResponseWriter, r *http.Request) {

		robotIds, ok := r.URL.Query()["robot_id"]
		if !ok || len(robotIds[0]) < 1 {
			http.Error(w, "missing robot_id in request", 400)
		}

		floorIds, ok := r.URL.Query()["floor_id"]
		if !ok || len(floorIds[0]) < 1 {
			http.Error(w, "missing floor_id in request", 400)
		}

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Printf("Serving on %v:%v", *serverHost, *serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", *serverHost, *serverPort), nil))
}

func startKafkaListener() {
	topics := strings.Split(*kafkaTopic, ",")
	kafkaConsumer, err := c.NewKafkaConsumer(*kafkaHost, *kafkaGroupID, &topics)
	if err != nil {
		log.Fatalf("Error creating KafkaConsumer: %v", err)
	}
	listener, err := kafka.NewListener(kafkaConsumer)
	if err != nil {
		log.Fatalf("Error creating Listener: %v", err)
	}
	listener.ListenForPositionUpdateEvents(func(p *common.PositionUpdateEvent) {
		log.Printf("Received message: %v", p)
	})
}

func getBytes(r *http.Request) (*[]byte, error) {
	bytes := make([]byte, r.ContentLength)
	size, err := r.Body.Read(bytes)
	if err != nil {
		return nil, err
	}
	if size < len(bytes) {
		return nil, fmt.Errorf("Could only read %v out of %v bytes", size, r.ContentLength)
	}
	return &bytes, nil
}

func parseFloor(r *http.Request) (*model.FloorDTO, error) {
	bytes, err := getBytes(r)
	if err != nil {
		return nil, err
	}
	floorDTO := model.FloorDTO{}
	json.Unmarshal(*bytes, &floorDTO)
	return &floorDTO, nil
}

func parseRobot(r *http.Request) (*model.RobotDTO, error) {
	bytes, err := getBytes(r)
	if err != nil {
		return nil, err
	}
	robotDTO := model.RobotDTO{}
	json.Unmarshal(*bytes, &robotDTO)

	return &robotDTO, nil
}
