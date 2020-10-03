package visualization

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ExtractBurstsAndGeneratePlots(burstsNumber int, deltas []time.Duration, csvFile *os.File, path string) {
	_, err := csvFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	df := dataframe.ReadCSV(csvFile)

	for burstIndex := 0; burstIndex < burstsNumber; burstIndex++ {
		burstDF := df.Filter(dataframe.F{Colname: "Burst ID", Comparator: series.Eq, Comparando: burstIndex})
		PlotBurstLatencies(
			filepath.Join(path, fmt.Sprintf("burst%d_delta%v.png", burstIndex, deltas[burstIndex])),
			burstDF.Col("Client Latency (ms)"),
		)
	}
}
