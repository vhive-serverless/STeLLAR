package experiment

import (
	"fmt"
	"lambda-benchmarking/client/experiment/benchmarking"
	"lambda-benchmarking/client/experiment/configuration"
	"lambda-benchmarking/client/experiment/visualization"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func RunExperiment(experimentsWaitGroup *sync.WaitGroup, outputDirectoryPath string, config configuration.ExperimentConfig,
	visualizationType string, randomized bool) {
	defer experimentsWaitGroup.Done()

	experimentDirectoryPath := createExperimentDirectory(outputDirectoryPath, config.Id)
	csvFile := createExperimentLatenciesFile(experimentDirectoryPath)
	defer csvFile.Close()

	burstDeltas := benchmarking.CreateBurstDeltas(config.FrequencySeconds, config.Bursts, randomized)
	burstRelativeDeltas := benchmarking.MakeBurstDeltasRelative(burstDeltas)

	log.Printf("Starting experiment %d...", config.Id)
	safeExperimentWriter := benchmarking.InitializeExperimentWriter(csvFile)
	benchmarking.RunProfiler(config, burstRelativeDeltas, safeExperimentWriter)

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
			burstRelativeDeltas,
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
