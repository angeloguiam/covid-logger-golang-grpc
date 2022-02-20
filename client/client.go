package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/logger/loggerpb"

	"google.golang.org/grpc"
)

type record struct {
	timestamp string
	source    string
	confirmed int
	recovered int
	death     int
}

func main() {
	const fileInput = "data.csv"

	fmt.Println("Client starting...")
	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error creating client connection: %v", err)
	}

	defer clientConn.Close()

	records, err := loadDataFromCSV(fileInput)

	if err != nil {
		log.Fatal("Unable to read input file CSV file")
	}

	dataLoggerServiceClient := loggerpb.NewDataLoggerServiceClient(clientConn)

	data := []record{}

	for i, row := range records {
		if i == 0 {
			continue
		}

		newRecord := record{
			timestamp: row[0],
			source:    row[1],
			confirmed: stringToInt(row[2]),
			recovered: stringToInt(row[3]),
			death:     stringToInt(row[4]),
		}

		data = append(data, newRecord)
		logData(dataLoggerServiceClient, newRecord)

	}
}

func stringToInt(value string) int {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("Cannot convert string to int: %v, %v", intValue, err)
	}

	return int(intValue)
}

func loadDataFromCSV(filename string) ([][]string, error) {

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal("Unable to read input file from %v", filename)
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal("Unable to parse file as CSV from %v", filename)
		return nil, err
	}

	return records, nil
}

func logData(loggerClient loggerpb.DataLoggerServiceClient, newRecord record) {
	fmt.Println("Function log is invoked.")

	req := &loggerpb.DataLoggerRequest{
		Timestamp: newRecord.timestamp,
		Source:    newRecord.source,
		Report: &loggerpb.Report{
			Confirmed: int32(newRecord.confirmed),
			Recovered: int32(newRecord.recovered),
			Death:     int32(newRecord.death),
		},
	}

	res, err := loggerClient.LogData(context.Background(), req)

	if err != nil {
		log.Fatalf("Error in response: %v\n", err)
	}

	fmt.Println(res.GetStatus())
}
