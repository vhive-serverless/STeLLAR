package networking

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	timeout = 15 * time.Minute
)

func CallAPIGateway(gatewayEndpoint string, lambdaIncrementLimit int, payloadLengthBytes int) string {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/benchmarking?LambdaIncrementLimit=%d&PayloadLengthBytes=%d",
			gatewayEndpoint,
			lambdaIncrementLimit,
			payloadLengthBytes,
		),
		nil,
	)
	if err != nil {
		log.Error(err)
	}

	req.Header.Add("x-api-key", CheckAndReturnEnvVar("AWS_API_GATEWAY_KEY"))

	//Increase context deadline for when number of configured requests
	//is large, which usually triggers `dial tcp: i/o timeout`
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Error(err)
	}

	return processResponse(resp, gatewayEndpoint)
}

type LambdaFunctionResponse struct {
	AwsRequestID string `json:"AwsRequestID"`
	Payload      []byte `json:"Payload"`
}

func processResponse(resp *http.Response, endpoint string) string {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		log.Errorf("API Gateway response from %s had status %s:\n %s", endpoint, resp.Status, string(bodyBytes))
	}

	var lambdaFunctionResponse LambdaFunctionResponse
	if err := json.Unmarshal(bytes, &lambdaFunctionResponse); err != nil {
		log.Error(err)
	}
	return lambdaFunctionResponse.AwsRequestID
}

func CheckAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Debugf("Environment variable %s is not set.", key)
	}
	return envVar
}
