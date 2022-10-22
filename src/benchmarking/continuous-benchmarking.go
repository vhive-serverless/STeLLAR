// MIT License
//
// Copyright (c) 2022 Dilina Dehigama and EASE Lab
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
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vhive-bench/setup"
)

type StatisticsRecord struct {
	ExperimentType    string `json:"experiment_type"`
	Date              string `json:"date"`
	SubType           string `json:"subtype"`
	Min               string `json:"min"`
	Max               string `json:"max"`
	Median            string `json:"median"`
	TailLatency       string `json:"tail_latency"`
	FirstQuartile     string `json:"first_quartile"`
	ThirdQuartile     string `json:"third_quartile"`
	StandardDeviation string `json:"standard_deviation"`
	PayloadSize       string `json:"payload_size"`
	BurstSize         string `json:"burst_size"`
	IATType           string `json:"IAT_type"`
	Count             int32  `json:"count"`
	Provider          string `json:"provider"`
}

func writeStatisticsToDB(sortedLatencies []float64, experiment setup.SubExperiment) {
	experimentID := experiment.ID
	log.Infof("[sub-experiment %d] Writing statistics to the database", experimentID)
	url := "https://51941s0gs7.execute-api.us-west-1.amazonaws.com/results"
	method := "POST"

	record := StatisticsRecord{
		ExperimentType:    experiment.Title,
		Date:              time.Now().Format("2006-01-02"),
		Min:               fmt.Sprintf("%.2f", stat.Quantile(0, stat.Empirical, sortedLatencies, nil)),
		Max:               fmt.Sprintf("%.2f", stat.Quantile(1, stat.Empirical, sortedLatencies, nil)),
		Median:            fmt.Sprintf("%.2f", stat.Quantile(0.50, stat.Empirical, sortedLatencies, nil)),
		TailLatency:       fmt.Sprintf("%.2f", stat.Quantile(0.99, stat.Empirical, sortedLatencies, nil)),
		FirstQuartile:     fmt.Sprintf("%.2f", stat.Quantile(0.25, stat.Empirical, sortedLatencies, nil)),
		ThirdQuartile:     fmt.Sprintf("%.2f", stat.Quantile(0.75, stat.Empirical, sortedLatencies, nil)),
		StandardDeviation: fmt.Sprintf("%.2f", stat.StdDev(sortedLatencies, nil)),
		PayloadSize:       strconv.Itoa(experiment.PayloadLengthBytes),
		BurstSize:         strings.Trim(strings.Join(strings.Fields(fmt.Sprint(experiment.BurstSizes)), ","), "[]"),
		IATType:           experiment.IATType,
		Count:             int32(len(sortedLatencies)),
		Provider:          "aws",
	}

	jsonRecord, _ := json.Marshal(record)

	payload := strings.NewReader(string(jsonRecord))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatalf("[sub-experiment %d] Could not create HTTP request: %s", experimentID, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("[sub-experiment %d] HTTP request was not successful! : %s", experimentID, err.Error())
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("[sub-experiment %d] Could not read HTTP response body: %s", experimentID, err.Error())
		return
	}
	log.Infof("[sub-experiment %d] Response received: %s", experimentID, string(body))
}
