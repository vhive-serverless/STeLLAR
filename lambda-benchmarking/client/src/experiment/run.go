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
	"strings"
	"sync"
)

func ExtractConfigurationAndRunExperiment(df dataframe.DataFrame, experimentIndex int, experimentsWaitGroup *sync.WaitGroup, outputDirectoryPath string, gateways []string, experimentsGatewayIndex int, visualization string) int {
	experimentDirectoryPath := createExperimentDirectory(outputDirectoryPath, experimentIndex)
	latenciesFile := createExperimentLatenciesFile(experimentDirectoryPath)
	safeExperimentWriter := benchmarking.InitializeExperimentWriter(latenciesFile)

	bursts, _ := df.Elem(experimentIndex, 0).Int()
	burstSize := strings.Split(df.Elem(experimentIndex, 1).String(), " ")
	iatType := df.Elem(experimentIndex, 2).String()
	payloadLengthBytes, _ := df.Elem(experimentIndex, 3).Int()
	functionIncrementLimits := strings.Split(df.Elem(experimentIndex, 4).String(), " ")
	frequencySeconds := df.Elem(experimentIndex, 5).Float()
	endpointsAssigned, _ := df.Elem(experimentIndex, 6).Int()

	go runExperiment(experimentsWaitGroup, latenciesFile, experimentDirectoryPath, configuration.ExperimentConfig{
		Bursts:                  bursts,
		BurstSizes:              burstSize,
		PayloadLengthBytes:      payloadLengthBytes,
		FrequencySeconds:        frequencySeconds,
		FunctionIncrementLimits: getIncrementLimits(functionIncrementLimits),
		GatewayEndpoints:        gateways[experimentsGatewayIndex : experimentsGatewayIndex+endpointsAssigned],
		Id:                      experimentIndex,
		IatType:                 iatType,
	}, visualization, safeExperimentWriter)

	if err := latenciesFile.Close(); err != nil {
		log.Fatal(err)
	}
	return endpointsAssigned
}

func createExperimentLatenciesFile(experimentDirectoryPath string) *os.File {
	csvFile, err := os.Create(filepath.Join(experimentDirectoryPath, "latencies.csv"))
	if err != nil {
		log.Fatal(err)
	}
	return csvFile
}

func createExperimentDirectory(path string, id int) string {
	experimentDirectoryPath := filepath.Join(path, fmt.Sprintf("experiment_%d", id))
	log.Infof("Creating directory for experiment %d at `%s`", id, experimentDirectoryPath)
	if err := os.MkdirAll(experimentDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return experimentDirectoryPath
}

func runExperiment(experimentsWaitGroup *sync.WaitGroup, latenciesFile *os.File, experimentDirectoryPath string,
	config configuration.ExperimentConfig, visualizationType string, safeExperimentWriter *benchmarking.SafeWriter) {
	defer experimentsWaitGroup.Done()

	burstDeltas := benchmarking.CreateIAT(config.FrequencySeconds, config.Bursts, config.IatType)

	log.Infof("Starting experiment %d...", config.Id)
	benchmarking.RunProfiler(config, burstDeltas, safeExperimentWriter)

	switch visualizationType {
	case "":
		log.Infof("Experiment %d: skipping visualization", config.Id)
	case "CDF":
		log.Infof("Experiment %d: creating %ss from CSV file `%s`", config.Id, visualizationType, latenciesFile.Name())
		visualization.GenerateVisualization(
			visualizationType,
			config,
			burstDeltas,
			latenciesFile,
			experimentDirectoryPath,
		)
	default:
		log.Errorf("Experiment %d: unrecognized visualization %s, skipping", config.Id, visualizationType)
	}
	log.Infof("Experiment %d: done", config.Id)
}
