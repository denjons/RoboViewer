package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pr "github.com/denjons/RoboViewer/common/kafka/producer"
	"github.com/denjons/RoboViewer/robot_gateway/client/grpc/positionreport"
	pb "github.com/denjons/RoboViewer/robot_gateway/client/grpc/positionreport"
	ev "github.com/denjons/RoboViewer/robot_gateway/event"
	"google.golang.org/grpc"
)

var (
	serverPort = flag.Int("server.port", 50001, "Server bind port.")
	serverHost = flag.String("server.host", "localhost", "server bind host.")
	kafkaHost  = flag.String("kafka.host", "localhost", "broker host")
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

		evErr := s.Handler.HandlePositionUpdate(positionUpdate)
		if evErr != nil {
			return evErr
		}
	}
}

func (s *server) ReportSession(context.Context, *positionreport.SessionUpdate) (*positionreport.SessionUpdateResponse, error) {
	return nil, nil
}

func main() {
	flag.Parse()

	kafkaProducer, err := pr.NewKafkaProducer(*kafkaHost)

	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}

	defer kafkaProducer.Close()

	handler := ev.NewEventHandler(kafkaProducer)

	var hostPort = fmt.Sprintf("%s:%d", *serverHost, *serverPort)

	log.Printf("Serving from: %v", hostPort)

	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("Failed to bind server: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPositionReportServer(s, &server{Handler: handler})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
