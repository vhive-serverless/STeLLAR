package experiment

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/benchmarking"
	"lambda-benchmarking/client/experiment/configuration"
	"lambda-benchmarking/client/experiment/visualization"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func TriggerExperiment(df dataframe.DataFrame, experimentIndex int, experimentsWaitGroup *sync.WaitGroup,
	outputDirectoryPath string, gateways []string, experimentsGatewayIndex int, visualization string) int {
	experimentDirectoryPath, latenciesFile := createExperimentOutput(outputDirectoryPath, experimentIndex)
	config, endpointsAssigned := extractConfiguration(df, experimentIndex, gateways, experimentsGatewayIndex)
	go runExperiment(experimentsWaitGroup, latenciesFile, experimentDirectoryPath, *config, visualization)

	return endpointsAssigned
}

func getIncrementLimits(incrementLimitStrings []string) []int64 {
	var functionIncrementLimits []int64
	for _, incrementLimitString := range incrementLimitStrings {
		parsedIncrementLimit, err := strconv.ParseInt(incrementLimitString, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		functionIncrementLimits = append(functionIncrementLimits, parsedIncrementLimit)
	}
	return functionIncrementLimits
}

func extractConfiguration(df dataframe.DataFrame, experimentIndex int, gateways []string, experimentsGatewayIndex int) (*configuration.ExperimentConfig, int) {
	bursts, _ := df.Elem(experimentIndex, 0).Int()
	burstSize := strings.Split(df.Elem(experimentIndex, 1).String(), " ")
	iatType := df.Elem(experimentIndex, 2).String()
	payloadLengthBytes, _ := df.Elem(experimentIndex, 3).Int()
	incrementLimitStrings := strings.Split(df.Elem(experimentIndex, 4).String(), " ")
	frequencySeconds := df.Elem(experimentIndex, 5).Float()
	endpointsAssigned, _ := df.Elem(experimentIndex, 6).Int()
	providerBenchmarked := df.Elem(experimentIndex, 7).String()

	return &configuration.ExperimentConfig{
		Bursts:                  bursts,
		BurstSizes:              burstSize,
		PayloadLengthBytes:      payloadLengthBytes,
		FrequencySeconds:        frequencySeconds,
		FunctionIncrementLimits: getIncrementLimits(incrementLimitStrings),
		GatewayEndpoints:        gateways[experimentsGatewayIndex : experimentsGatewayIndex+endpointsAssigned],
		Id:                      experimentIndex,
		IatType:                 iatType,
		Provider:                providerBenchmarked,
	}, endpointsAssigned
}

func createExperimentOutput(path string, id int) (string, *os.File) {
	directoryPath := filepath.Join(path, fmt.Sprintf("experiment_%d", id))
	log.Infof("Experiment %d: Creating directory at `%s`", id, directoryPath)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	latenciesPath := filepath.Join(directoryPath, "latencies.csv")
	log.Infof("Experiment %d: Creating latencies file at `%s`", id, latenciesPath)
	csvFile, err := os.Create(latenciesPath)
	if err != nil {
		log.Fatal(err)
	}

	return directoryPath, csvFile
}

func runExperiment(experimentsWaitGroup *sync.WaitGroup, latenciesFile *os.File, experimentDirectoryPath string, config configuration.ExperimentConfig, visualizationType string) {
	defer experimentsWaitGroup.Done()
	log.Infof("Experiment %d: Starting...", config.Id)
	burstDeltas := generateIAT(config.FrequencySeconds, config.Bursts, config.IatType, config.Id)
	benchmarking.RunProfiler(config, burstDeltas, benchmarking.InitializeExperimentWriter(latenciesFile))
	visualization.GenerateVisualization(visualizationType, config, burstDeltas, latenciesFile, experimentDirectoryPath)

	if err := latenciesFile.Close(); err != nil {
		log.Fatal(err)
	}
	log.Infof("Experiment %d: done", config.Id)
}
