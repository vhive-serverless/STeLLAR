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

func (writer *SafeWriter) GenerateLatencyRecord(reqChannel chan<- int, id int, execMilliseconds int, payloadLength int) {
	start := time.Now()

	log.Printf("Making request with id %d to API Gateway", id)
	writer.WriteRowToFile(
		networking.CallAPIGateway(execMilliseconds, payloadLength),
		strconv.FormatInt(time.Since(start).Milliseconds(), 10),
	)

	reqChannel <- id
}

func (writer *SafeWriter) WriteRowToFile(column1 string, column2 string) {
	writer.mux.Lock()
	if err := writer.Writer.Write([]string{column1, column2}); err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
