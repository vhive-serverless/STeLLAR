package visualization

import (
	"fmt"
	"github.com/go-gota/gota/series"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot/plotutil"
	"lambda-benchmarking/client/experiment/configuration"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func plotBurstLatenciesHistogram(plotPath string, latencySeries series.Series, burstIndex int, duration time.Duration) {
	plotInstance, err := plot.New()
	if err != nil {
		log.Fatalf("Creating a new histogram plot failed with error %s", err.Error())
		return
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
		log.Error(err)
	}

	plotInstance.Add(histogram)
	if err := plotInstance.Save(5*vg.Inch, 5*vg.Inch, plotPath); err != nil {
		panic(err)
	}
}

func plotLatenciesCDF(plotPath string, latencySeries series.Series, config configuration.Experiment) {
	plotInstance, err := plot.New()
	if err != nil {
		log.Fatalf("Creating a new CDF plot failed with error %s", err.Error())
		return
	}

	plotInstance.Title.Text = fmt.Sprintf("Freq ~%vs, Burst sizes %s", config.CooldownSeconds, config.BurstSizes)
	plotInstance.Y.Label.Text = "portion of requests"
	plotInstance.Y.Min = 0.
	plotInstance.Y.Max = 1.
	plotInstance.X.Label.Text = "latency (ms)"
	plotInstance.X.Min = 0.
	plotInstance.X.Max = 2000.0

	latencies := latencySeries.Float()
	sort.Float64s(latencies)

	// Uncomment below for hard X limit
	//var maxIndexKept int
	//for maxIndexKept = 0; maxIndexKept < len(latencies) && latencies[maxIndexKept] <= plotInstance.X.Max; maxIndexKept++ {
	//}
	//latencies = latencies[:maxIndexKept]

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
	if err := plotInstance.Save(5*vg.Inch, 5*vg.Inch, plotPath); err != nil {
		panic(err)
	}
}
