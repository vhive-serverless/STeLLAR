package main

import (
	"flag"
	"fmt"
	"io"
	"lambda-benchmarking/client/benchmarking"
	"lambda-benchmarking/client/visualization"
	"log"
	"os"
	"path/filepath"
	"time"
)

//Note: those variables are pointers
var requestsFlag = flag.Int("requests", 1, "Number of outstanding requests for this run.")
var payloadLengthFlag = flag.Int("payloadLengthBytes", 8, "Length of the payload generated by the lambda function.")
var outputPathFlag = flag.String("outputPath", "latency-samples", "The path where latency samples should be written.")
var frequencySecondsFlag = flag.Int("frequencySeconds", -1, "Frequency at which the latency profiler operates.")
var randomizedFlag = flag.Bool("randomized", false, "Sample deltas from a Gaussian with shifted mean (not scaled, stddev still 1).")
var burstsNumberFlag = flag.Int("bursts", 5, "Number of bursts which the latency profiler will trigger.")
var lambdaIncrementLimitFlag = flag.Int("lambdaIncrementLimit", 5e7, "Increment limit for the lambda function to busy spin on.")
var gatewayEndpointFlag = flag.String("gatewayEndpoint", "", "The API Endpoint to make requests to.")
var visualizationFlag = flag.String("visualization", "", "The type of visualization to create (histogram, cdf).")

func main() {
	flag.Parse()

	outputDirectoryPath := filepath.Join(*outputPathFlag, time.Now().Format(time.RFC850))
	log.Printf("Creating working directory at %s", outputDirectoryPath)
	if err := os.Mkdir(outputDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(filepath.Join(outputDirectoryPath, "run_logs.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	stdoutFileMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(stdoutFileMultiWriter)

	log.Printf("Started benchmarking HTTP client on %v.", time.Now().Format(time.RFC850))
	log.Printf("Parameters entered: %d requests in a burst, %dbytes payload length, %d busy spin counter, %d profiler run frequency, output path was set to `%s`.",
		*requestsFlag, *payloadLengthFlag, *lambdaIncrementLimitFlag, *frequencySecondsFlag, *outputPathFlag)

	csvFile, err := os.Create(filepath.Join(outputDirectoryPath, fmt.Sprintf(
		"%dbursts_%dreqs_freq%ds_payload%db_counter%d.csv",
		*burstsNumberFlag,
		*requestsFlag,
		*frequencySecondsFlag,
		*payloadLengthFlag,
		*lambdaIncrementLimitFlag)))
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	burstDeltas := benchmarking.CreateBurstDeltas(*frequencySecondsFlag, *burstsNumberFlag, *randomizedFlag)
	relativeBurstDeltas := benchmarking.MakeBurstDeltasRelative(burstDeltas)

	log.Println("Running profiler...")
	benchmarking.SafeWriterInstance.Initialize(csvFile)
	benchmarking.TriggerRelativeAsyncBurstGroups(*gatewayEndpointFlag, relativeBurstDeltas, *requestsFlag, *lambdaIncrementLimitFlag, *payloadLengthFlag)

	log.Println("Flushing results to CSV file...")
	benchmarking.SafeWriterInstance.Writer.Flush()

	if *visualizationFlag == "" {
		log.Println("Skipping visualization...")
	} else {
		log.Printf("Creating %ss from CSV file `%s`", *visualizationFlag, csvFile.Name())
		visualization.GenerateVisualization(
			*visualizationFlag,
			*burstsNumberFlag,
			burstDeltas,
			relativeBurstDeltas,
			csvFile,
			outputDirectoryPath,
		)
	}

	log.Println("Exiting...")
}
