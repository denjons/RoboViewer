package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	c "github.com/denjons/RoboViewer/common/kafka/consumer"
	client "github.com/denjons/RoboViewer/robot_service/client"
	db "github.com/denjons/RoboViewer/robot_service/database"
	kafka "github.com/denjons/RoboViewer/robot_service/kafka"
	service "github.com/denjons/RoboViewer/robot_service/service"
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

	parent := context.Background()
	repository, err := db.NewRepository(parent, *dbURL)
	if err != nil {
		log.Fatalf("Could not create repository: %v", err)
	}
	defer repository.Stop()
	sessionService := createSessionService(repository)
	startKafkaListener(sessionService)
	robotService := createRobotService(repository)

	log.Printf("Running database migrations")
	if err := db.MigrateDatabase(*dbMigration, *dbURL); err != nil {
		log.Fatalf("Could not migrate database schema: %v", err)
	}

	log.Printf("Creating listeners")

	http.HandleFunc("/floor/create", func(w http.ResponseWriter, r *http.Request) {
		floorDTO, err := parseFloor(r)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not parse json body: %v", err), 400)
		}

		log.Printf("received create floor request: %v", floorDTO)

		floorID, err := robotService.CreateFloor(floorDTO)

		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not create Floor: %v", err), 500)
		}

		jsonRespone, err := json.Marshal(floorID)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not parse Response: %v", err), 500)
		}

		w.Write(jsonRespone)
	})

	http.HandleFunc("/robot/create", func(w http.ResponseWriter, r *http.Request) {
		robotDTO, err := parseRobot(r)

		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not parse json body: %v", err), 400)
		}

		log.Printf("received create robot request: %v", robotDTO)

		robotID, err := robotService.CreateRobot(robotDTO)

		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not create Robot: %v", err), 500)
		}

		jsonRespone, err := json.Marshal(robotID)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not parse Response: %v", err), 500)
		}

		w.Write(jsonRespone)
	})

	http.HandleFunc("/floor/map/robot", func(w http.ResponseWriter, r *http.Request) {
		rfMap, err := parseRobotToFlooeMap(r)

		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not parse json body: %v", err), 400)
		}

		log.Printf("received map robot to floor request: %v", rfMap)

		mapErr := robotService.MapRobotToFloor(rfMap)

		if mapErr != nil {
			log.Printf("Error: %v", err)
			http.Error(w, fmt.Sprintf("Could not map robot to floor %v", mapErr), 400)
		}

		w.Write([]byte("OK"))
	})

	log.Printf("Serving on %v:%v", *serverHost, *serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", *serverHost, *serverPort), nil))
}

func startKafkaListener(service *service.SessionEventService) {
	topics := strings.Split(*kafkaTopic, ",")
	kafkaConsumer, err := c.NewKafkaConsumer(*kafkaHost, *kafkaGroupID, &topics)
	if err != nil {
		log.Fatalf("Error creating KafkaConsumer: %v", err)
	}
	listener, err := kafka.NewListener(kafkaConsumer)
	if err != nil {
		log.Fatalf("Error creating Listener: %v", err)
	}
	listener.ListenForPositionUpdateEvents(service.HandlePositionUpdateEvent)
}

func createSessionService(repository *db.Repository) *service.SessionEventService {
	sessionCache, err := db.NewSessionCache(repository)
	if err != nil {
		log.Fatalf("Could not create SessionCache: %v", err)
	}

	eventService, err := service.NewSessionEventService(sessionCache)
	if err != nil {
		log.Fatalf("Could not create SessionEventService: %v", err)
	}

	return eventService
}

func createRobotService(repository *db.Repository) *service.RobotService {
	service, err := service.NewRobotService(repository)
	if err != nil {
		log.Fatalf("Could not create RobotService: %v", err)
	}
	return service
}

func parseFloor(r *http.Request) (*client.FloorDTO, error) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	floorDTO := &client.FloorDTO{}
	json.Unmarshal(bytes, floorDTO)
	return floorDTO, nil
}

func parseRobot(r *http.Request) (*client.RobotDTO, error) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	robotDTO := &client.RobotDTO{}
	json.Unmarshal(bytes, robotDTO)

	return robotDTO, nil
}

func parseRobotToFlooeMap(r *http.Request) (*client.RobotToFloorMapDTO, error) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	rfMap := &client.RobotToFloorMapDTO{}
	json.Unmarshal(bytes, rfMap)

	return rfMap, nil
}
