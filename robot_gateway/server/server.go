package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/denjons/RoboViewer/common/grpc/positionreport"
	"google.golang.org/grpc"
)

var (
	serverPort = flag.Int("server.port", 50001, "Server bind port.")
	serverHost = flag.String("server.host", "localhost", "server bind host.")
)

type server struct {
	pb.UnimplementedPositionReportServer
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

		log.Printf("Received position X: %v, Y: %v", &positionUpdate.Position.X, &positionUpdate.Position.Y)
	}
}

func main() {
	flag.Parse()

	var hostPort = fmt.Sprintf("%s:%d", *serverHost, *serverPort)

	log.Printf("Serving from: %v", hostPort)

	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("Failed to bind server: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPositionReportServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
