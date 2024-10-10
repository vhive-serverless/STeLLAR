// main.go

package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"stellar/benchmarking"
	"stellar/setup"
	"stellar/setup/deployment/connection"
	"stellar/setup/deployment/connection/amazon"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var awsUserArnNumber = flag.String("a", "356764711652", "This is used in AWS benchmarking for client authentication.")
var outputPathFlag = flag.String("o", "latency-samples", "The directory path where latency samples should be written.")
var configPathFlag = flag.String("c", "experiments/tests/aws/data-transfer.json", "Configuration file with experiment details.")
var endpointsDirectoryPathFlag = flag.String("g", "endpoints", "Directory containing provider endpoints to be used.")
var specificExperimentFlag = flag.Int("r", -1, "Only run this particular experiment.")
var logLevelFlag = flag.String("l", "info", "Select logging level.")
var writeToDatabaseFlag = flag.Bool("db", false, "This bool flag specifies whether statistics should be written to the database")
var serverlessDeployment = flag.Bool("s", true, "Use serverless.com framework for deployment.")
var warmFlag = flag.Bool("w", false, "Warm up the serverless function with 1 invocation before recording statistics. For continuous benchmarking.")

// Create a new variable for the first if statement
var azureDeployment = flag.Bool("azure", true, "Run Azure deployment.")

func main() {
	startTime := time.Now()
	randomSeed := startTime.Unix()
	rand.Seed(randomSeed) // comment line for reproducible inter-arrival times

	flag.Parse() // Move flag.Parse() before setupLogging()

	// Initialize logging
	logFile := setupLogging("logs")
	defer logFile.Close()

	if *azureDeployment {
		// Call the RunAzureDeployment function
		setup.RunAzureDeployment()
	} else {
		// Existing code that was causing errors, moved into the else block
		outputDirectoryPath := filepath.Join(*outputPathFlag, strconv.FormatInt(time.Now().Unix(), 10))
		log.Infof("Creating directory for this run at `%s`", outputDirectoryPath)
		if err := os.MkdirAll(outputDirectoryPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		// Re-initialize logging with the new output directory
		logFile = setupLogging(outputDirectoryPath)
		defer logFile.Close()

		log.Infof("Started benchmarking HTTP client on %v with random seed %d.",
			time.Now().UTC().Format(time.RFC850), randomSeed)
		log.Infof("Selected endpoints directory path: %s", *endpointsDirectoryPathFlag)
		log.Infof("Selected config path: %s", *configPathFlag)
		log.Infof("Selected output path: %s", *outputPathFlag)
		log.Infof("Selected experiment (-1 for all): %d", *specificExperimentFlag)

		config := setup.ExtractConfiguration(*configPathFlag)

		amazon.UserARNNumber = *awsUserArnNumber

		// We find the busy-spinning time based on the host where the tool is run, i.e., not AWS or other providers
		setup.FindBusySpinIncrements(&config)

		// Pick between deployment methods
		connection.Initialize(config.Provider, *endpointsDirectoryPathFlag, "./setup/deployment/raw-code/functions/producer-consumer/api-template.json")
		if *serverlessDeployment {
			serverlessDirPath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/", config.Provider)
			setup.ProvisionFunctionsServerless(&config, serverlessDirPath)
			log.Infof("number of routes %d, number of endpoints %d", len(config.SubExperiments[0].Routes), len(config.SubExperiments[0].Endpoints))
			benchmarking.TriggerSubExperiments(config, outputDirectoryPath, *specificExperimentFlag, *writeToDatabaseFlag, *warmFlag)

			log.Info("Starting functions removal from cloud.")
			setup.RemoveService(&config, serverlessDirPath)
		} else {
			setup.ProvisionFunctions(config)
			benchmarking.TriggerSubExperiments(config, outputDirectoryPath, *specificExperimentFlag, *writeToDatabaseFlag, *warmFlag)
		}
	}

	log.Infof("Done in %v, exiting...", time.Since(startTime))
}

func setupLogging(path string) *os.File {
	loggingPath := filepath.Join(path, "run_logs.txt")
	log.Debugf("Creating log file for this run at `%s`", loggingPath)

	// Ensure the directory exists
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(loggingPath)
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
	default:
		// Set default log level if flag is not set
		log.SetLevel(log.InfoLevel)
	}

	stdoutFileMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(stdoutFileMultiWriter)

	return logFile
}
