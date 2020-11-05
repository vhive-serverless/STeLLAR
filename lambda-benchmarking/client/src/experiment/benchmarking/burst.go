package benchmarking

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/configuration"
	"lambda-benchmarking/client/experiment/networking"
	"net/http"
	"sync"
	"time"
)

func sendBurst(config configuration.SubExperiment, burstID int, requests int, gatewayEndpointID string,
	assignedFunctionIncrementLimit int64, safeExperimentWriter *SafeWriter) {
	request := networking.GenerateRequest(config, gatewayEndpointID, assignedFunctionIncrementLimit)

	log.Infof("SubExperiment %d: starting burst %d, making %d requests with increment limit %d to (%s).",
		config.ID,
		burstID,
		requests,
		assignedFunctionIncrementLimit,
		request.URL.Hostname(),
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go generateLatencyRecord(&requestsWaitGroup, config.Provider, *request, safeExperimentWriter, burstID)
	}
	requestsWaitGroup.Wait()
	log.Infof("SubExperiment %d: received all responses for burst %d.", config.ID, burstID)
}

func generateLatencyRecord(requestsWaitGroup *sync.WaitGroup, provider string, request http.Request,
	safeExperimentWriter *SafeWriter, burstID int) {
	defer requestsWaitGroup.Done()

	startTime := time.Now()
	respBody := networking.MakeHTTPRequest(request)
	endTime := time.Now()

	var responseID string
	switch provider {
	case "aws":
		responseID = networking.GetAWSRequestID(respBody)
	default:
		responseID = ""
	}

	safeExperimentWriter.recordLatencyRecord(request.URL.Hostname(), startTime, endTime, responseID, burstID)
}
