package main

import (
	"flag"
	"github.com/go-gota/gota/dataframe"
	"io"
	"lambda-benchmarking/client/experiment"
	"lambda-benchmarking/client/experiment/configuration"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//Note: those variables are pointers
var randomizationFlag = flag.Bool("randomized", true, "If true, sample deltas from a scaled and shifted standard normal distribution.")
var visualizationFlag = flag.String("visualization", "CDF", "The type of visualization to create (per-burst histogram \"histogram\" or empirical CDF \"CDF\").")
var outputPathFlag = flag.String("outputPath", "latency-samples", "The path where latency samples should be written.")
var configPathFlag = flag.String("configPath", "config.csv", "Configuration file with details of experiments.")
var gatewaysPathFlag = flag.String("gatewaysPath", "gateways.csv", "File containing ids of gateways to be used.")

func main() {
	//rand.Seed(time.Now().Unix()) not sure if we want irreproducible deltas?
	flag.Parse()

	outputDirectoryPath := filepath.Join(*outputPathFlag, time.Now().Format(time.RFC850))
	log.Printf("Creating directory for this run at `%s`", outputDirectoryPath)
	if err := os.Mkdir(outputDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFile := setupClientLogging(outputDirectoryPath)
	defer logFile.Close()

	log.Printf("Started benchmarking HTTP client on %v.", time.Now().Format(time.RFC850))
	log.Printf("Parameters entered: visualization %v, randomization %v, config path was set to `%s`, output path was set to `%s`.",
		*visualizationFlag, *randomizationFlag, *configPathFlag, *outputPathFlag)

	gatewaysFile, err := os.Open(*gatewaysPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(gatewaysFile)
	gateways := df.Col("Gateway IDs").Records()
	experimentsGatewayIndex:=0

	configFile, err := os.Open(*configPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	df = dataframe.ReadCSV(configFile)

	var experimentsWaitGroup sync.WaitGroup
	for experimentIndex := 0; experimentIndex < df.Nrow(); experimentIndex++ {
		bursts, _ := df.Elem(experimentIndex, 0).Int()
		burstSize := strings.Split(df.Elem(experimentIndex, 1).String(), " ")
		payloadLengthBytes, _ := df.Elem(experimentIndex, 2).Int()
		lambdaIncrementLimit, _ := df.Elem(experimentIndex, 3).Int()
		frequencySeconds, _ := df.Elem(experimentIndex, 4).Int()
		endpointsAssigned, _ := df.Elem(experimentIndex, 5).Int()

		experimentsWaitGroup.Add(1)
		go experiment.RunExperiment(&experimentsWaitGroup, outputDirectoryPath, configuration.ExperimentConfig{
			Bursts:               bursts,
			BurstSizes:           burstSize,
			PayloadLengthBytes:   payloadLengthBytes,
			FrequencySeconds:     frequencySeconds,
			LambdaIncrementLimit: lambdaIncrementLimit,
			GatewayEndpoints:     gateways[experimentsGatewayIndex:experimentsGatewayIndex+endpointsAssigned],
			Id:                   experimentIndex,
		}, *visualizationFlag, *randomizationFlag)

		experimentsGatewayIndex+=endpointsAssigned
	}
	experimentsWaitGroup.Wait()

	log.Println("Exiting...")
}

func setupClientLogging(path string) *os.File {
	logFile, err := os.Create(filepath.Join(path, "run_logs.txt"))
	if err != nil {
		log.Fatal(err)
	}
	stdoutFileMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(stdoutFileMultiWriter)
	return logFile
}
