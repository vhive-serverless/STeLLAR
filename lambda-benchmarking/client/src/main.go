package main

import (
	"flag"
	"github.com/go-gota/gota/dataframe"
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/experiment"
	"lambda-benchmarking/client/experiment/configuration"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var outputPathFlag = flag.String("o", "latency-samples", "The path where latency samples should be written.")
var configPathFlag = flag.String("c", "config.csv", "Configuration file with details of experiments.")
var gatewaysPathFlag = flag.String("g", "gateways.csv", "File containing ids of gateways to be used.")
var runExperimentFlag = flag.Int("r", -1, "Only run this particular experiment.")
var sequentialFlag = flag.Bool("s", true, "Run the experiments sequentially or simultaneously.")
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

	logFile := setupClientLogging(outputDirectoryPath)
	defer logFile.Close()

	log.Infof("Started benchmarking HTTP client on %v (random seed %d).",
		time.Now().UTC().Format(time.RFC850), randomSeed)
	log.Infof("Selected gateways path: %s", *gatewaysPathFlag)
	log.Infof("Selected config path: %s", *configPathFlag)
	log.Infof("Selected output path: %s", *outputPathFlag)
	log.Infof("Selected experiment (-1 for all): %d", *runExperimentFlag)

	log.Debug("Creating Ctrl-C handler")
	setupCtrlCHandler()

	experiments := readInstructions()

	triggerExperiments(experiments, outputDirectoryPath)

	log.Info("Exiting...")
}

func triggerExperiments(experiments []configuration.SubExperiment, outputDirectoryPath string) {
	var experimentsWaitGroup sync.WaitGroup

	switch *runExperimentFlag {
	case -1: // run all experiments
		for experimentIndex := 0; experimentIndex < len(experiments); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			go experiment.TriggerExperiment(&experimentsWaitGroup, experiments[experimentIndex], outputDirectoryPath)
		}
	default:
		if *runExperimentFlag < 0 || *runExperimentFlag >= len(experiments) {
			log.Fatalf("Parameter `runExperiment` is invalid: %d", *runExperimentFlag)
		}

		experimentsWaitGroup.Add(1)
		go experiment.TriggerExperiment(&experimentsWaitGroup, experiments[*runExperimentFlag], outputDirectoryPath)
	}

	experimentsWaitGroup.Wait()
}

func readInstructions() []configuration.SubExperiment {
	log.Debugf("Reading gateways file for this run from `%s`", *gatewaysPathFlag)
	gatewaysFile, err := os.Open(*gatewaysPathFlag)
	if err != nil {
		log.Fatalf("Could not read gateways file: %s", err.Error())
	}
	gatewaysDF := dataframe.ReadCSV(gatewaysFile)
	gateways := gatewaysDF.Col("Gateway ID").Records()

	log.Debugf("Reading config file for this run from `%s`", *configPathFlag)
	configFile, err := os.Open(*configPathFlag)
	if err != nil {
		log.Fatalf("Could not read config file: %s", err.Error())
	}

	experimentsGatewayIndex := 0
	experiments := configuration.Extract(configFile)
	for index, exp := range experiments {
		experiments[index].ID = index
		assignGatewaysToExperiment(gateways, experimentsGatewayIndex, exp.GatewaysNumber, gatewaysFile.Name(), &experiments[index])
		experimentsGatewayIndex += exp.GatewaysNumber
	}

	log.Debugf("Extracted %d experiments from given configuration file.", len(experiments))
	return experiments
}

func assignGatewaysToExperiment(gateways []string, currExpGatewayIndex int, experimentGatewaysReq int,
	gatewaysFileName string, experiment *configuration.SubExperiment) {
	if currExpGatewayIndex+experimentGatewaysReq > len(gateways) {
		log.Warnf("Not enough gateways were found in the given gateways file `%s`, found %d but trying to assign from %d to %d. Experiment `%s` will be assigned 0 gateways.",
			gatewaysFileName, len(gateways), currExpGatewayIndex, currExpGatewayIndex+experimentGatewaysReq, experiment.Title)
	}
	experiment.GatewayEndpoints = gateways[currExpGatewayIndex : currExpGatewayIndex+experimentGatewaysReq]
}

// setupCtrlCHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS.
func setupCtrlCHandler() {
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
