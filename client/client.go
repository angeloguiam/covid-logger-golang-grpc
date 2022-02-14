package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/logger/loggerpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client starting...")
	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error creating client connection: %v", err)
	}

	defer clientConn.Close()

	dataLoggerServiceClient := loggerpb.NewDataLoggerServiceClient(clientConn)

	logData(dataLoggerServiceClient)
	fmt.Println(dataLoggerServiceClient)
}

func logData(loggerClient loggerpb.DataLoggerServiceClient) {
	fmt.Println("Function log is invoked.")

	req := &loggerpb.DataLoggerRequest{
		Timestamp: time.Now().String(),
		Source:    "PGH",
		Report: &loggerpb.Report{
			Cases:     10,
			Death:     10,
			Recovered: 5,
		},
	}

	res, err := loggerClient.LogData(context.Background(), req)

	if err != nil {
		log.Fatalf("Error in response: %v\n", err)
	}

	fmt.Println(res.GetStatus())
}
