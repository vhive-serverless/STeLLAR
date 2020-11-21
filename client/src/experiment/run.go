// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
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

package experiment

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/configuration"
	"lambda-benchmarking/client/experiment/benchmarking"
	"lambda-benchmarking/client/experiment/visualization"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//TriggerSubExperiments will run the sub-experiments specified by the passed configuration object. It creates
//a directory for each sub-experiment, as well as separate visualizations and latency files.
func TriggerSubExperiments(config configuration.Configuration, outputDirectoryPath string, specificExperiment int) {
	var experimentsWaitGroup sync.WaitGroup

	switch specificExperiment {
	case -1: // run all experiments
		for experimentIndex := 0; experimentIndex < len(config.SubExperiments); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			go triggerExperiment(&experimentsWaitGroup, config.SubExperiments[experimentIndex], outputDirectoryPath)

			if config.Sequential {
				experimentsWaitGroup.Wait()
			}
		}
	default:
		if specificExperiment < 0 || specificExperiment >= len(config.SubExperiments) {
			log.Fatalf("Parameter `runExperiment` is invalid: %d", specificExperiment)
		}

		experimentsWaitGroup.Add(1)
		go triggerExperiment(&experimentsWaitGroup, config.SubExperiments[specificExperiment], outputDirectoryPath)
	}

	experimentsWaitGroup.Wait()
}

func triggerExperiment(experimentsWaitGroup *sync.WaitGroup, experiment configuration.SubExperiment, outputDirectoryPath string) {
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

func generateIAT(experiment configuration.SubExperiment) []time.Duration {
	step := 1.0
	maxStep := experiment.CooldownSeconds
	runningDelta := math.Min(maxStep, experiment.CooldownSeconds)

	log.Debugf("SubExperiment %d: Generating %s IATs", experiment.ID, experiment.IATType)
	burstDeltas := make([]time.Duration, experiment.Bursts)
	for i := range burstDeltas {
		switch experiment.IATType {
		case "stochastic":
			burstDeltas[i] = time.Duration(getSpinTime(experiment.CooldownSeconds)*1000) * time.Millisecond
		case "deterministic":
			burstDeltas[i] = time.Duration(experiment.CooldownSeconds) * time.Second
		case "step":
			// TODO: TEST THIS and allow customization for runningDelta & step
			if i == 0 {
				burstDeltas[0] = time.Duration(runningDelta) * time.Second
			} else {
				burstDeltas[i] = time.Duration(math.Min(maxStep, runningDelta)) * time.Second
			}
			runningDelta += step
		default:
			log.Errorf("Unrecognized inter-arrival time type %s, using default: stochastic", experiment.IATType)
			burstDeltas[i] = time.Duration(getSpinTime(experiment.CooldownSeconds)*1000) * time.Millisecond
		}
	}
	return burstDeltas
}

// use a shifted and scaled exponential distribution to guarantee a minimum sleep time
func getSpinTime(frequencySeconds float64) float64 {
	rateParameter := 1 / math.Log(frequencySeconds) // experimentally deduced formula
	return frequencySeconds + rand.ExpFloat64()/rateParameter
}
