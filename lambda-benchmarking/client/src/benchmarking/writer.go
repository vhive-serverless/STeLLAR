package benchmarking

import (
	"encoding/csv"
	"lambda-benchmarking/client/networking"
	"log"
	"strconv"
	"sync"
	"time"
)

type SafeWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

func (writer *SafeWriter) GenerateLatencyRecord(requestsWaitGroup *sync.WaitGroup, execMilliseconds int, payloadLength int, burstId int) {
	defer requestsWaitGroup.Done()
	start := time.Now()

	writer.WriteRowToFile(
		networking.CallAPIGateway(execMilliseconds, payloadLength),
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
