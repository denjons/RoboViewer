package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/denjons/RoboViewer/common/grpc/positionreport"
	"google.golang.org/grpc"
)

var (
	serverHost = flag.String("server.host", "localhost:50001", "server bind host.")
	updates    = []*pb.PositionUpdate{
		{SequenceNumber: &pb.SequenceNumber{Count: 1}, Position: &pb.Position{X: 1, Y: 1}, RobotId: &pb.RobotId{Id: 1}, SessionId: &pb.SessionId{Id: 1}},
		{SequenceNumber: &pb.SequenceNumber{Count: 2}, Position: &pb.Position{X: 2, Y: 1}, RobotId: &pb.RobotId{Id: 1}, SessionId: &pb.SessionId{Id: 1}},
		{SequenceNumber: &pb.SequenceNumber{Count: 3}, Position: &pb.Position{X: 2, Y: 2}, RobotId: &pb.RobotId{Id: 1}, SessionId: &pb.SessionId{Id: 1}},
		{SequenceNumber: &pb.SequenceNumber{Count: 4}, Position: &pb.Position{X: 3, Y: 2}, RobotId: &pb.RobotId{Id: 1}, SessionId: &pb.SessionId{Id: 1}},
		{SequenceNumber: &pb.SequenceNumber{Count: 5}, Position: &pb.Position{X: 4, Y: 2}, RobotId: &pb.RobotId{Id: 1}, SessionId: &pb.SessionId{Id: 1}},
	}
)

func main() {

	conn, err := grpc.Dial(*serverHost, grpc.WithInsecure(), grpc.WithBlock())
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

	for i := int32(1); i <= 1000; i++ {
		value := &pb.PositionUpdate{SequenceNumber: &pb.SequenceNumber{Count: 1}, Position: &pb.Position{X: i, Y: (i + 1)}, RobotId: &pb.RobotId{Id: 1}, SessionId: &pb.SessionId{Id: 1}}
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
