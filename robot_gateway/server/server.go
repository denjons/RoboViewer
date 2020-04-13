package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pr "github.com/denjons/RoboViewer/common/kafka/producer"
	rgClient "github.com/denjons/RoboViewer/robot_gateway/client"
	pb "github.com/denjons/RoboViewer/robot_gateway/client/grpc/positionreport"
	ev "github.com/denjons/RoboViewer/robot_gateway/event"
	"google.golang.org/grpc"
)

var (
	serverPort = flag.Int("server.port", 50001, "Server bind port.")
	serverHost = flag.String("server.host", "localhost", "server bind host.")
	kafkaHost  = flag.String("kafka.host", "localhost", "broker host")
	kafkaTopic = flag.String("kafka.topic", "position-events", "Topic for postion events")
)

type server struct {
	pb.UnimplementedPositionReportServer
	Handler *ev.KafkaEventHandler
}

func (s *server) ReportPosition(stream pb.PositionReport_ReportPositionServer) error {
	for {
		positionUpdate, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.PositionUpdateResponse{
				ResponseStatus: pb.ResponseStatus_OK,
				StatusMessage:  "Received position updates",
			})
		}
		if err != nil {
			return err
		}

		log.Printf("Received position X: %v, Y: %v", positionUpdate.Position.X, positionUpdate.Position.Y)

		evErr := s.Handler.HandleEvent(rgClient.ConvertToPositionUpdateEvent(positionUpdate))
		if evErr != nil {
			return evErr
		}
	}
}

func createEventHandler() *ev.KafkaEventHandler {
	kafkaProducer, err := pr.NewKafkaProducer(*kafkaHost, *kafkaTopic)

	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}

	handler := ev.NewEventHandler(kafkaProducer)

	return handler
}

func main() {
	flag.Parse()

	var hostPort = fmt.Sprintf("%s:%d", *serverHost, *serverPort)

	log.Printf("Serving from: %v", hostPort)

	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("Failed to bind server: %v", err)
	}

	handler := createEventHandler()

	s := grpc.NewServer()
	pb.RegisterPositionReportServer(s, &server{Handler: handler})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
