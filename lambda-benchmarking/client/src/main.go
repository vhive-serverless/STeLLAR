package main

import (
	"flag"
	"github.com/go-gota/gota/dataframe"
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/experiment"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var visualizationFlag = flag.String("visualization", "CDF", "The type of visualization to create (per-burst histogram \"histogram\" or empirical CDF \"CDF\").")
var outputPathFlag = flag.String("outputPath", "latency-samples", "The path where latency samples should be written.")
var configPathFlag = flag.String("configPath", "config.csv", "Configuration file with details of experiments.")
var gatewaysPathFlag = flag.String("gatewaysPath", "gateways.csv", "File containing ids of gateways to be used.")
var runExperimentFlag = flag.Int("runExperiment", -1, "Only run this particular experiment.")
var logLevelFlag = flag.String("logLevel", "info", "Select logging level.")

func main() {
	rand.Seed(time.Now().Unix()) // comment line for reproducible deltas
	flag.Parse()

	outputDirectoryPath := filepath.Join(*outputPathFlag, time.Now().Format(time.RFC850))
	log.Infof("Creating directory for this run at `%s`", outputDirectoryPath)
	if err := os.MkdirAll(outputDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFile := setupClientLogging(outputDirectoryPath)
	defer logFile.Close()

	log.Infof("Started benchmarking HTTP client on %v.", time.Now().UTC().Format(time.RFC850))
	log.Infof(`Parameters entered: visualization %v, config path was set to %s,
		output path was set to %s, runExperiment is %d`, *visualizationFlag, *configPathFlag,
		*outputPathFlag, *runExperimentFlag)

	gatewaysFile, err := os.Open(*gatewaysPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(gatewaysFile)
	gateways := df.Col("Gateway ID").Records()
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
			gateways, experimentsGatewayIndex, *visualizationFlag)
	} else {
		for experimentIndex := 0; experimentIndex < df.Nrow(); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			endpointsAssigned := experiment.ExtractConfigurationAndRunExperiment(df, experimentIndex, &experimentsWaitGroup,
				outputDirectoryPath, gateways, experimentsGatewayIndex, *visualizationFlag)

			experimentsGatewayIndex += endpointsAssigned
		}
	}
	experimentsWaitGroup.Wait()

	log.Infof("Exiting...")
}

func setupClientLogging(path string) *os.File {
	logFile, err := os.Create(filepath.Join(path, "run_logs.txt"))
	if err != nil {
		log.Fatal(err)
	}

	switch *logLevelFlag {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}

	stdoutFileMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(stdoutFileMultiWriter)

	return logFile
}
