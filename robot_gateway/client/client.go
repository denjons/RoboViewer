package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "github.com/denjons/RoboViewer/common/grpc/positionreport"
	rsClient "github.com/denjons/RoboViewer/robot_service/client"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

var (
	grpcHost = flag.String("grpc.host", "localhost:50001", "grpc bind host.")
	httpHost = flag.String("http.host", "http://localhost:8080", "http bind host")
)

func main() {

	robotID := createRobot("test robot", 4, 4)
	floorID := createFloor("test floor", 10, make([]int, 100))

	mapRobotToFloor(robotID, floorID)

	conn, err := grpc.Dial(*grpcHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPositionReportClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	stream, StreamErr := client.ReportPosition(ctx)

	if StreamErr != nil {
		log.Fatalf("Could not open stream: %v", StreamErr)
	}

	log.Println("sending position updates")

	sessionID := uuid.NewV1().String()
	for i := int32(1); i <= 1000; i++ {
		value := &pb.PositionUpdate{SequenceNumber: &pb.SequenceNumber{Value: 1}, Position: &pb.Position{X: i, Y: (i + 1)},
			RobotId: &pb.RobotId{Value: robotID.Value}, SessionId: &pb.SessionId{Value: sessionID}}
		if sendErr := stream.Send(value); sendErr != nil {
			log.Fatalf("Failed to send a update: %v", err)
		}
	}

	reply, closeErr := stream.CloseAndRecv()

	if closeErr != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, closeErr, nil)
	} else {
		log.Println("sending complete")
	}

	log.Printf("Report Responde: %v", reply)
}

func createRobot(name string, width int, height int) *rsClient.RobotID {
	robotDTO := &rsClient.RobotDTO{Name: name, Width: width, Height: height}

	robotJSON, err := json.Marshal(robotDTO)

	if err != nil {
		log.Fatalf("Could not parse robot request %v", err)
	}

	bytes := post(&robotJSON, "/robot/create")

	robotID := &rsClient.RobotID{}
	json.Unmarshal(bytes, robotID)

	return robotID
}

func mapRobotToFloor(robotID *rsClient.RobotID, floorID *rsClient.FloorID) {
	rfMap := &rsClient.RobotToFloorMapDTO{RobotID: *robotID, FloorID: *floorID}

	rfMapJSON, err := json.Marshal(rfMap)

	if err != nil {
		log.Fatalf("Could not parse robot request %v", err)
	}

	post(&rfMapJSON, "/floor/map/robot")
}

func createFloor(name string, width int, grid []int) *rsClient.FloorID {
	floorDTO := &rsClient.FloorDTO{Name: name, Width: width, Grid: grid}

	floorJSON, err := json.Marshal(floorDTO)

	if err != nil {
		log.Fatalf("Could not parse robot request %v", err)
	}

	bytes := post(&floorJSON, "/floor/create")

	floorID := &rsClient.FloorID{}
	json.Unmarshal(bytes, floorID)

	return floorID
}

func post(request *[]byte, path string) []byte {
	response, err := http.Post(*httpHost+path, "application/json", bytes.NewBuffer(*request))

	if err != nil {
		log.Fatalf("Could not send robot request %v", err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("Got response %v", response.StatusCode)
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Could read response bytes %v", err)
	}
	return bytes
}
