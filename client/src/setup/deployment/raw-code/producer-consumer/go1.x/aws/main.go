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
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	lambda.Start(vhiveBenchProducerConsumer)
}

func vhiveBenchProducerConsumer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	timestampMilli := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	lambdaIncrementLimit, err := strconv.Atoi(request.QueryStringParameters["LambdaIncrementLimit"])
	if err != nil {
		log.Fatalf("Could not parse LambdaIncrementLimit: %s", err)
	}
	log.Printf("Running function up to increment limit (%d)...", lambdaIncrementLimit)
	for i := 0; i < lambdaIncrementLimit; i++ {
	}

	var updatedTimestampChain []string
	if timestampChainStringForm, ok := request.QueryStringParameters["TimestampChain"]; !ok {
		updatedTimestampChain = []string{strconv.FormatInt(timestampMilli, 10)}
	} else {
		updatedTimestampChain = append(stringArrayToArrayOfString(timestampChainStringForm), strconv.FormatInt(timestampMilli, 10))
	}

	dataTransferChainIDs := stringArrayToArrayOfString(request.QueryStringParameters["DataTransferChainIDs"])
	if len(dataTransferChainIDs) > 0 && dataTransferChainIDs[0] != "" {
		log.Printf("More functions (%d) in the chain, invoking next one...", len(dataTransferChainIDs))

		payloadLengthBytes, err := strconv.Atoi(request.QueryStringParameters["PayloadLengthBytes"])
		if err != nil {
			log.Fatalf("Could not parse PayloadLengthBytes: %s", err)
		}

		log.Printf("Generating transfer payload (length: %d bytes)", payloadLengthBytes)
		generatedTransferPayload := make([]byte, payloadLengthBytes)
		rand.Read(generatedTransferPayload)

		response := invokeNextFunction(map[string]string{
			"LambdaIncrementLimit": request.QueryStringParameters["LambdaIncrementLimit"],
			"PayloadLengthBytes":   request.QueryStringParameters["PayloadLengthBytes"],
			"TimestampChain":       fmt.Sprintf("%v", updatedTimestampChain),
			"TransferPayload":      string(generatedTransferPayload),
			"DataTransferChainIDs": fmt.Sprintf("%v", dataTransferChainIDs[1:]),
		}, dataTransferChainIDs[0])

		updatedTimestampChain = extractTimestampChain(response)
	}

	// ctx context.Context provides runtime Gateway information
	// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
	lc, _ := lambdacontext.FromContext(ctx)
	output, err := json.Marshal(lambdaFunctionOutput{
		AwsRequestID:   lc.AwsRequestID,
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

func extractTimestampChain(response *lambdaSDK.InvokeOutput) []string {
	var reply map[string]interface{}
	err := json.Unmarshal(response.Payload, &reply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response into map[string]interface{}: %s", err)
	}

	var parsedReply lambdaFunctionOutput
	err = json.Unmarshal([]byte(reply["body"].(string)), &parsedReply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response body into functionOutput: %s", err)
	}

	return parsedReply.TimestampChain
}

func invokeNextFunction(parameters map[string]string, functionID string) *lambdaSDK.InvokeOutput {
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: &region,
	})
	if err != nil {
		log.Fatalf("Could not create a new session: %s", err)
	}

	client := lambdaSDK.New(sess)

	type Payload struct {
		QueryStringParameters map[string]string `json:"queryStringParameters"`
	}
	nextFunctionPayload, err := json.Marshal(Payload{QueryStringParameters: parameters})
	if err != nil {
		log.Fatalf("Could not marshal nextFunctionPayload: %s", err)
	}

	log.Printf("Invoking next function: vHive_%s", functionID)
	result, err := client.Invoke(&lambdaSDK.InvokeInput{
		FunctionName:   aws.String(fmt.Sprintf("vHive_%s", functionID)),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        nextFunctionPayload,
	})
	if err != nil {
		log.Fatalf("Could not invoke lambda: %s", err)
	}

	return result
}

type lambdaFunctionOutput struct {
	AwsRequestID   string   `json:"AwsRequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

func stringArrayToArrayOfString(str string) []string {
	str = strings.Split(str, "]")[0]
	str = strings.Split(str, "[")[1]
	return strings.Split(str, " ")
}
