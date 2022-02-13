package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/logger/loggerpb"
)

type logger struct{}

func (l *logger) Log(ctx context.Context, req *loggerpb.DataLoggerRequest) (*loggerpb.DataLoggerResponse, error) {
	fmt.Println("Log function is invoked")

	date := req.GetTimestamp()
	source := req.GetSource()
	report := req.GetReport()

	log.Printf("%v: Report from %v: %v", date, source, report)

	res := &loggerpb.DataLoggerResponse{
		Summary: &loggerpb.Report{
			Cases:     10,
			Death:     10,
			Recovered: 10,
		},
		Status: loggerpb.Status_OK,
	}

	return res, nil
}

func main() {
	fmt.Println("Satarting server...")
	address := "0.0.0.0:50051"

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listen to %s: %v", address, err)
	}

	server := grpc.NewServer()
	loggerpb.RegisterDataLoggerServiceServer(server, &logger{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on %s: %v", address, err)
	}

}
