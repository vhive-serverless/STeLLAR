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

package benchhttp

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"vhive-bench/setup"
	"vhive-bench/setup/deployment/connection/amazon"
)

//ProducerConsumerResponse is the structure holding the response from a producer-consumer function
type ProducerConsumerResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

//ExtractProducerConsumerResponse will process an HTTP response body coming from a producer-consumer function
func ExtractProducerConsumerResponse(respBody []byte) ProducerConsumerResponse {
	var response ProducerConsumerResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		log.Errorf("ExtractProducerConsumerResponse encountered an error: %v", err)
	}
	return response
}

func appendProducerConsumerParameters(provider string, request *http.Request, payloadLengthBytes int,
	assignedFunctionIncrementLimit int64, gatewayEndpoint setup.EndpointInfo, storageTransfer bool) *http.Request {
	switch provider {
	case "aws":
		request.URL.Path = "/prod/benchmarking"
	case "azure":
		// Example Azure Functions URL:
		// vhive-bench.azurewebsites.net/api/hellopy-19?code=2FXks0D4k%2FmEvTc6RNQmfIBa%2FBvN2OPxaxgh4fVVFQbVaencM1PLTw%3D%3D

		path := strings.Split(gatewayEndpoint.ID, request.Host)[1] // path is after the host
		request.URL.Path = strings.Split(path, "?")[0]             // but before the raw query
	case "google":
		// Example Google Cloud Functions URL:
		// us-west2-zinc-hour-315914.cloudfunctions.net/hellopy-1

		request.URL.Path = strings.Split(gatewayEndpoint.ID, request.Host)[1] // path is after the host
		// there is no raw query
	default:
		log.Fatalf("Unrecognized provider %q", provider)
	}

	request.URL.RawQuery = fmt.Sprintf("IncrementLimit=%d&PayloadLengthBytes=%d&DataTransferChainIDs=%v",
		assignedFunctionIncrementLimit,
		payloadLengthBytes,
		gatewayEndpoint.DataTransferChainIDs,
	)

	if provider == "azure" {
		queryCode := strings.Split(gatewayEndpoint.ID, "code=")[1]
		request.URL.RawQuery += fmt.Sprintf("&code=%v", queryCode)
	}

	if storageTransfer {
		request.URL.RawQuery += fmt.Sprintf("&Bucket=%v&StorageTransfer=true", amazon.AWSSingletonInstance.S3Bucket)
	}

	return request
}
