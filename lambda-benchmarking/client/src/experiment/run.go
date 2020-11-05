package experiment

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/benchmarking"
	"lambda-benchmarking/client/experiment/configuration"
	"lambda-benchmarking/client/experiment/visualization"
	"os"
	"path/filepath"
	"sync"
)

func TriggerExperiment(experimentsWaitGroup *sync.WaitGroup, experiment configuration.SubExperiment, outputDirectoryPath string) {
	defer experimentsWaitGroup.Done()
	experimentDirectoryPath, latenciesFile := createExperimentOutput(outputDirectoryPath, experiment)
	runExperiment(latenciesFile, experimentDirectoryPath, experiment)
}

func createExperimentOutput(path string, experiment configuration.SubExperiment) (string, *os.File) {
	directoryPath := filepath.Join(path, fmt.Sprintf("%d_%s", experiment.Id, experiment.Title))
	log.Infof("SubExperiment %d: Creating directory at `%s`", experiment.Id, directoryPath)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	latenciesPath := filepath.Join(directoryPath, "latencies.csv")
	log.Infof("SubExperiment %d: Creating latencies file at `%s`", experiment.Id, latenciesPath)
	csvFile, err := os.Create(latenciesPath)
	if err != nil {
		log.Fatal(err)
	}

	return directoryPath, csvFile
}

func runExperiment(latenciesFile *os.File, experimentDirectoryPath string, experiment configuration.SubExperiment) {
	log.Infof("SubExperiment %d: Starting...", experiment.Id)
	burstDeltas := generateIAT(experiment)
	benchmarking.RunProfiler(experiment, burstDeltas, benchmarking.InitializeExperimentWriter(latenciesFile))
	visualization.GenerateVisualization(experiment, burstDeltas, latenciesFile, experimentDirectoryPath)

	if err := latenciesFile.Close(); err != nil {
		log.Fatal(err)
	}
	log.Infof("SubExperiment %d successfully finished.", experiment.Id)
}
