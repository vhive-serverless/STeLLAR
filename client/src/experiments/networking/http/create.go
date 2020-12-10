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

package http

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//CreateRequest will generate an HTTP request according to the provider passed in the sub-experiment
//configuration object.
func CreateRequest(provider string, payloadLengthBytes int, gatewayEndpointID string, assignedFunctionIncrementLimit int64) *http.Request {
	switch provider {
	case "aws":
		return generateAWSRequest(payloadLengthBytes, gatewayEndpointID, assignedFunctionIncrementLimit)
	default:
		return generateRequest(http.MethodGet, provider)
	}
}

func generateRequest(method string, hostname string) *http.Request {
	request, err := http.NewRequest(method, fmt.Sprintf("https://%s", hostname), nil)
	if err != nil {
		log.Error(err)
	}
	return request
}

func generateAWSRequest(payloadLengthBytes int, gatewayEndpointID string, assignedFunctionIncrementLimit int64) *http.Request {
	request := generateRequest(
		http.MethodPost,
		fmt.Sprintf("%s.execute-api.%s.amazonaws.com", gatewayEndpointID, awsRegion),
	)

	request.URL.Path = "/prod/benchmarking"
	request.URL.RawQuery = fmt.Sprintf("LambdaIncrementLimit=%d&PayloadLengthBytes=%d",
		assignedFunctionIncrementLimit,
		payloadLengthBytes,
	)

	_, err := getAWSSignerSingleton().Sign(request, nil, "execute-api", awsRegion, time.Now())
	if err != nil {
		log.Error(err)
	}

	return request
}
