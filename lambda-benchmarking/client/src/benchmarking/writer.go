package benchmarking

import (
	"encoding/csv"
	"lambda-benchmarking/client/networking"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type SafeWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

var SafeWriterInstance *SafeWriter

func (writer *SafeWriter) Initialize(file *os.File) {
	SafeWriterInstance = &SafeWriter{Writer: csv.NewWriter(file)}
	// writer.WriteRowToFile would fail because the instance Initialize was called on didn't have the Writer initialized
	SafeWriterInstance.WriteRowToFile(
		"AWS Request ID",
		"Sent At",
		"Received At",
		"Client Latency (ms)",
		"Burst ID",
	)
}

func (writer *SafeWriter) GenerateLatencyRecord(gatewayEndpoint string, requestsWaitGroup *sync.WaitGroup, lambdaIncrementLimit int, payloadLength int, burstId int) {
	defer requestsWaitGroup.Done()
	start := time.Now()

	writer.WriteRowToFile(
		networking.CallAPIGateway(gatewayEndpoint, lambdaIncrementLimit, payloadLength),
		start.Format(time.StampNano),
		time.Now().Format(time.StampNano),
		strconv.FormatInt(time.Since(start).Milliseconds(), 10),
		strconv.Itoa(burstId),
	)
}

func (writer *SafeWriter) WriteRowToFile(
	AwsRequestID string,
	SentAt string,
	ReceivedAt string,
	ClientLatencyMs string,
	BurstID string) {
	writer.mux.Lock()
	if err := writer.Writer.Write([]string{AwsRequestID, SentAt, ReceivedAt, ClientLatencyMs, BurstID}); err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
