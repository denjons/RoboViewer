package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/denjons/RoboViewer/common/grpc/positionreport"
	pr "github.com/denjons/RoboViewer/common/kafka/producer"
	model "github.com/denjons/RoboViewer/common/model/position"
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
			return stream.SendAndClose(&pb.PeportResponse{
				ResponseStatus: pb.ResponseStatus_OK,
				StatusMessage:  "Received position updates",
			})
		}
		if err != nil {
			return err
		}

		log.Printf("Received position X: %v, Y: %v", positionUpdate.Position.X, positionUpdate.Position.Y)

		evErr := s.Handler.HandleEvent(parse(positionUpdate))
		if evErr != nil {
			return evErr
		}
	}
}

func parse(positionUpdate *pb.PositionUpdate) *model.PositionUpdateEvent {
	return &model.PositionUpdateEvent{
		Sequence:  positionUpdate.SequenceNumber.Count,
		RobotID:   positionUpdate.RobotId.Id,
		SessionID: positionUpdate.SessionId.Id,
		Position:  &model.Position{X: positionUpdate.Position.X, Y: positionUpdate.Position.Y},
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
