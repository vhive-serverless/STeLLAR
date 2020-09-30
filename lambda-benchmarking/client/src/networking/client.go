package networking

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	timeout = 8 * time.Second
)

func CallAPIGateway(execMilliseconds int, payloadLengthBytes int) string {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/benchmarking?ExecMilliseconds=%d&PayloadLengthBytes=%d",
			CheckAndReturnEnvVar("AWS_API_GATEWAY_ENDPOINT"),
			execMilliseconds,
			payloadLengthBytes,
		),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("x-api-key", CheckAndReturnEnvVar("AWS_API_GATEWAY_KEY"))

	//Increase context deadline for when number of configured requests
	//is large, which usually triggers `dial tcp: i/o timeout`
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}

	return processResponse(resp)
}

type LambdaFunctionResponse struct {
	AwsRequestID string `json:"AwsRequestID"`
	Payload      []byte `json:"Payload"`
}

func processResponse(resp *http.Response) string {
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Gateway response had status %s", resp.Status)
	}


	var lambdaFunctionResponse LambdaFunctionResponse
	if err := json.Unmarshal(bytes, &lambdaFunctionResponse); err != nil {
		log.Fatal(err)
	}
	return lambdaFunctionResponse.AwsRequestID
}

func CheckAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Fatalf("Environment variable %s is not set.", key)
	}
	return envVar
}
