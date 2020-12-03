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
	maxLatencyPlotted = 2000.0
	isStrict          = true
	fileName          = "cdf"
	iatSec            = 600
	burstSizesString  = "1"
)

func main() {
	plotInstance, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotInstance.Title.Text = fmt.Sprintf("IAT ~%dsec, Burst sizes [%s]", iatSec, burstSizesString)
	plotInstance.Y.Label.Text = "Portion of requests"
	plotInstance.Y.Min = 0.
	plotInstance.Y.Max = 1.
	plotInstance.X.Label.Text = "Latency (ms)"
	plotInstance.X.Min = 0.
	plotInstance.X.Max = maxLatencyPlotted

	configFile, err := os.Open("latencies.csv")
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(configFile)

	latencies := df.Col("Client Latency (ms)").Float()
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

	if err := plotInstance.Save(5*vg.Inch, 5*vg.Inch, fmt.Sprintf("%s.png", fileName)); err != nil {
		panic(err)
	}
}
