package benchmarking

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/configuration"
	"lambda-benchmarking/client/experiment/networking"
	"net/http"
	"sync"
	"time"
)

func sendBurst(config configuration.ExperimentConfig, burstId int, requests int, gatewayEndpointID string,
	assignedFunctionIncrementLimit int64, safeExperimentWriter *SafeWriter) {
	request := generateRequest(config, gatewayEndpointID, assignedFunctionIncrementLimit)

	log.Infof("Experiment %d: starting burst %d, making %d requests with increment limit %d to (%s).",
		config.Id,
		burstId,
		requests,
		assignedFunctionIncrementLimit,
		request.URL.Hostname(),
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go generateLatencyRecord(&requestsWaitGroup, config.Provider, *request, safeExperimentWriter, burstId)
	}
	requestsWaitGroup.Wait()
	log.Infof("Experiment %d: received all responses for burst %d.", config.Id, burstId)
}

func generateLatencyRecord(requestsWaitGroup *sync.WaitGroup, provider string, request http.Request, safeExperimentWriter *SafeWriter, burstId int) {
	defer requestsWaitGroup.Done()

	startTime := time.Now()
	resp := networking.MakeHTTPRequest(request)
	endTime := time.Now()

	var responseID string
	switch provider {
	case "aws":
		responseID = networking.GetAWSRequestID(resp)
	case "test":
		fallthrough
	default:
		responseID = ""
	}

	safeExperimentWriter.RecordLatencyRecord(request.URL.Hostname(), startTime, endTime, responseID, burstId)
}

func generateRequest(config configuration.ExperimentConfig, gatewayEndpointID string, assignedFunctionIncrementLimit int64) *http.Request {
	switch config.Provider {
	case "aws":
		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com", gatewayEndpointID, networking.Region),
			nil,
		)
		if err != nil {
			log.Error(err)
		}

		request.URL.Path = "/prod/benchmarking"
		request.URL.RawQuery = fmt.Sprintf("LambdaIncrementLimit=%d&PayloadLengthBytes=%d",
			assignedFunctionIncrementLimit,
			config.PayloadLengthBytes,
		)

		_, err = networking.GetAWSSignerSingleton().Sign(request, nil, "execute-api", networking.Region, time.Now())
		if err != nil {
			log.Error(err)
		}

		return request
	case "test":
		request, err := http.NewRequest(http.MethodGet, "https://www.google.com/", nil)
		if err != nil {
			log.Error(err)
		}
		return request
	}

	log.Fatalf("Unrecognized provider %s", config.Provider)
	return nil
}
