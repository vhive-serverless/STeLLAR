package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"log"
	"time"
)

// To be exported, field names in event structs must be capitalized
type BenchmarkRequest struct {
	ExecMilliseconds int `json:"ExecMilliseconds"` // "ExecMilliseconds" is a key in this request
}

type BenchmarkResponse struct {
	ProducerOutput []byte `json:"ProducerOutput"` // "ProducerOutput" is a key in this response
}

func BenchmarkingProducer(ctx context.Context, request BenchmarkRequest) (BenchmarkResponse, error) {
	// Note: on AWS, lambda runtimes are rounded up to the nearest 100ms for usage purposes
	executionTime := time.Duration(request.ExecMilliseconds) * time.Millisecond

	//ctx context.Context provides runtime information for your Lambda
	//function invocation (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf(`Received client request with AwsRequestID %s, InvokedFunctionArn %s, 
		and desired execution time %d`, lc.AwsRequestID, lc.InvokedFunctionArn, executionTime)

	executionTimer := time.NewTimer(executionTime)
	<-executionTimer.C
	return BenchmarkResponse{ProducerOutput: []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")}, nil
}

func main() {
	lambda.Start(BenchmarkingProducer)
}
