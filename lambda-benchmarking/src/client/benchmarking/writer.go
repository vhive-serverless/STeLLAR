package benchmarking

import (
	"encoding/csv"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type BenchmarkWriter struct {
	writer *csv.Writer
	mux    sync.Mutex
}

func (writer *BenchmarkWriter) generateLatencyRecord(reqChannel chan<- int, id int) {
	log.Printf("Making request to API Gateway with id %d\n", id)

	start := time.Now()
	responseDuration := rand.Intn(250) // Simulate response duration for now
	time.Sleep(time.Duration(responseDuration) * time.Millisecond)

	writer.writeRowToFile(strconv.Itoa(id), strconv.FormatInt(time.Since(start).Milliseconds(), 10) + "ms")

	reqChannel <- id
}

func (writer *BenchmarkWriter) writeRowToFile(column1 string, column2 string) {
	writer.mux.Lock()
	err := writer.writer.Write([]string{column1, column2})
	if err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
