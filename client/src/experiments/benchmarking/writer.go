// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package benchmarking

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
	"time"
)

//SafeWriter makes file-writing safe for concurrent use by using a mutual exclusion lock.
type SafeWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

//NewExperimentWriter will create a new dedicated writer for this experiment as well as write the first header row
//to the given latencies file.
func NewExperimentWriter(file *os.File) *SafeWriter {
	log.Debugf("Creating experiment writer to file `%s`", file.Name())
	safeExperimentWriter := &SafeWriter{Writer: csv.NewWriter(file)}
	// writer.writeRowToFile would fail because the instance Initialize was called on didn't have the Writer initialized
	safeExperimentWriter.writeRowToFile(
		"Request ID",
		"Host",
		"Sent At",
		"Received At",
		"Client Latency (ms)",
		"Burst ID",
	)
	return safeExperimentWriter
}

func (writer *SafeWriter) recordLatencyRecord(host string, startTime time.Time, endTime time.Time, responseID string, burstID int) {
	writer.writeRowToFile(
		responseID,
		host,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		strconv.FormatInt(endTime.Sub(startTime).Milliseconds(), 10),
		strconv.Itoa(burstID),
	)
}

func (writer *SafeWriter) writeRowToFile(awsRequestID string, host string, sentAt string, receivedAt string, clientLatencyMs string, burstID string) {
	writer.mux.Lock()
	if err := writer.Writer.Write([]string{awsRequestID, host, sentAt, receivedAt, clientLatencyMs, burstID}); err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
