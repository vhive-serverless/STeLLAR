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

//TriggerExperiment will run the sub-experiment specified by the passed configuration object. It creates
//its own directory, as well as separate visualizations and latency files.
func TriggerExperiment(experimentsWaitGroup *sync.WaitGroup, experiment configuration.SubExperiment, outputDirectoryPath string) {
	defer experimentsWaitGroup.Done()
	experimentDirectoryPath, latenciesFile := createExperimentOutput(outputDirectoryPath, experiment)
	runExperiment(latenciesFile, experimentDirectoryPath, experiment)
}

func createExperimentOutput(path string, experiment configuration.SubExperiment) (string, *os.File) {
	directoryPath := filepath.Join(path, fmt.Sprintf("%d_%s", experiment.ID, experiment.Title))
	log.Infof("SubExperiment %d: Creating directory at `%s`", experiment.ID, directoryPath)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	latenciesPath := filepath.Join(directoryPath, "latencies.csv")
	log.Infof("SubExperiment %d: Creating latencies file at `%s`", experiment.ID, latenciesPath)
	csvFile, err := os.Create(latenciesPath)
	if err != nil {
		log.Fatal(err)
	}

	return directoryPath, csvFile
}

func runExperiment(latenciesFile *os.File, experimentDirectoryPath string, experiment configuration.SubExperiment) {
	log.Infof("SubExperiment %d: Starting...", experiment.ID)
	burstDeltas := generateIAT(experiment)
	benchmarking.RunProfiler(experiment, burstDeltas, benchmarking.NewExperimentWriter(latenciesFile))
	visualization.GenerateVisualization(experiment, burstDeltas, latenciesFile, experimentDirectoryPath)

	if err := latenciesFile.Close(); err != nil {
		log.Fatal(err)
	}
	log.Infof("SubExperiment %d successfully finished.", experiment.ID)
}
