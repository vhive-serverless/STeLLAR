package visualization

import (
	"fmt"
	"github.com/go-gota/gota/series"
	"log"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func PlotBurstLatencies(plotPath string, latencySeries series.Series, burstIndex int, duration time.Duration) {
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
	if err := plotInstance.Save(5*vg.Inch, 5*vg.Inch, plotPath); err != nil {
		panic(err)
	}
}
