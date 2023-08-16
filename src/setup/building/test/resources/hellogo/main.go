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

type HelloGoResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

func main() {
	lambda.Start(LambdaHandler)
}

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	incrementLimit := extractIncrementLimit(&request)

	simulateWork(incrementLimit)

	reqId := "no-context"

	if ctx != nil {
		// ctx context.Context provides runtime Gateway information
		// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
		lc, _ := lambdacontext.FromContext(ctx)
		reqId = lc.AwsRequestID
	}

	httpOutput, err := json.Marshal(HelloGoResponse{
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

func extractIncrementLimit(requestHTTP *events.APIGatewayProxyRequest) int {
	var incrementLimitString string
	incrementLimitString = requestHTTP.QueryStringParameters["IncrementLimit"]
	incrementLimit, err := strconv.Atoi(incrementLimitString)
	if err != nil {
		log.Warnf("Could not parse IncrementLimit parameter: %s", err.Error())
		incrementLimit = 0
	}
	return incrementLimit
}

// simulateWork will keep the CPU busy-spinning
func simulateWork(incrementLimit int) {
	log.Infof("Running function up to increment limit (%d)...", incrementLimit)
	for i := 0; i < incrementLimit; i++ {
	}
}
