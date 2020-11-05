package networking

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"lambda-benchmarking/client/experiment/configuration"
	"net/http"
	"time"
)

const (
	timeout = 15 * time.Minute
)

//MakeHTTPRequest will send an HTTP request, check its status code and return the response.
func MakeHTTPRequest(req http.Request) *http.Response {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Fatalf("HTTP request failed with error %s", err.Error())
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Could not read HTTP response body: %s", err.Error())
		}
		log.Errorf("Response from %s had status %s:\n %s", req.URL.Hostname(), resp.Status, string(bodyBytes))
	}

	return resp
}

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
