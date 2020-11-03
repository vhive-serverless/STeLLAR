package main

import (
	"flag"
	"github.com/go-gota/gota/dataframe"
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/experiment"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var visualizationFlag = flag.String("visualization", "all", "The type of visualization to create (`histogram`, `CDF`, `all`).")
var outputPathFlag = flag.String("outputPath", "latency-samples", "The path where latency samples should be written.")
var configPathFlag = flag.String("configPath", "config.csv", "Configuration file with details of experiments.")
var gatewaysPathFlag = flag.String("gatewaysPath", "gateways.csv", "File containing ids of gateways to be used.")
var runExperimentFlag = flag.Int("runExperiment", -1, "Only run this particular experiment.")
var logLevelFlag = flag.String("logLevel", "info", "Select logging level.")

func main() {
	randomSeed := time.Now().Unix()
	rand.Seed(randomSeed) // comment line for reproducible deltas
	flag.Parse()

	outputDirectoryPath := filepath.Join(*outputPathFlag, time.Now().Format(time.RFC850))
	log.Infof("Creating directory for this run at `%s`", outputDirectoryPath)
	if err := os.MkdirAll(outputDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFile := setupClientLogging(outputDirectoryPath)
	defer logFile.Close()

	log.Infof("Started benchmarking HTTP client on %v (random seed %d).", time.Now().UTC().Format(time.RFC850), randomSeed)
	log.Infof(`Parameters entered: visualization %v, gateways path was set to %s, config path was set to %s, output path was set to %s, runExperiment is %d`, *visualizationFlag, *gatewaysPathFlag, *configPathFlag,
		*outputPathFlag, *runExperimentFlag)

	log.Debug("Creating Ctrl-C handler")
	SetupCtrlCHandler()

	gateways, configDF := readInstructions()
	experimentsGatewayIndex := 0

	var experimentsWaitGroup sync.WaitGroup
	if *runExperimentFlag != -1 {
		if *runExperimentFlag < 0 || *runExperimentFlag >= configDF.Nrow() {
			panic("runExperiment parameter is invalid")
		}
		experimentsWaitGroup.Add(1)
		experiment.TriggerExperiment(configDF, *runExperimentFlag, &experimentsWaitGroup, outputDirectoryPath, gateways, experimentsGatewayIndex, *visualizationFlag)
	} else {
		for experimentIndex := 0; experimentIndex < configDF.Nrow(); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			endpointsAssigned := experiment.TriggerExperiment(configDF, experimentIndex, &experimentsWaitGroup,
				outputDirectoryPath, gateways, experimentsGatewayIndex, *visualizationFlag)

			experimentsGatewayIndex += endpointsAssigned
		}
	}
	experimentsWaitGroup.Wait()

	log.Info("Exiting...")
}

func readInstructions() ([]string, dataframe.DataFrame) {
	log.Debugf("Reading config file for this run from `%s`", *configPathFlag)
	configFile, err := os.Open(*configPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	configDF := dataframe.ReadCSV(configFile)

	log.Debugf("Reading gateways file for this run from `%s`", *gatewaysPathFlag)
	gatewaysFile, err := os.Open(*gatewaysPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	gatewaysDF := dataframe.ReadCSV(gatewaysFile)
	gateways := gatewaysDF.Col("Gateway ID").Records()

	return gateways, configDF
}

// SetupCtrlCHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS.
func SetupCtrlCHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Ctrl+C pressed in Terminal")
		log.Info("Exiting...")
		os.Exit(0)
	}()
}

func setupClientLogging(path string) *os.File {
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
