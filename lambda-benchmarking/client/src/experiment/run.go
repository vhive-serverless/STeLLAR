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

func TriggerExperiment(experimentsWaitGroup *sync.WaitGroup, experiment configuration.Experiment, outputDirectoryPath string, visualization string) {
	defer experimentsWaitGroup.Done()
	experimentDirectoryPath, latenciesFile := createExperimentOutput(outputDirectoryPath, experiment.Id)
	runExperiment(latenciesFile, experimentDirectoryPath, experiment, visualization)
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

func runExperiment(latenciesFile *os.File, experimentDirectoryPath string, experiment configuration.Experiment, visualizationType string) {
	log.Infof("Experiment %d: Starting...", experiment.Id)
	burstDeltas := generateIAT(experiment)
	benchmarking.RunProfiler(experiment, burstDeltas, benchmarking.InitializeExperimentWriter(latenciesFile))
	visualization.GenerateVisualization(visualizationType, experiment, burstDeltas, latenciesFile, experimentDirectoryPath)

	if err := latenciesFile.Close(); err != nil {
		log.Fatal(err)
	}
	log.Infof("Experiment %d successfully finished.", experiment.Id)
}
