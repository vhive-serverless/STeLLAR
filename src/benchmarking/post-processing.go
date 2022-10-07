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
	"encoding/csv"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
	"stellar/benchmarking/visualization"
	"stellar/setup"
)

func postProcessing(experiment setup.SubExperiment, latenciesFile *os.File, burstDeltas []time.Duration, experimentDirectoryPath string, statisticsFile *os.File) {
	log.Debugf("[sub-experiment %d] Reading written latencies from file %s", experiment.ID, latenciesFile.Name())

	_, err := latenciesFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Error(err)
	}

	latenciesDF := dataframe.ReadCSV(latenciesFile)

	sortedLatencies := latenciesDF.Col("Client Latency (ms)").Float()
	sort.Float64s(sortedLatencies)

	visualization.Generate(experiment, burstDeltas, latenciesDF, sortedLatencies, experimentDirectoryPath)
	generateStatistics(statisticsFile, experiment.ID, sortedLatencies)
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
