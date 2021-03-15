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

package visualization

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"vhive-bench/client/setup"
)

const defaultColdThreshold = 300.

//Generate will create plots, charts, histograms etc. according to the
//visualization passed in the sub-experiment configuration object.
func Generate(experiment setup.SubExperiment, deltas []time.Duration, latenciesDF dataframe.DataFrame,
	sortedLatencies []float64, path string) {
	switch experiment.Visualization {
	case "all":
		log.Infof("[sub-experiment %d] Generating all visualizations", experiment.ID)
		generateCDFs(experiment, sortedLatencies, path)
		generateHistograms(experiment, latenciesDF, path, deltas)
		generateBarCharts(experiment, latenciesDF, defaultColdThreshold, path)
	case "bar":
		log.Infof("[sub-experiment %d] Generating burst bar chart visualization", experiment.ID)
		generateBarCharts(experiment, latenciesDF, defaultColdThreshold, path)
	case "cdf":
		log.Infof("[sub-experiment %d] Generating CDF visualization", experiment.ID)
		generateCDFs(experiment, sortedLatencies, path)
	case "histogram":
		log.Infof("[sub-experiment %d] Generating histograms visualizations (per-burst)", experiment.ID)
		generateHistograms(experiment, latenciesDF, path, deltas)
	case "none":
		log.Warnf("[sub-experiment %d] No visualization selected, skipping", experiment.ID)
	default:
		if strings.Contains(experiment.Visualization, "bar") {
			coldThreshold, err := strconv.ParseFloat(strings.Split(experiment.Visualization, "-")[1], 64)
			if err != nil {
				log.Errorf("[sub-experiment %d] Could not parse bar chart threshold latency, using default.", experiment.ID)
				coldThreshold = defaultColdThreshold
			}

			log.Infof("[sub-experiment %d] Generating bar chart visualization (cold threshold %vms)",
				experiment.ID,
				coldThreshold,
			)
			generateBarCharts(experiment, latenciesDF, coldThreshold, path)
		} else {
			log.Errorf("[sub-experiment %d] Unrecognized visualization `%s`, skipping", experiment.ID, experiment.Visualization)
		}
	}
}

func generateBarCharts(experiment setup.SubExperiment, latenciesDF dataframe.DataFrame, coldThreshold float64, path string) {
	log.Debugf("[sub-experiment %d] Plotting characterization bar chart", experiment.ID)
	plotBurstsBarChart(filepath.Join(path, "bursts_characterization.png"), experiment, coldThreshold, latenciesDF)
}

func generateHistograms(experiment setup.SubExperiment, latenciesDF dataframe.DataFrame, path string, deltas []time.Duration) {
	histogramsDirectoryPath := filepath.Join(path, "histograms")
	log.Infof("[sub-experiment %d] Creating directory for histograms at `%s`", experiment.ID, histogramsDirectoryPath)
	if err := os.MkdirAll(histogramsDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	log.Debugf("[sub-experiment %d] Plotting latency histograms for each burst", experiment.ID)
	for burstIndex := 0; burstIndex < experiment.Bursts; burstIndex++ {
		burstDF := latenciesDF.Filter(dataframe.F{Colname: "Burst ID", Comparator: series.Eq, Comparando: burstIndex})
		plotBurstLatenciesHistogram(
			filepath.Join(histogramsDirectoryPath, fmt.Sprintf("burst%d_delta%v.png", burstIndex, deltas[burstIndex])),
			burstDF.Col("Client Latency (ms)").Float(),
			burstIndex,
			deltas[burstIndex],
		)
	}
}

func generateCDFs(config setup.SubExperiment, sortedLatencies []float64, path string) {
	log.Debugf("[sub-experiment %d] Plotting latencies CDF", config.ID)
	plotLatenciesCDF(
		filepath.Join(path, "empirical_CDF.png"),
		sortedLatencies,
		config,
	)
}
