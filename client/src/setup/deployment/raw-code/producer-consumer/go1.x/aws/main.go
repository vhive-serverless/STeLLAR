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
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSDK "github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"net/http"
	"os"
	"vhive-bench/client/setup/deployment/raw-code/producer-consumer/go1.x/common"
)

const namingPrefix = "vHive-bench_"

func main() {
	lambda.Start(producerConsumer)
}

func producerConsumer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var updatedTimestampChain []string
	if timestampChainStringForm, hasChainParameter := request.QueryStringParameters["TimestampChain"]; !hasChainParameter {
		// First function in the chain
		request.QueryStringParameters["TransferPayload"] = string(common.GeneratePayload(request.QueryStringParameters["PayloadLengthBytes"]))
		updatedTimestampChain = common.AppendTimestampToChain([]string{})
	} else {
		updatedTimestampChain = common.AppendTimestampToChain(common.StringArrayToArrayOfString(timestampChainStringForm))
	}

	common.SimulateWork(request.QueryStringParameters["IncrementLimit"])

	dataTransferChainIDs := common.StringArrayToArrayOfString(request.QueryStringParameters["DataTransferChainIDs"])
	if common.FunctionsLeftInChain(dataTransferChainIDs) {
		log.Printf("There are %d functions left in the chain, invoking next one...", len(dataTransferChainIDs))

		responsePayload := invokeNextFunction(map[string]string{
			"IncrementLimit":       request.QueryStringParameters["IncrementLimit"],
			"TimestampChain":       fmt.Sprintf("%v", updatedTimestampChain),
			"TransferPayload":      request.QueryStringParameters["TransferPayload"],
			"DataTransferChainIDs": fmt.Sprintf("%v", dataTransferChainIDs[1:]),
		}, dataTransferChainIDs[0])

		updatedTimestampChain = common.ExtractJSONTimestampChain(responsePayload)
	}

	// ctx context.Context provides runtime Gateway information
	// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
	lc, _ := lambdacontext.FromContext(ctx)
	output, err := json.Marshal(common.ProducerConsumerResponse{
		RequestID:      lc.AwsRequestID,
		TimestampChain: updatedTimestampChain,
	})
	if err != nil {
		log.Fatalf("Could not marshal function output: %s", err)
	}

	return events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      http.StatusOK,
		Body:            string(output),
	}, nil
}

func invokeNextFunction(parameters map[string]string, functionID string) []byte {
	type Payload struct {
		QueryStringParameters map[string]string `json:"queryStringParameters"`
	}
	nextFunctionPayload, err := json.Marshal(Payload{QueryStringParameters: parameters})
	if err != nil {
		log.Fatalf("Could not marshal nextFunctionPayload: %s", err)
	}

	log.Printf("Invoking next function: %s%s", namingPrefix, functionID)
	client := authenticateClient()
	result, err := client.Invoke(&lambdaSDK.InvokeInput{
		FunctionName:   aws.String(fmt.Sprintf("%s%s", namingPrefix, functionID)),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        nextFunctionPayload,
	})
	if err != nil {
		log.Fatalf("Could not invoke lambda: %s", err)
	}

	return result.Payload
}

func authenticateClient() *lambdaSDK.Lambda {
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: &region,
	})
	if err != nil {
		log.Fatalf("Could not create a new session: %s", err)
	}

	return lambdaSDK.New(sess)
}
