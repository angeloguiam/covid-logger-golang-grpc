package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/logger/loggerpb"

	"google.golang.org/grpc"
)

type logger struct{}

type record struct {
	timestamp string
	source    string
	confirmed int
	recovered int
	death     int
}

func (l *logger) LogData(ctx context.Context, req *loggerpb.DataLoggerRequest) (*loggerpb.DataLoggerResponse, error) {
	const fileOutput = "record_from_client.csv"

	fmt.Println("Log function is invoked")

	newRecord := record{
		timestamp: req.GetTimestamp(),
		source:    req.GetSource(),
		confirmed: int(req.GetReport().GetConfirmed()),
		recovered: int(req.GetReport().GetRecovered()),
		death:     int(req.GetReport().GetDeath()),
	}

	saveRecordToFile(fileOutput, newRecord)

	res := &loggerpb.DataLoggerResponse{
		Summary: &loggerpb.Report{
			Confirmed: 10,
			Death:     10,
			Recovered: 10,
		},
		Status: loggerpb.Status_OK,
	}

	return res, nil
}

func saveRecordToFile(filename string, newRecord record) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalf("Error opening file: %v", filename)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	row := []string{newRecord.timestamp, newRecord.source, strconv.Itoa(newRecord.confirmed), strconv.Itoa(newRecord.recovered), strconv.Itoa(newRecord.death)}
	writer.Write(row)

	fmt.Println(row)
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
