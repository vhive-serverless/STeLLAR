package experiment

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"lambda-benchmarking/client/experiment/benchmarking"
	"lambda-benchmarking/client/experiment/configuration"
	"lambda-benchmarking/client/experiment/visualization"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func ExtractConfigurationAndRunExperiment(df dataframe.DataFrame, experimentIndex int, experimentsWaitGroup *sync.WaitGroup,
	outputDirectoryPath string, gateways []string, experimentsGatewayIndex int, visualization string, randomization bool) int {
	bursts, _ := df.Elem(experimentIndex, 0).Int()
	burstSize := strings.Split(df.Elem(experimentIndex, 1).String(), " ")
	payloadLengthBytes, _ := df.Elem(experimentIndex, 2).Int()
	lambdaIncrementLimit, _ := df.Elem(experimentIndex, 3).Int()
	frequencySeconds, _ := df.Elem(experimentIndex, 4).Int()
	endpointsAssigned, _ := df.Elem(experimentIndex, 5).Int()

	go runExperiment(experimentsWaitGroup, outputDirectoryPath, configuration.ExperimentConfig{
		Bursts:               bursts,
		BurstSizes:           burstSize,
		PayloadLengthBytes:   payloadLengthBytes,
		FrequencySeconds:     frequencySeconds,
		LambdaIncrementLimit: lambdaIncrementLimit,
		GatewayEndpoints:     gateways[experimentsGatewayIndex : experimentsGatewayIndex+endpointsAssigned],
		Id:                   experimentIndex,
	}, visualization, randomization)
	return endpointsAssigned
}

func runExperiment(experimentsWaitGroup *sync.WaitGroup, outputDirectoryPath string, config configuration.ExperimentConfig,
	visualizationType string, randomized bool) {
	defer experimentsWaitGroup.Done()

	experimentDirectoryPath := createExperimentDirectory(outputDirectoryPath, config.Id)
	csvFile := createExperimentLatenciesFile(experimentDirectoryPath)
	defer csvFile.Close()

	burstDeltas := benchmarking.CreateBurstDeltas(config.FrequencySeconds, config.Bursts, randomized)

	log.Printf("Starting experiment %d...", config.Id)
	safeExperimentWriter := benchmarking.InitializeExperimentWriter(csvFile)
	benchmarking.RunProfiler(config, burstDeltas, safeExperimentWriter)

	log.Printf("Experiment %d: flushing results to CSV file.", config.Id)
	safeExperimentWriter.Writer.Flush()

	if visualizationType == "" {
		log.Printf("Experiment %d: skipping visualization", config.Id)
	} else {
		log.Printf("Experiment %d: creating %ss from CSV file `%s`", config.Id, visualizationType,
			csvFile.Name())
		visualization.GenerateVisualization(
			visualizationType,
			config,
			burstDeltas,
			csvFile,
			experimentDirectoryPath,
		)
	}
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
	log.Printf("Creating directory for experiment %d at `%s`", id, experimentDirectoryPath)
	if err := os.Mkdir(experimentDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return experimentDirectoryPath
}
