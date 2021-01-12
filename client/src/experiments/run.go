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

package experiments

import (
	"encoding/csv"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"
	"vhive-bench/client/experiments/benchmarking"
	"vhive-bench/client/experiments/visualization"
	"vhive-bench/client/setup"
)

//TriggerSubExperiments will run the sub-experiments specified by the passed configuration object. It creates
//a directory for each sub-experiment, as well as separate visualizations and latency files.
func TriggerSubExperiments(config setup.Configuration, outputDirectoryPath string, specificExperiment int) {
	var experimentsWaitGroup sync.WaitGroup

	switch specificExperiment {
	case -1: // run all experiments
		for experimentIndex := 0; experimentIndex < len(config.SubExperiments); experimentIndex++ {
			experimentsWaitGroup.Add(1)
			go triggerExperiment(&experimentsWaitGroup, config.Provider, config.SubExperiments[experimentIndex], outputDirectoryPath)

			if config.Sequential {
				experimentsWaitGroup.Wait()
			}
		}
	default:
		if specificExperiment < 0 || specificExperiment >= len(config.SubExperiments) {
			log.Fatalf("Parameter `runExperiment` is invalid: %d", specificExperiment)
		}

		experimentsWaitGroup.Add(1)
		go triggerExperiment(&experimentsWaitGroup, config.Provider, config.SubExperiments[specificExperiment], outputDirectoryPath)
	}

	experimentsWaitGroup.Wait()
}

func triggerExperiment(experimentsWaitGroup *sync.WaitGroup, provider string, experiment setup.SubExperiment, outputDirectoryPath string) {
	log.Infof("[sub-experiment %d] Starting...", experiment.ID)
	defer experimentsWaitGroup.Done()

	experimentDirectoryPath, latenciesFile, statisticsFile := createExperimentOutput(outputDirectoryPath, experiment)
	defer latenciesFile.Close()
	defer statisticsFile.Close()

	burstDeltas := generateIAT(experiment)
	benchmarking.RunProfiler(provider, experiment, burstDeltas, benchmarking.NewExperimentWriter(latenciesFile))

	latenciesDF := readLatenciesFromFile(experiment.ID, latenciesFile)

	sortedLatencies := latenciesDF.Col("Client Latency (ms)").Float()
	sort.Float64s(sortedLatencies)

	visualization.Generate(experiment, burstDeltas, latenciesDF, sortedLatencies, experimentDirectoryPath)
	generateStatistics(statisticsFile, experiment.ID, sortedLatencies)

	log.Infof("[sub-experiment %d] Successfully finished.", experiment.ID)
}

func createExperimentOutput(path string, experiment setup.SubExperiment) (string, *os.File, *os.File) {
	directoryPath := filepath.Join(path, experiment.Title)
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

	return directoryPath, latenciesFile, statisticsFile
}

func generateStatistics(file *os.File, experimentID int, sortedLatencies []float64) {
	log.Debugf("[sub-experiment %d] Generating result statistics...", experimentID)

	statisticsWriter := csv.NewWriter(file)

	if err := statisticsWriter.Write([]string{"Count", "Mean", "Standard Deviation", "Min", "25%ile", "50%ile",
		"75%ile", "95%ile", "Max"}); err != nil {
		log.Errorf("[sub-experiment %d] Could not write statistics header to file: %s", experimentID, err.Error())
	}

	if err := statisticsWriter.Write([]string{
		strconv.Itoa(len(sortedLatencies)),
		fmt.Sprintf("%.2f", stat.Mean(sortedLatencies, nil)),
		fmt.Sprintf("%.2f", stat.StdDev(sortedLatencies, nil)),
		fmt.Sprintf("%.2f", stat.Quantile(0, stat.Empirical, sortedLatencies, nil)),
		fmt.Sprintf("%.2f", stat.Quantile(0.25, stat.Empirical, sortedLatencies, nil)),
		fmt.Sprintf("%.2f", stat.Quantile(0.50, stat.Empirical, sortedLatencies, nil)),
		fmt.Sprintf("%.2f", stat.Quantile(0.75, stat.Empirical, sortedLatencies, nil)),
		fmt.Sprintf("%.2f", stat.Quantile(0.95, stat.Empirical, sortedLatencies, nil)),
		fmt.Sprintf("%.2f", stat.Quantile(1, stat.Empirical, sortedLatencies, nil)),
	}); err != nil {
		log.Errorf("[sub-experiment %d] Could not write statistics to file: %s", experimentID, err.Error())
	}

	statisticsWriter.Flush()
}

func readLatenciesFromFile(experimentID int, csvFile *os.File) dataframe.DataFrame {
	log.Debugf("[sub-experiment %d] Reading written latencies from file %s", experimentID, csvFile.Name())

	_, err := csvFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Error(err)
	}

	latenciesDF := dataframe.ReadCSV(csvFile)
	return latenciesDF
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

// use a shifted and scaled exponential distribution to guarantee a minimum sleep time
func getSpinTime(frequencySeconds float64) float64 {
	rateParameter := 1 / math.Log(frequencySeconds) // experimentally deduced formula
	return frequencySeconds + rand.ExpFloat64()/rateParameter
}
