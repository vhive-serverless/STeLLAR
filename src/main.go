// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"stellar/benchmarking"
	"stellar/setup"
	"stellar/setup/deployment/connection"
	"stellar/setup/deployment/connection/amazon"
	"time"
)

var awsUserArnNumber = flag.String("a", "356764711652", "This is used in AWS benchmarking for client authentication.")
var outputPathFlag = flag.String("o", "latency-samples", "The directory path where latency samples should be written.")
var configPathFlag = flag.String("c", "experiments/tests/aws/data-transfer.json", "Configuration file with experiment details.")
var endpointsDirectoryPathFlag = flag.String("g", "endpoints", "Directory containing provider endpoints to be used.")
var specificExperimentFlag = flag.Int("r", -1, "Only run this particular experiment.")
var logLevelFlag = flag.String("l", "info", "Select logging level.")
var writeToDatabaseFlag = flag.Bool("db", false, "This bool flag specifies whether statistics should be written to the database")
var serverlessDeployment = flag.Bool("s", true, "Use serverless.com framework for deployment. ")

func main() {
	startTime := time.Now()
	randomSeed := startTime.Unix()
	rand.Seed(randomSeed) // comment line for reproducible inter-arrival times
	flag.Parse()

	outputDirectoryPath := filepath.Join(*outputPathFlag, time.Now().Format(time.RFC850))
	log.Infof("Creating directory for this run at `%s`", outputDirectoryPath)
	if err := os.MkdirAll(outputDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFile := setupLogging(outputDirectoryPath)
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
		log.Infof("number of routes %d, numebr of endpoints %d", len(config.SubExperiments[0].Routes), len(config.SubExperiments[0].Endpoints))
		benchmarking.TriggerSubExperiments(config, outputDirectoryPath, *specificExperimentFlag)

		log.Info("Starting functions removal from cloud.")
		setup.RemoveService(config.Provider, serverlessDirPath, len(config.SubExperiments))
	} else {
		setup.ProvisionFunctions(config)
		benchmarking.TriggerSubExperiments(config, outputDirectoryPath, *specificExperimentFlag)
	}

	log.Infof("Done in %v, exiting...", time.Since(startTime))
}

func setupLogging(path string) *os.File {
	loggingPath := filepath.Join(path, "run_logs.txt")
	log.Debugf("Creating log file for this run at `%s`", loggingPath)
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
	}

	stdoutFileMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(stdoutFileMultiWriter)

	return logFile
}
