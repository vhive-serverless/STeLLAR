package main

import (
	"flag"
	"github.com/go-gota/gota/dataframe"
	"io"
	"lambda-benchmarking/client/experiment"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//Note: those variables are pointers
var randomizationFlag = flag.Bool("randomized", true, "If true, sample deltas from a scaled and shifted standard normal distribution.")
var visualizationFlag = flag.String("visualization", "CDF", "The type of visualization to create (per-burst histogram \"histogram\" or empirical CDF \"CDF\").")
var outputPathFlag = flag.String("outputPath", "latency-samples", "The path where latency samples should be written.")
var configPathFlag = flag.String("configPath", "config.csv", "Configuration file with details of experiments.")
var gatewaysPathFlag = flag.String("gatewaysPath", "gateways.csv", "File containing ids of gateways to be used.")
var runExperimentFlag = flag.Int("runExperiment", -1, "Only run this particular experiment.")

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

	log.Printf("Started benchmarking HTTP client on %v.", time.Now().UTC().Format(time.RFC850))
	log.Printf(`Parameters entered: visualization %v, randomization %v, config path was set to %s,
		output path was set to %s, runExperiment is %d`, *visualizationFlag, *randomizationFlag, *configPathFlag,
		*outputPathFlag, *runExperimentFlag)

	gatewaysFile, err := os.Open(*gatewaysPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(gatewaysFile)
	gateways := df.Col("Gateway IDs").Records()
	experimentsGatewayIndex := 0

	configFile, err := os.Open(*configPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	df = dataframe.ReadCSV(configFile)

	var experimentsWaitGroup sync.WaitGroup
	if *runExperimentFlag != -1 {
		if *runExperimentFlag < 0 || *runExperimentFlag >= df.Nrow() {
			panic("runExperiment parameter is invalid")
		}
		experimentsWaitGroup.Add(1)
		experiment.ExtractConfigurationAndRunExperiment(df, *runExperimentFlag, &experimentsWaitGroup, outputDirectoryPath,
			gateways, experimentsGatewayIndex, *visualizationFlag, *randomizationFlag)
	} else {
		for experimentIndex := 0; experimentIndex < df.Nrow(); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			endpointsAssigned := experiment.ExtractConfigurationAndRunExperiment(df, experimentIndex, &experimentsWaitGroup,
				outputDirectoryPath, gateways, experimentsGatewayIndex, *visualizationFlag, *randomizationFlag)

			experimentsGatewayIndex += endpointsAssigned
		}
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
