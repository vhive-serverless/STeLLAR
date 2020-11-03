package visualization

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/experiment/configuration"
	"os"
	"path/filepath"
	"time"
)

func GenerateVisualization(visualizationType string, config configuration.ExperimentConfig, deltas []time.Duration, csvFile *os.File, path string) {
	switch visualizationType {
	case "all":
		log.Infof("Experiment %d: creating all visualizations", config.Id)
		generateCDFs(config, csvFile, path)
		generateHistograms(config, csvFile, path, deltas)
	case "CDF":
		log.Infof("Experiment %d: creating CDF visualizations", config.Id)
		generateCDFs(config, csvFile, path)
	case "histogram":
		log.Infof("Experiment %d: creating histograms visualizations (per-burst)", config.Id)
		generateHistograms(config, csvFile, path, deltas)
	case "":
		fallthrough
	default:
		log.Errorf("Experiment %d: unrecognized visualization `%s`, skipping", config.Id, visualizationType)
	}
}

func generateHistograms(config configuration.ExperimentConfig, csvFile *os.File, path string, deltas []time.Duration) {
	log.Debugf("Experiment %d: reading written latencies file %s", config.Id, csvFile.Name())
	latenciesDF := readWrittenLatenciesFile(csvFile)

	log.Debugf("Experiment %d: plotting latencies burst histograms", config.Id)
	for burstIndex := 0; burstIndex < config.Bursts; burstIndex++ {
		burstDF := latenciesDF.Filter(dataframe.F{Colname: "Burst ID", Comparator: series.Eq, Comparando: burstIndex})
		plotBurstLatenciesHistogram(
			filepath.Join(path, fmt.Sprintf("burst%d_delta%v.png", burstIndex, deltas[burstIndex])),
			burstDF.Col("Client Latency (ms)"),
			burstIndex,
			deltas[burstIndex],
		)
	}
}

func generateCDFs(config configuration.ExperimentConfig, csvFile *os.File, path string) {
	log.Debugf("Experiment %d: reading written latencies file %s", config.Id, csvFile.Name())
	latenciesDF := readWrittenLatenciesFile(csvFile)

	log.Debugf("Experiment %d: plotting latencies CDF", config.Id)
	plotLatenciesCDF(
		filepath.Join(path, "empirical_CDF.png"),
		latenciesDF.Col("Client Latency (ms)"),
		config,
	)
}

func readWrittenLatenciesFile(csvFile *os.File) dataframe.DataFrame {
	_, err := csvFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Error(err)
	}

	latenciesDF := dataframe.ReadCSV(csvFile)
	return latenciesDF
}
