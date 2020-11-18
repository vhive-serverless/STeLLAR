package visualization

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/configuration"
	"os"
	"path/filepath"
	"time"
)

//GenerateVisualization will generate files representing plots, charts etc. according to the
//visualization passed in the sub-experiment configuration object.
func GenerateVisualization(experiment configuration.SubExperiment, deltas []time.Duration, csvFile *os.File, path string) {
	log.Debugf("SubExperiment %d: reading written latencies file %s", experiment.ID, csvFile.Name())
	latenciesDF := readWrittenLatenciesFile(csvFile)

	switch experiment.Visualization {
	case "all":
		log.Infof("SubExperiment %d: generating all visualizations", experiment.ID)
		generateCDFs(experiment, latenciesDF, path)
		generateHistograms(experiment, latenciesDF, path, deltas)
		generateBarCharts(experiment, latenciesDF, path)
	case "bar":
		log.Infof("SubExperiment %d: generating burst bar chart visualization", experiment.ID)
		generateBarCharts(experiment, latenciesDF, path)
	case "cdf":
		log.Infof("SubExperiment %d: generating CDF visualization", experiment.ID)
		generateCDFs(experiment, latenciesDF, path)
	case "histogram":
		log.Infof("SubExperiment %d: generating histograms visualizations (per-burst)", experiment.ID)
		generateHistograms(experiment, latenciesDF, path, deltas)
	case "none":
		log.Warnf("SubExperiment %d: no visualization selected, skipping", experiment.ID)
	default:
		log.Errorf("SubExperiment %d: unrecognized visualization `%s`, skipping", experiment.ID, experiment.Visualization)
	}
}

func generateBarCharts(experiment configuration.SubExperiment, latenciesDF dataframe.DataFrame, path string) {
	log.Debugf("SubExperiment %d: plotting characterization bar chart", experiment.ID)
	plotBurstsBarChart(filepath.Join(path, "bursts_characterization.png"), experiment, latenciesDF)
}

func generateHistograms(config configuration.SubExperiment, latenciesDF dataframe.DataFrame, path string, deltas []time.Duration) {
	histogramsDirectoryPath := filepath.Join(path, "histograms")
	log.Infof("Creating directory for histograms at `%s`", histogramsDirectoryPath)
	if err := os.MkdirAll(histogramsDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	log.Debugf("SubExperiment %d: plotting latency histograms for each burst", config.ID)
	for burstIndex := 0; burstIndex < config.Bursts; burstIndex++ {
		burstDF := latenciesDF.Filter(dataframe.F{Colname: "Burst ID", Comparator: series.Eq, Comparando: burstIndex})
		plotBurstLatenciesHistogram(
			filepath.Join(histogramsDirectoryPath, fmt.Sprintf("burst%d_delta%v.png", burstIndex, deltas[burstIndex])),
			burstDF.Col("Client Latency (ms)").Float(),
			burstIndex,
			deltas[burstIndex],
		)
	}
}

func generateCDFs(config configuration.SubExperiment, latenciesDF dataframe.DataFrame, path string) {
	log.Debugf("SubExperiment %d: plotting latencies CDF", config.ID)
	plotLatenciesCDF(
		filepath.Join(path, "empirical_CDF.png"),
		latenciesDF.Col("Client Latency (ms)").Float(),
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
