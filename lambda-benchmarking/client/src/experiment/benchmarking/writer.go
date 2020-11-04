package benchmarking

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
	"time"
)

type SafeWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

//InitializeExperimentWriter will create a new dedicated writer for this experiment as well as write the first header row
//to the given latencies file.
func InitializeExperimentWriter(file *os.File) *SafeWriter {
	log.Debugf("Creating experiment writer to file `%s`", file.Name())
	safeExperimentWriter := &SafeWriter{Writer: csv.NewWriter(file)}
	// writer.writeRowToFile would fail because the instance Initialize was called on didn't have the Writer initialized
	safeExperimentWriter.writeRowToFile(
		"AWS Request ID",
		"Host",
		"Sent At",
		"Received At",
		"Client Latency (ms)",
		"Burst ID",
	)
	return safeExperimentWriter
}

func (writer *SafeWriter) recordLatencyRecord(host string, startTime time.Time, endTime time.Time, responseID string, burstId int) {
	writer.writeRowToFile(
		responseID,
		host,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		strconv.FormatInt(endTime.Sub(startTime).Milliseconds(), 10),
		strconv.Itoa(burstId),
	)
}

func (writer *SafeWriter) writeRowToFile(awsRequestID string, host string, sentAt string, receivedAt string, clientLatencyMs string, burstID string) {
	writer.mux.Lock()
	if err := writer.Writer.Write([]string{awsRequestID, host, sentAt, receivedAt, clientLatencyMs, burstID}); err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
