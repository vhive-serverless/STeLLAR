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

package visualization

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"strings"
	"time"
	"vhive-bench/setup"
)

func plotBurstsBarChart(plotPath string, experiment setup.SubExperiment, coldThreshold float64, latenciesDF dataframe.DataFrame) {
	plotInstance := plot.New()

	plotInstance.Title.Text = fmt.Sprintf("Bursts Characterization (%vms warm threshold, cooldown ~%vs)",
		coldThreshold, experiment.IATSeconds)
	plotInstance.X.Label.Text = "Burst Sizes (Sequential)"
	plotInstance.Y.Label.Text = "Requests"

	coldResponses := plotter.Values{}
	warmResponses := plotter.Values{}
	for burstIndex := 0; burstIndex < experiment.Bursts; burstIndex++ {
		burstDF := latenciesDF.Filter(dataframe.F{Colname: "Burst ID", Comparator: series.Eq, Comparando: burstIndex})
		burstLatencies := burstDF.Col("Client Latency (ms)").Float()

		// This always generated same proportion of cold/warm, wrong:
		// sort.Float64s(burstLatencies)
		// coldThreshold := stat.Quantile(0.8, stat.Empirical, burstLatencies, nil)

		burstColdResponses := 0
		burstWarmResponses := 0
		for _, burst := range burstLatencies {
			if burst >= coldThreshold {
				burstColdResponses++
			} else {
				burstWarmResponses++
			}
		}

		coldResponses = append(coldResponses, float64(burstColdResponses))
		warmResponses = append(warmResponses, float64(burstWarmResponses))
	}

	w := vg.Points(20)

	barsWarm, err := plotter.NewBarChart(warmResponses, w)
	if err != nil {
		log.Errorf("[sub-experiment %d] Could not plot warm requests bars in bar chart: %s", experiment.ID, err.Error())
		return
	}
	barsWarm.LineStyle.Width = vg.Length(0)
	barsWarm.Color = plotutil.Color(3) // orange

	barsCold, err := plotter.NewBarChart(coldResponses, w)
	if err != nil {
		log.Errorf("[sub-experiment %d] Could not plot cold requests bars in bar chart: %s", experiment.ID, err.Error())
		return
	}
	barsCold.LineStyle.Width = vg.Length(0)
	barsCold.Color = plotutil.Color(2) // light blue
	barsCold.Offset = -w

	plotInstance.Add(barsWarm, barsCold)
	plotInstance.Legend.Add("Warm Requests", barsWarm)
	plotInstance.Legend.Add("Cold Requests", barsCold)
	plotInstance.Legend.Left = true
	plotInstance.Legend.Top = true

	augmentedBurstSizes := experiment.BurstSizes
	for i := experiment.Bursts - len(experiment.BurstSizes); i > 0; i-- {
		augmentedBurstSizes = append(augmentedBurstSizes, experiment.BurstSizes[len(experiment.BurstSizes)-1])
	}
	plotInstance.NominalX(strings.Split(strings.Trim(fmt.Sprint(augmentedBurstSizes), "[]"), " ")...)

	if err := plotInstance.Save(10*vg.Inch, 5*vg.Inch, plotPath); err != nil {
		log.Errorf("[sub-experiment %d] Could not save bar chart: %s", experiment.ID, err.Error())
	}
}

func plotBurstLatenciesHistogram(plotPath string, burstLatencies []float64, burstIndex int, duration time.Duration) {
	plotInstance := plot.New()

	plotInstance.Title.Text = fmt.Sprintf("Burst %v Histogram (%v since last)", burstIndex, duration)
	plotInstance.X.Label.Text = "Latency (ms)"
	plotInstance.Y.Label.Text = "Requests"

	latencies := make(plotter.Values, len(burstLatencies))
	
	copy(latencies, burstLatencies)

	histogram, err := plotter.NewHist(latencies, 1<<5)
	if err != nil {
		log.Error(err)
	}

	plotInstance.Add(histogram)
	if err := plotInstance.Save(5*vg.Inch, 5*vg.Inch, plotPath); err != nil {
		log.Errorf("Could not save bursts histogram: %s", err.Error())
	}
}

func plotLatenciesCDF(plotPath string, sortedLatencies []float64, experiment setup.SubExperiment) {
	plotInstance := plot.New()

	plotInstance.Title.Text = fmt.Sprintf("IAT ~%vs, Burst sizes %v", experiment.IATSeconds, experiment.BurstSizes)
	plotInstance.Y.Label.Text = "Portion of requests"
	plotInstance.Y.Min = 0.
	plotInstance.Y.Max = 1.
	plotInstance.X.Label.Text = "Latency (ms)"
	plotInstance.X.Min = 0.
	plotInstance.X.Max = 2000.0

	// Uncomment below for hard X limit
	//var maxIndexKept int
	//for maxIndexKept = 0; maxIndexKept < len(sortedLatencies) && sortedLatencies[maxIndexKept] <= plotInstance.X.Max; maxIndexKept++ {
	//}
	//sortedLatencies = sortedLatencies[:maxIndexKept]

	latenciesToPlot := make(plotter.XYs, len(sortedLatencies))
	for i := 0; i < len(sortedLatencies); i++ {
		latenciesToPlot[i].X = sortedLatencies[i]
		latenciesToPlot[i].Y = stat.CDF(
			sortedLatencies[i],
			stat.Empirical,
			sortedLatencies,
			nil,
		)
	}

	err := plotutil.AddLinePoints(plotInstance, latenciesToPlot)
	if err != nil {
		log.Errorf("[sub-experiment %d] Could not add line points to CDF plot: %s", experiment.ID, err.Error())
	}

	// Save the plot to a PNG file.
	if err := plotInstance.Save(5*vg.Inch, 5*vg.Inch, plotPath); err != nil {
		log.Errorf("[sub-experiment %d] Could not save CDF plot: %s", experiment.ID, err.Error())
	}
}
