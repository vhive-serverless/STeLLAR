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
	plotInstance.X.Label.Text = "portion of requests"
	plotInstance.Y.Label.Text = "latency (ms)"

	latencies := latencySeries.Float()
	sort.Float64s(latencies)

	latenciesToPlot := make(plotter.XYs, len(latencies))
	for i := 0; i < len(latencies); i++ {
		latenciesToPlot[i].X = stat.CDF(
			latencies[i],
			stat.Empirical,
			latencies,
			nil,
		)
		latenciesToPlot[i].Y = latencies[i]
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
