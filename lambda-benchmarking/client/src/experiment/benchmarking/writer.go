package benchmarking

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/networking"
	"os"
	"strconv"
	"sync"
	"time"
)

type SafeWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

func InitializeExperimentWriter(file *os.File) *SafeWriter {
	safeExperimentWriter := &SafeWriter{Writer: csv.NewWriter(file)}
	// writer.WriteRowToFile would fail because the instance Initialize was called on didn't have the Writer initialized
	safeExperimentWriter.WriteRowToFile(
		"AWS Request ID",
		"Gateway Endpoint",
		"Sent At",
		"Received At",
		"Client Latency (ms)",
		"Burst ID",
	)
	return safeExperimentWriter
}

func (writer *SafeWriter) GenerateLatencyRecord(gatewayEndpointURL string, requestsWaitGroup *sync.WaitGroup, lambdaIncrementLimit int, payloadLength int, burstId int) {
	defer requestsWaitGroup.Done()
	start := time.Now()

	writer.WriteRowToFile(
		networking.CallAPIGateway(gatewayEndpointURL, lambdaIncrementLimit, payloadLength),
		gatewayEndpointURL,
		start.Format(time.StampNano),
		time.Now().Format(time.StampNano),
		strconv.FormatInt(time.Since(start).Milliseconds(), 10),
		strconv.Itoa(burstId),
	)
}

func (writer *SafeWriter) WriteRowToFile(
	AwsRequestID string,
	gatewayEndpoint string,
	SentAt string,
	ReceivedAt string,
	ClientLatencyMs string,
	BurstID string) {
	writer.mux.Lock()
	if err := writer.Writer.Write([]string{AwsRequestID, gatewayEndpoint, SentAt, ReceivedAt, ClientLatencyMs, BurstID}); err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
