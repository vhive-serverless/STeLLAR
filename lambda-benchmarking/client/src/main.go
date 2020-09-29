package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"lambda-benchmarking/client/benchmarking"
	"log"
	"os"
	"time"
)

//Note: those variables are pointers
var requestsFlag = flag.Int("requests", 1, "The number of outstanding requests.")
var execMillisecondsFlag = flag.Int("execMilliseconds", 80, "The number of milliseconds for the lambda function to busy spin.")
var outputPathFlag = flag.String("outputPath", "latency-samples", "The path where latency samples should be written for this run.")

func main() {
	flag.Parse()
	log.Printf("Started benchmarking HTTP client started with %d outstanding requests and %dms busy spin. Output path was set to: %s",
		*requestsFlag, *execMillisecondsFlag, *outputPathFlag)

	file, err := os.Create(fmt.Sprintf("%s/%s.csv", *outputPathFlag, time.Now().Format(time.RFC850)))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create writer, make it safe for concurrent use and flush it when program finishes
	writer := csv.NewWriter(file)
	defer writer.Flush()
	safeWriter := benchmarking.SafeWriter{Writer: writer}
	safeWriter.WriteRowToFile("AwsRequestID", "Client latency (ms)")

	// Generate as many requests (and thus latency records) as specified in the arguments
	reqChannel := make(chan int)
	for i := 0; i < *requestsFlag; i++ {
		go safeWriter.GenerateLatencyRecord(reqChannel, i, *execMillisecondsFlag)
	}

	// Wait until all responses have been received before returning and thus flushing records to disk
	for *requestsFlag > 0 {
		requestFinished := <-reqChannel
		*requestsFlag--
		log.Printf("Received response for request %d, %d responses remaining\n", requestFinished, *requestsFlag)
	}
}
