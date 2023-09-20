// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package benchmarking

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"stellar/benchmarking/writers"
	"stellar/setup"
	"sync"
	"time"
)

// TriggerSubExperiments will run the sub-experiments specified by the passed configuration object. It creates
// a directory for each sub-experiment, as well as separate visualizations and latency files.
func TriggerSubExperiments(config setup.Configuration, outputDirectoryPath string, specificExperiment int, writeToDatabase bool, warmExperiment bool) {
	var experimentsWaitGroup sync.WaitGroup

	switch specificExperiment {
	case -1: // run all experiments
		for experimentIndex := 0; experimentIndex < len(config.SubExperiments); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			go triggerSubExperiment(&experimentsWaitGroup, config.Provider, config.SubExperiments[experimentIndex], outputDirectoryPath, writeToDatabase, warmExperiment)

			if config.Sequential {
				experimentsWaitGroup.Wait()
			}
		}
	default:
		if specificExperiment < 0 || specificExperiment >= len(config.SubExperiments) {
			log.Fatalf("Parameter `runSubExperiment` is invalid: %d", specificExperiment)
		}

		experimentsWaitGroup.Add(1)
		go triggerSubExperiment(&experimentsWaitGroup, config.Provider, config.SubExperiments[specificExperiment], outputDirectoryPath, writeToDatabase, warmExperiment)
	}

	experimentsWaitGroup.Wait()
}

func triggerSubExperiment(experimentsWaitGroup *sync.WaitGroup, provider string, experiment setup.SubExperiment, outputDirectoryPath string, writeToDatabase bool, warmExperiment bool) {
	log.Infof("[sub-experiment %d] Starting...", experiment.ID)
	defer experimentsWaitGroup.Done()

	experimentDirectoryPath, latenciesFile, statisticsFile, dataTransfersFile := createSubExperimentOutput(outputDirectoryPath, experiment)
	defer latenciesFile.Close()
	defer statisticsFile.Close()
	if dataTransfersFile != nil {
		defer dataTransfersFile.Close()
	}

	burstDeltas := generateIAT(experiment)

	log.Infof("[sub-experiment %d] Started benchmarking, scheduling %d bursts with IAT ~%vs and %d gateways (bursts/gateways*freq=%v)",
		experiment.ID, experiment.Bursts, experiment.IATSeconds, len(experiment.Endpoints),
		float64(experiment.Bursts)/float64(len(experiment.Endpoints))*experiment.IATSeconds)

	latenciesWriter := writers.NewRTTLatencyWriter(latenciesFile)
	dataTransferWriter := writers.NewDataTransferWriter(dataTransfersFile, experiment.DataTransferChainLength)

	runSubExperiment(experiment, burstDeltas, provider, latenciesWriter, dataTransferWriter)

	postProcessing(experiment, latenciesFile, burstDeltas, experimentDirectoryPath, statisticsFile, writeToDatabase, warmExperiment)

	log.Infof("[sub-experiment %d] Successfully finished.", experiment.ID)
}

func createSubExperimentOutput(path string, experiment setup.SubExperiment) (string, *os.File, *os.File, *os.File) {
	detailedTitle := fmt.Sprintf("%s-memory%dMB-img%dMB-IAT%vs-burst%d-st%s-payload%dKB", experiment.Title,
		int(experiment.FunctionMemoryMB), int(experiment.FunctionImageSizeMB), experiment.IATSeconds, experiment.BurstSizes[0],
		experiment.DesiredServiceTimes[0], experiment.PayloadLengthBytes/1024.0)

	directoryPath := filepath.Join(path, detailedTitle)
	log.Infof("[sub-experiment %d] Creating directory at `%s`", experiment.ID, directoryPath)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	latenciesPath := filepath.Join(directoryPath, "latencies.csv")
	log.Infof("[sub-experiment %d] Creating latencies file at `%s`", experiment.ID, latenciesPath)
	latenciesFile, err := os.Create(latenciesPath)
	if err != nil {
		log.Fatalf("[sub-experiment %d] Could not create statistics file: %s", experiment.ID, err.Error())
	}

	statisticsPath := filepath.Join(directoryPath, "statistics.csv")
	log.Infof("[sub-experiment %d] Creating statistics file at `%s`", experiment.ID, statisticsPath)
	statisticsFile, err := os.Create(statisticsPath)
	if err != nil {
		log.Fatalf("[sub-experiment %d] Could not create statistics file: %s", experiment.ID, err.Error())
	}

	if experiment.DataTransferChainLength > 1 {
		dataTransfersPath := filepath.Join(directoryPath, "data-transfers.csv")
		log.Infof("[sub-experiment %d] Creating data transfers file at `%s`", experiment.ID, dataTransfersPath)
		dataTransfersFile, err := os.Create(dataTransfersPath)
		if err != nil {
			log.Fatalf("[sub-experiment %d] Could not create data transfers file: %s", experiment.ID, err.Error())
		}
		return directoryPath, latenciesFile, statisticsFile, dataTransfersFile
	}

	return directoryPath, latenciesFile, statisticsFile, nil
}

func generateIAT(experiment setup.SubExperiment) []time.Duration {
	step := 1.0
	maxStep := experiment.IATSeconds
	runningDelta := math.Min(maxStep, experiment.IATSeconds)

	log.Debugf("[sub-experiment %d] Generating %s IATs", experiment.ID, experiment.IATType)
	burstDeltas := make([]time.Duration, experiment.Bursts)
	for i := range burstDeltas {
		switch experiment.IATType {
		case "stochastic":
			burstDeltas[i] = time.Duration(getSpinTime(experiment.IATSeconds)*1000) * time.Millisecond
		case "deterministic":
			burstDeltas[i] = time.Duration(experiment.IATSeconds) * time.Second
		case "step":
			// TODO: TEST THIS and allow customization for runningDelta & step
			if i == 0 {
				burstDeltas[0] = time.Duration(runningDelta) * time.Second
			} else {
				burstDeltas[i] = time.Duration(math.Min(maxStep, runningDelta)) * time.Second
			}
			runningDelta += step
		default:
			log.Errorf("[sub-experiment %d] Unrecognized inter-arrival time type %s, using default: stochastic", experiment.ID, experiment.IATType)
			burstDeltas[i] = time.Duration(getSpinTime(experiment.IATSeconds)*1000) * time.Millisecond
		}
	}
	return burstDeltas
}

// Use a shifted and scaled exponential distribution to guarantee a minimum sleep time
func getSpinTime(frequencySeconds float64) float64 {
	rateParameter := 1 / math.Log(frequencySeconds) // experimentally deduced formula
	return frequencySeconds + rand.ExpFloat64()/rateParameter
}
