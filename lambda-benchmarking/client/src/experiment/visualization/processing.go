package visualization

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"io"
	"lambda-benchmarking/client/experiment/configuration"
	"log"
	"os"
	"path/filepath"
	"time"
)

func GenerateVisualization(visualizationType string, config configuration.ExperimentConfig, deltas []time.Duration, csvFile *os.File, path string) {
	_, err := csvFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	df := dataframe.ReadCSV(csvFile)

	if visualizationType == "histogram" {
		for burstIndex := 0; burstIndex < config.Bursts; burstIndex++ {
			burstDF := df.Filter(dataframe.F{Colname: "Burst ID", Comparator: series.Eq, Comparando: burstIndex})
			plotBurstLatenciesHistogram(
				filepath.Join(path, fmt.Sprintf("burst%d_delta%v.png", burstIndex, deltas[burstIndex])),
				burstDF.Col("Client Latency (ms)"),
				burstIndex,
				deltas[burstIndex],
			)
		}
	} else {
		plotLatenciesCDF(
			filepath.Join(path, fmt.Sprintf("empirical_CDF.png")),
			df.Col("Client Latency (ms)"),
			config,
		)
	}
}
