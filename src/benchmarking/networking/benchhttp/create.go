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

package benchhttp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"stellar/setup"
	"stellar/setup/deployment/connection/amazon"
	"strings"
	"time"
)

// CreateRequest will generate an HTTP request according to the provider passed in the sub-experiment
// configuration object.
func CreateRequest(provider string, payloadLengthBytes int, gatewayEndpoint setup.EndpointInfo, assignedFunctionIncrementLimit int64, storageTransfer bool, route string) *http.Request {
	var request *http.Request

	switch provider {
	case "aws":
		request = createGeneralHttpsRequest(
			http.MethodGet,
			fmt.Sprintf("%s.execute-api.%s.amazonaws.com", gatewayEndpoint.ID, amazon.AWSRegion),
		)

		appendProducerConsumerParameters(provider, request, payloadLengthBytes, assignedFunctionIncrementLimit,
			gatewayEndpoint, storageTransfer, route)

		_, err := amazon.AWSSingletonInstance.RequestSigner.Sign(request, nil, "execute-api", amazon.AWSRegion, time.Now())
		if err != nil {
			log.Fatalf("Could not sign AWS HTTP request: %s", err.Error())
		}
	case "azure":
		// Example Azure Functions URL:
		// stellar.azurewebsites.net/api/hellopy-19?code=2FXks0D4k%2FmEvTc6RNQmfIBa%2FBvN2OPxaxgh4fVVFQbVaencM1PLTw%3D%3D
		request = createGeneralHttpsRequest(
			http.MethodGet,
			fmt.Sprintf("%s.azurewebsites.net", gatewayEndpoint.ID),
		)

		appendProducerConsumerParameters(provider, request, payloadLengthBytes, assignedFunctionIncrementLimit, gatewayEndpoint, storageTransfer, route)
	case "google":
		// Example Google Cloud Functions URL:
		// us-west2-zinc-hour-315914.cloudfunctions.net/hellopy-1
		request = createGeneralHttpsRequest(http.MethodGet, strings.Split(gatewayEndpoint.ID, "/")[0])

		appendProducerConsumerParameters(provider, request, payloadLengthBytes, assignedFunctionIncrementLimit, gatewayEndpoint, storageTransfer, route)
	case "cloudflare":
		fallthrough
	case "gcr":
		request = createGeneralHttpsRequest(http.MethodGet, gatewayEndpoint.ID)

		appendProducerConsumerParameters(provider, request, payloadLengthBytes, assignedFunctionIncrementLimit, gatewayEndpoint, storageTransfer, route)
	case "aliyun":
		// Example Alibaba Cloud URL:
		// http://5cfeb440ed6d4ad69ae29d8408aa606e-ap-southeast-1.alicloudapi.com/foo
		request = createGeneralHttpRequest(
			http.MethodGet,
			fmt.Sprintf("%s-us-west-1.alicloudapi.com", gatewayEndpoint.ID),
		)

		appendProducerConsumerParameters(provider, request, payloadLengthBytes, assignedFunctionIncrementLimit, gatewayEndpoint, storageTransfer, route)
	default:
		return createGeneralHttpsRequest(http.MethodGet, provider)
	}

	return request
}

func createGeneralHttpsRequest(method string, hostname string) *http.Request {
	request, err := http.NewRequest(method, fmt.Sprintf("https://%s", hostname), nil)
	if err != nil {
		log.Fatalf("Could not create HTTPS request: %s", err.Error())
	}
	return request
}

func createGeneralHttpRequest(method string, hostname string) *http.Request {
	request, err := http.NewRequest(method, fmt.Sprintf("http://%s", hostname), nil)
	if err != nil {
		log.Fatalf("Could not create HTTP request: %s", err.Error())
	}
	return request
}
