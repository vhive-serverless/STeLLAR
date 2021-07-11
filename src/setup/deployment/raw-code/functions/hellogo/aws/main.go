// MIT License
//
// Copyright (c) 2021 Theodor Amariucai
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

package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

//ProducerConsumerResponse is the structure that we expect a consumer-producer function response to follow
type ProducerConsumerResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

func main() {
	lambda.Start(producerConsumer)
}

func producerConsumer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	incrementLimit := extractIncrementLimit(&request)

	simulateWork(incrementLimit)

	reqId := "no-context"

	if ctx != nil {
		// ctx context.Context provides runtime Gateway information
		// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
		lc, _ := lambdacontext.FromContext(ctx)
		reqId = lc.AwsRequestID
	}

	httpOutput, err := json.Marshal(ProducerConsumerResponse{
		RequestID:      reqId,
		TimestampChain: []string{},
	})
	if err != nil {
		log.Fatalf("Could not marshal function output: %s", err)
	}

	return events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      http.StatusOK,
		Body:            string(httpOutput),
	}, nil
}

func extractIncrementLimit(requestHTTP *events.APIGatewayProxyRequest) string {
	var incrementLimit string
	incrementLimit = requestHTTP.QueryStringParameters["IncrementLimit"]
	return incrementLimit
}

//simulateWork will keep the CPU busy-spinning
func simulateWork(incrementLimitString string) {
	incrementLimit, err := strconv.Atoi(incrementLimitString)
	if err != nil {
		log.Fatalf("Could not parse IncrementLimit parameter: %s", err.Error())
	}

	log.Infof("Running function up to increment limit (%d)...", incrementLimit)
	for i := 0; i < incrementLimit; i++ {
	}
}
