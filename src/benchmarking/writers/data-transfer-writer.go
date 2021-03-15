// MIT License
//
// Copyright (c) 2021 Theodor Amariucai
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

package writers

import (
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

//DataTransferWriter records serverless data transfer latencies. It is safe for concurrent use as it uses a mutual exclusion lock.
type DataTransferWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

//NewDataTransferWriter will create a new dedicated writer for this experiment as well as write the first header row.
func NewDataTransferWriter(file *os.File, chainLength int) *DataTransferWriter {
	if file == nil { // If experiment doesn't target data transfer, writer can be nil
		return nil
	}

	log.Debugf("Creating experiment writer to file `%s`", file.Name())
	safeExperimentWriter := &DataTransferWriter{Writer: csv.NewWriter(file)}

	timestampTitles := []string{"Function 0 Timestamp"}
	for i := 1; i < chainLength; i++ {
		timestampTitles = append(timestampTitles, fmt.Sprintf("Function %d Timestamp", i))
	}

	safeExperimentWriter.WriteDataTransferRow(
		"Request ID",
		"Host",
		"Burst ID",
		timestampTitles...,
	)

	return safeExperimentWriter
}

//WriteDataTransferRow records a data transfer timestamp chain to disk.
func (writer *DataTransferWriter) WriteDataTransferRow(awsRequestID string, host string, burstID string, timestamps ...string) {
	writer.mux.Lock()
	if err := writer.Writer.Write(append([]string{awsRequestID, host, burstID}, timestamps...)); err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
