package visualization

import (
	"fmt"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot/plotutil"
	"log"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func plotBurstLatenciesHistogram(plotPath string, latencySeries series.Series, burstIndex int, duration time.Duration) {
	plotInstance, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotInstance.Title.Text = fmt.Sprintf("Burst %v Histogram (%v since last)", burstIndex, duration)
	plotInstance.X.Label.Text = "latency (ms)"
	plotInstance.Y.Label.Text = "requests"

	latencies := make(plotter.Values, latencySeries.Len())
	for i := 0; i < latencySeries.Len(); i++ {
		latencies[i] = latencySeries.Float()[i]
	}

	histogram, err := plotter.NewHist(latencies, 1<<5)
	if err != nil {
		log.Fatal(err)
	}

	plotInstance.Add(histogram)
	if err := plotInstance.Save(6*vg.Inch, 6*vg.Inch, plotPath); err != nil {
		panic(err)
	}
}

func plotLatenciesCDF(plotPath string, latencySeries series.Series) {
	plotInstance, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotInstance.Title.Text = fmt.Sprintf("Empirical CDF of Latencies")
	plotInstance.Y.Label.Text = "portion of requests"
	plotInstance.X.Label.Text = "latency (ms)"
	plotInstance.X.Min = 0

	latencies := latencySeries.Float()
	sort.Float64s(latencies)

	plotInstance.X.Max = latencies[len(latencies)-1]

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

	// Save the plot to a PNG file.
	if err := plotInstance.Save(6*vg.Inch, 6*vg.Inch, plotPath); err != nil {
		panic(err)
	}
}
