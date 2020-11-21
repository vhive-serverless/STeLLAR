// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
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
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/configuration"
	"lambda-benchmarking/client/experiment"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var outputPathFlag = flag.String("o", "latency-samples", "The path where latency samples should be written.")
var configPathFlag = flag.String("c", "config.csv", "Configuration file with details of experiments.")
var endpointsDirectoryPathFlag = flag.String("g", "endpoints", "Directory containing provider endpoints to be used.")
var specificExperimentFlag = flag.Int("r", -1, "Only run this particular experiment.")
var logLevelFlag = flag.String("l", "info", "Select logging level.")

func main() {
	randomSeed := time.Now().Unix()
	rand.Seed(randomSeed) // comment line for reproducible deltas
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

	setupCtrlCHandler()

	config := configuration.ReadInstructions(*endpointsDirectoryPathFlag, *configPathFlag)

	experiment.TriggerSubExperiments(config, outputDirectoryPath, *specificExperimentFlag)

	log.Info("Exiting...")
}

//setupCtrlCHandler creates a 'listener' on a new goroutine which will notify the
//program if it receives an interrupt from the OS.
func setupCtrlCHandler() {
	log.Debug("Creating Ctrl-C handler")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Ctrl+C pressed in Terminal")
		log.Info("Exiting...")
		os.Exit(0)
	}()
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
