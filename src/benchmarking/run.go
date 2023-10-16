// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
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
	log "github.com/sirupsen/logrus"
	"stellar/benchmarking/networking/benchgrpc"
	"stellar/benchmarking/networking/benchhttp"
	"stellar/benchmarking/writers"
	"stellar/setup"
	"stellar/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

// runSubExperiment will trigger bursts sequentially to each available gateway for a given experiment, then sleep for the
// selected interval, and repeat.
func runSubExperiment(experiment setup.SubExperiment, burstDeltas []time.Duration, provider string, latenciesWriter *writers.RTTLatencyWriter, dataTransferWriter *writers.DataTransferWriter) {
	burstID := 0
	deltaIndex := 0
	errorThreshold := (experiment.Bursts) * (experiment.BurstSizes[util.IntegerMin(deltaIndex, len(experiment.BurstSizes)-1)]) / 10
	errorCount := 0
	mu := sync.Mutex{}
	for burstID < experiment.Bursts {
		time.Sleep(burstDeltas[deltaIndex])
		// Send one burst to each available gateway (the more gateways used, the faster the experiment)
		for gatewayID := 0; gatewayID < len(experiment.Endpoints) && burstID < experiment.Bursts; gatewayID++ {
			// Every refresh period, we cycle through burst sizes if they're dynamic i.e. more than 1 element
			incrementLimit := experiment.BusySpinIncrements[util.IntegerMin(deltaIndex, len(experiment.BusySpinIncrements)-1)]
			burstSize := experiment.BurstSizes[util.IntegerMin(deltaIndex, len(experiment.BurstSizes)-1)]
			log.Infof("%d", len(experiment.Routes))
			sendBurst(provider, experiment, burstID, burstSize, experiment.Endpoints[gatewayID], incrementLimit, latenciesWriter, dataTransferWriter, experiment.Routes[gatewayID], &errorCount, &mu)
			mu.Lock()
			if errorCount > errorThreshold {
				log.Fatalf("Too many errors (%d) occurred, aborting experiment.", errorCount)
			}
			mu.Unlock()
			burstID++
		}

		deltaIndex++
		log.Debugf("[sub-experiment %d] All %d gateways have been used for bursts, flushing and sleeping for %v...", experiment.ID, len(experiment.Endpoints), burstDeltas[deltaIndex-1])
		latenciesWriter.Writer.Flush()
		if dataTransferWriter != nil {
			dataTransferWriter.Writer.Flush()
		}
	}
}

func sendBurst(provider string, config setup.SubExperiment, burstID int, requests int, gatewayEndpoint setup.EndpointInfo,
	incrementLimit int64, latenciesWriter *writers.RTTLatencyWriter, dataTransfersWriter *writers.DataTransferWriter, route string, errorCount *int, mu *sync.Mutex) {

	log.Infof("[sub-experiment %d] Starting burst %d, making %d requests with increment limit %d to gateway with ID %q of provider %q.",
		config.ID,
		burstID,
		requests,
		incrementLimit,
		gatewayEndpoint.ID,
		provider,
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go executeRequestAndWriteResults(&requestsWaitGroup, provider, incrementLimit, latenciesWriter, dataTransfersWriter, burstID,
			config.PayloadLengthBytes, gatewayEndpoint, config.StorageTransfer, route, errorCount, mu)
	}

	requestsWaitGroup.Wait()
	log.Infof("[sub-experiment %d] Received all responses for burst %d.", config.ID, burstID)
}

func executeRequestAndWriteResults(requestsWaitGroup *sync.WaitGroup, provider string, incrementLimit int64,
	latenciesWriter *writers.RTTLatencyWriter, dataTransfersWriter *writers.DataTransferWriter, burstID int,
	payloadLengthBytes int, gatewayEndpoint setup.EndpointInfo, storageTransfer bool, route string, errorCount *int, mu *sync.Mutex) {
	defer requestsWaitGroup.Done()

	var reqSentTime, reqReceivedTime time.Time
	var responseID, hostname string
	var timestampChain []string
	var ok bool

	switch provider {
	case "vhive":
		var stringArrayTimeStampChain string
		stringArrayTimeStampChain, reqSentTime, reqReceivedTime = benchgrpc.ExecuteRequest(payloadLengthBytes, gatewayEndpoint, incrementLimit, storageTransfer)

		timestampChain = stringArrayToArrayOfString(stringArrayTimeStampChain)
		hostname = gatewayEndpoint.ID
		responseID = "N/A"
	case "aws":
		fallthrough
	case "azure":
		fallthrough
	case "gcr":
		fallthrough
	case "cloudflare":
		fallthrough
	case "aliyun":
		fallthrough
	case "google":
		request := benchhttp.CreateRequest(provider, payloadLengthBytes, gatewayEndpoint, incrementLimit, storageTransfer, route)
		log.Debugf("Created HTTP request with URL (%q), Body (%q)", (*request).URL, (*request).Body)

		var respBody []byte
		ok, respBody, reqSentTime, reqReceivedTime = benchhttp.ExecuteRequest(*request)
		if !ok {
			defer mu.Unlock()
			log.Errorf("Request failed, skipping...")
			mu.Lock()
			*errorCount++
			return
		}
		response := benchhttp.ExtractProducerConsumerResponse(respBody)

		timestampChain = response.TimestampChain
		hostname = request.URL.Hostname()
		responseID = response.RequestID
	default:
		log.Fatalf("Unrecognized provider %q, benchmarking module cannot run.", provider)
	}

	if dataTransfersWriter != nil {
		dataTransfersWriter.WriteDataTransferRow(
			responseID,
			hostname,
			strconv.Itoa(burstID),
			timestampChain...,
		)
	}

	clientLatencyMs := strconv.FormatInt(reqReceivedTime.Sub(reqSentTime).Milliseconds(), 10)
	log.Debugf("Received HTTP response after %sms.", clientLatencyMs)
	latenciesWriter.WriteRTTLatencyRow(
		responseID,
		hostname,
		reqSentTime.Format(time.RFC3339),
		reqReceivedTime.Format(time.RFC3339),
		clientLatencyMs,
		strconv.Itoa(burstID),
	)
}

// stringArrayToArrayOfString will process, e.g., "[14 35 8]" into []string{14, 35, 8}
func stringArrayToArrayOfString(str string) []string {
	log.Debugf("stringArrayToArrayOfString argument was %q", str)
	str = strings.Split(str, "]")[0]
	str = strings.Split(str, "[")[1]
	return strings.Split(str, " ")
}
