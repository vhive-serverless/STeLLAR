package main

import (
	"flag"
	"github.com/go-gota/gota/dataframe"
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/experiment"
	"lambda-benchmarking/client/experiment/configuration"
	"lambda-benchmarking/client/prompts"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var outputPathFlag = flag.String("o", "latency-samples", "The path where latency samples should be written.")
var configPathFlag = flag.String("c", "config.csv", "Configuration file with details of experiments.")
var gatewaysPathFlag = flag.String("g", "gateways.csv", "File containing ids of gateways to be used.")
var runExperimentFlag = flag.Int("r", -1, "Only run this particular experiment.")
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

	config := readInstructions()

	triggerExperiments(config, outputDirectoryPath)

	log.Info("Exiting...")
}

func triggerExperiments(config configuration.Configuration, outputDirectoryPath string) {
	var experimentsWaitGroup sync.WaitGroup

	switch *runExperimentFlag {
	case -1: // run all experiments
		for experimentIndex := 0; experimentIndex < len(config.SubExperiments); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			go experiment.TriggerExperiment(&experimentsWaitGroup, config.SubExperiments[experimentIndex], outputDirectoryPath)

			if config.Sequential {
				experimentsWaitGroup.Wait()
			}
		}
	default:
		if *runExperimentFlag < 0 || *runExperimentFlag >= len(config.SubExperiments) {
			log.Fatalf("Parameter `runExperiment` is invalid: %d", *runExperimentFlag)
		}

		experimentsWaitGroup.Add(1)
		go experiment.TriggerExperiment(&experimentsWaitGroup, config.SubExperiments[*runExperimentFlag], outputDirectoryPath)
	}

	experimentsWaitGroup.Wait()
}

func readInstructions() configuration.Configuration {
	log.Debugf("Reading gateways file for this run from `%s`", *gatewaysPathFlag)
	gatewaysFile, err := os.Open(*gatewaysPathFlag)
	if err != nil {
		log.Fatalf("Could not read gateways file: %s", err.Error())
	}
	gatewaysDF := dataframe.ReadCSV(gatewaysFile)
	memoryToGatewayIDs, memoryToLastAssignedIndex := mapMemoryToGateways(gatewaysDF)

	log.Debugf("Reading config file for this run from `%s`", *configPathFlag)
	configFile, err := os.Open(*configPathFlag)
	if err != nil {
		log.Fatalf("Could not read config file: %s", err.Error())
	}

	config := configuration.Extract(configFile)
	for index := range config.SubExperiments {
		config.SubExperiments[index].ID = index
		assignGatewaysToExperiment(memoryToGatewayIDs, memoryToLastAssignedIndex, &config.SubExperiments[index])

		// Issue warning if sending too many requests in a single burst
		const manyRequestsInBurstThreshold = 2000
		for _, burstSize := range config.SubExperiments[index].BurstSizes {
			if burstSize > manyRequestsInBurstThreshold {
				log.Warnf("Experiment %d has a burst of size %d, NIC (Network Interface Controller) contention may occur.",
					index, burstSize)
			}
		}

		// Issue warning if generating too many files
		const manyFilesWarningThreshold = 500
		chosenVisualization := config.SubExperiments[index].Visualization
		burstsNumber := config.SubExperiments[index].Bursts
		if burstsNumber >= manyFilesWarningThreshold && (chosenVisualization == "all" || chosenVisualization == "histogram") {
			log.Warnf("Generating histograms for each burst, this will create a large number (%d) of new files.",
				burstsNumber)
		}
	}

	log.Debugf("Extracted %d sub-experiments from given configuration file.", len(config.SubExperiments))
	return config
}

func mapMemoryToGateways(gatewaysDF dataframe.DataFrame) (map[int64][]string, map[int64]int) {
	memoryToListOfGatewayIDs := make(map[int64][]string)
	memoryToLastAssignedIndex := make(map[int64]int)
	for idx, record := range gatewaysDF.Records() {
		if idx == 0 {
			continue
		}

		desiredMemory, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Fatalf("Could not parse function memory %s in configuration file.", record[1])
		}
		memoryToLastAssignedIndex[desiredMemory] = 0
		memoryToListOfGatewayIDs[desiredMemory] = append(memoryToListOfGatewayIDs[desiredMemory], record[0])
	}
	return memoryToListOfGatewayIDs, memoryToLastAssignedIndex
}

func assignGatewaysToExperiment(gateways map[int64][]string, memoryToLastAssignedIndex map[int64]int, experiment *configuration.SubExperiment) {
	lastAssignedIndexExcl := memoryToLastAssignedIndex[experiment.FunctionMemoryMB]
	newLastAssignedIndexExcl := lastAssignedIndexExcl + experiment.GatewaysNumber
	nrGatewaysWithDesiredMemory := len(gateways[experiment.FunctionMemoryMB])

	if newLastAssignedIndexExcl > nrGatewaysWithDesiredMemory {
		remainingGatewaysToAssign := nrGatewaysWithDesiredMemory - lastAssignedIndexExcl
		log.Errorf("Not enough remaining gateways were found in the given gateways file with requested memory %dMB, found %d but trying to assign from %d to %d. Experiment `%s` will be assigned %d gateways.",
			experiment.FunctionMemoryMB, remainingGatewaysToAssign,
			lastAssignedIndexExcl, newLastAssignedIndexExcl,
			experiment.Title, remainingGatewaysToAssign)

		if remainingGatewaysToAssign <= 0 {
			log.Fatalf("Cannot assign %d gateways to an experiment.", remainingGatewaysToAssign)
		}

		if !prompts.PromptForConfirmation("Would you like to continue with this setting?") {
			os.Exit(0)
		}

		newLastAssignedIndexExcl = nrGatewaysWithDesiredMemory
	}
	memoryToLastAssignedIndex[experiment.FunctionMemoryMB] = newLastAssignedIndexExcl
	experiment.GatewayEndpoints = gateways[experiment.FunctionMemoryMB][lastAssignedIndexExcl:newLastAssignedIndexExcl]
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
