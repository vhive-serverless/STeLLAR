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
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
	"time"
	"vhive-bench/client/benchmarking/networking/benchgrpc"
	"vhive-bench/client/benchmarking/networking/benchhttp"
	"vhive-bench/client/benchmarking/writers"
	"vhive-bench/client/setup"
	"vhive-bench/client/util"
)

//runSubExperiment will trigger bursts sequentially to each available gateway for a given experiment, then sleep for the
//selected interval, and repeat.
func runSubExperiment(experiment setup.SubExperiment, burstDeltas []time.Duration, provider string, latenciesWriter *writers.RTTLatencyWriter, dataTransferWriter *writers.DataTransferWriter) {
	burstID := 0
	deltaIndex := 0
	for burstID < experiment.Bursts {
		time.Sleep(burstDeltas[deltaIndex])

		// Send one burst to each available gateway (the more gateways used, the faster the experiment)
		for gatewayID := 0; gatewayID < len(experiment.GatewayEndpoints) && burstID < experiment.Bursts; gatewayID++ {
			// Every refresh period, we cycle through burst sizes if they're dynamic i.e. more than 1 element
			incrementLimit := experiment.FunctionIncrementLimits[util.IntegerMin(deltaIndex, len(experiment.FunctionIncrementLimits)-1)]
			burstSize := experiment.BurstSizes[util.IntegerMin(deltaIndex, len(experiment.BurstSizes)-1)]
			sendBurst(provider, experiment, burstID, burstSize, experiment.GatewayEndpoints[gatewayID], incrementLimit, latenciesWriter, dataTransferWriter)
			burstID++
		}

		deltaIndex++
		log.Debugf("[sub-experiment %d] All %d gateways have been used for bursts, flushing and sleeping for %v...", experiment.ID, len(experiment.GatewayEndpoints), burstDeltas[deltaIndex-1])
		latenciesWriter.Writer.Flush()
		if dataTransferWriter != nil {
			dataTransferWriter.Writer.Flush()
		}
	}
}

func sendBurst(provider string, config setup.SubExperiment, burstID int, requests int, gatewayEndpoint setup.GatewayEndpoint,
	incrementLimit int64, latenciesWriter *writers.RTTLatencyWriter, dataTransfersWriter *writers.DataTransferWriter) {

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
			config.PayloadLengthBytes, gatewayEndpoint, config.StorageTransfer)
	}

	requestsWaitGroup.Wait()
	log.Infof("[sub-experiment %d] Received all responses for burst %d.", config.ID, burstID)
}

func executeRequestAndWriteResults(requestsWaitGroup *sync.WaitGroup, provider string, incrementLimit int64,
	latenciesWriter *writers.RTTLatencyWriter, dataTransfersWriter *writers.DataTransferWriter, burstID int,
	payloadLengthBytes int, gatewayEndpoint setup.GatewayEndpoint, S3Transfer bool) {
	defer requestsWaitGroup.Done()

	var reqSentTime, reqReceivedTime time.Time
	var responseID, hostname string
	var timestampChain []string

	switch provider {
	case "vhive":
		var stringArrayTimeStampChain string
		stringArrayTimeStampChain, reqSentTime, reqReceivedTime = benchgrpc.ExecuteRequest(payloadLengthBytes, gatewayEndpoint, incrementLimit, S3Transfer)

		timestampChain = stringArrayToArrayOfString(stringArrayTimeStampChain)
		hostname = gatewayEndpoint.ID
		responseID = "N/A"
	case "aws":
		request := benchhttp.CreateRequest(provider, payloadLengthBytes, gatewayEndpoint, incrementLimit, S3Transfer)

		var respBody []byte
		respBody, reqSentTime, reqReceivedTime = benchhttp.ExecuteRequest(*request)
		response := benchhttp.ExtractProducerConsumerResponse(respBody)

		timestampChain = response.TimestampChain
		hostname = request.URL.Hostname()
		responseID = response.RequestID
	default:
		responseID = ""
	}

	if dataTransfersWriter != nil {
		dataTransfersWriter.WriteDataTransferRow(
			responseID,
			hostname,
			strconv.Itoa(burstID),
			timestampChain...,
		)
	}

	latenciesWriter.WriteRTTLatencyRow(
		responseID,
		hostname,
		reqSentTime.Format(time.RFC3339),
		reqReceivedTime.Format(time.RFC3339),
		strconv.FormatInt(reqReceivedTime.Sub(reqSentTime).Milliseconds(), 10),
		strconv.Itoa(burstID),
	)
}

//stringArrayToArrayOfString will process, e.g., "[14 35 8]" into []string{14, 35, 8}
func stringArrayToArrayOfString(str string) []string {
	str = strings.Split(str, "]")[0]
	str = strings.Split(str, "[")[1]
	return strings.Split(str, " ")
}
