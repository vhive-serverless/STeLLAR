package networking

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/configuration"
	"net/http"
	"time"
)

//GenerateRequest will generate an HTTP request according to the provider passed in the sub-experiment
//configuration object.
func GenerateRequest(experiment configuration.SubExperiment, gatewayEndpointID string, assignedFunctionIncrementLimit int64) *http.Request {
	switch experiment.Provider {
	case "aws":
		return generateAWSRequest(experiment, gatewayEndpointID, assignedFunctionIncrementLimit)
	default:
		return generateCustomRequest(experiment.Provider)
	}
}

func generateCustomRequest(hostname string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/", hostname), nil)
	if err != nil {
		log.Error(err)
	}
	return request
}

func generateAWSRequest(config configuration.SubExperiment, gatewayEndpointID string, assignedFunctionIncrementLimit int64) *http.Request {
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com", gatewayEndpointID, awsRegion),
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

	_, err = getAWSSignerSingleton().Sign(request, nil, "execute-api", awsRegion, time.Now())
	if err != nil {
		log.Error(err)
	}

	return request
}

