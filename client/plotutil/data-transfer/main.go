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

package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot/plotutil"
	"log"
	"os"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	maxLatencyPlotted = 500.0
	isStrict          = true
	transferSize      = "1MB"
	iatSec            = 10
	requests          = 3000
	memory            = 10
	memoryUnit        = "GB"
)

func main() {
	plotInstance, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotInstance.Title.Text = fmt.Sprintf("Transfer Size %s, Memory %d%s, IAT ~%dsec, %d indiv. reqs.", transferSize, memory, memoryUnit, iatSec, requests)
	plotInstance.Y.Label.Text = "Portion of requests"
	plotInstance.Y.Min = 0.
	plotInstance.Y.Max = 1.
	plotInstance.X.Label.Text = "Transfer Latency (ms)"
	plotInstance.X.Min = 0.
	plotInstance.X.Max = maxLatencyPlotted

	configFile, err := os.Open("producer-consumer/data-transfers.csv")
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(configFile)

	timestamps0 := df.Col("Function 0 Timestamp").Float()
	timestamps1 := df.Col("Function 1 Timestamp").Float()

	latencies := make([]float64, len(timestamps0))
	for idx := range timestamps0 {
		latencies[idx] = timestamps1[idx] - timestamps0[idx]
	}
	sort.Float64s(latencies)

	if isStrict {
		var maxIndexKept int
		for maxIndexKept = 0; maxIndexKept < len(latencies) && latencies[maxIndexKept] <= maxLatencyPlotted; maxIndexKept++ {
		}
		latencies = latencies[:maxIndexKept]
	}

	latenciesToPlot := make(plotter.XYs, len(latencies))
	for i := 0; i < len(latencies); i++ {
		latenciesToPlot[i].X = latencies[i]
		latenciesToPlot[i].Y = stat.CDF(
			latencies[i],
			stat.Empirical,
			latencies,
			nil,
		)
	}

	err = plotutil.AddLinePoints(plotInstance, latenciesToPlot)
	if err != nil {
		panic(err)
	}

	if err = plotInstance.Save(5*vg.Inch, 5*vg.Inch, fmt.Sprintf("producer-consumer/%s.png", plotInstance.Title.Text)); err != nil {
		panic(err)
	}
}
