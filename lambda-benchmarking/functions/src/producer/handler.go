package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ProducerOutput struct {
	AwsRequestID string `json:"AwsRequestID"`
	Payload      []byte `json:"Payload"`
}

func BenchmarkingProducer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	executionDuration, err := strconv.Atoi(request.QueryStringParameters["ExecMilliseconds"])
	if err != nil {
		return serverError(err)
	}

	time.Sleep(time.Duration(executionDuration) * time.Millisecond)

	randomPayload, err := generatePayload(err, request)
	if err != nil {
		return serverError(err)
	}

	// ctx context.Context provides runtime Gateway information
	// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
	lc, _ := lambdacontext.FromContext(ctx)

	// The APIGatewayProxyResponse.Body fields needs to be a string, so we marshal the Payload into JSON
	output, err := json.Marshal(ProducerOutput{
		Payload:      randomPayload,
		AwsRequestID: lc.AwsRequestID,
	})
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      http.StatusOK,
		Body:            string(output),
	}, nil
}

func generatePayload(err error, request events.APIGatewayProxyRequest) ([]byte, error) {
	payloadLength, err := strconv.Atoi(request.QueryStringParameters["PayloadLengthBytes"])
	if err != nil {
		return nil, err
	}

	randomPayload := make([]byte, payloadLength)
	rand.Read(randomPayload)
	return randomPayload, nil
}

//Note: on AWS, lambda runtimes are rounded up to the nearest 100ms for usage purposes
//Note: CORS is required to call your API from a webpage that isn’t hosted on the same domain
func main() {
	lambda.Start(BenchmarkingProducer)
}

//Add a helper for handling errors. This logs any error to os.Stderr and returns a 500
//Internal Server Error response that the AWS API Gateway understands.
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}
