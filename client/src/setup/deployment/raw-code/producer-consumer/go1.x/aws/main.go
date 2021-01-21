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
	currentTimestamp := generateTimestampMilli()

	lambdaIncrementLimit, err := strconv.Atoi(request.QueryStringParameters["LambdaIncrementLimit"])
	if err != nil {
		log.Fatalf("Could not parse LambdaIncrementLimit: %s", err)
	}

	payloadLengthBytes, err := strconv.Atoi(request.QueryStringParameters["PayloadLengthBytes"])
	if err != nil {
		log.Fatalf("Could not parse PayloadLengthBytes: %s", err)
	}

	dataTransferChainIDs := request.QueryStringParameters["DataTransferChainIDs"]

	var timestampChain []string
	if timestampChainStringForm, ok := request.QueryStringParameters["TimestampChain"]; !ok {
		timestampChain = []string{strconv.FormatInt(currentTimestamp, 10)}
	} else {
		timestampChain = append(stringToArrayOfString(timestampChainStringForm), strconv.FormatInt(currentTimestamp, 10))
	}

	log.Printf("Running function up to increment limit (%d)...", lambdaIncrementLimit)
	for i := 0; i < lambdaIncrementLimit; i++ {
	}

	functionsLeftToInvoke := stringToArrayOfString(dataTransferChainIDs)
	if len(functionsLeftToInvoke) > 0 && functionsLeftToInvoke[0] != "" {
		log.Printf("More functions (%d) in the chain, invoking next one...", len(functionsLeftToInvoke))

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
		payload, err := json.Marshal(Payload{QueryStringParameters: map[string]string{
			"LambdaIncrementLimit": strconv.Itoa(lambdaIncrementLimit),
			"PayloadLengthBytes":   strconv.Itoa(payloadLengthBytes),
			"TimestampChain":       fmt.Sprintf("%v", timestampChain),
			"DataTransferChainIDs": fmt.Sprintf("%v", functionsLeftToInvoke[1:]),
		}})
		if err != nil {
			log.Fatalf("Could not marshal payload: %s", err)
		}

		log.Println(fmt.Sprintf("vHive_%s", functionsLeftToInvoke[0]))
		result, err := client.Invoke(&lambdaSDK.InvokeInput{
			FunctionName:   aws.String(fmt.Sprintf("vHive_%s", functionsLeftToInvoke[0])),
			InvocationType: aws.String("RequestResponse"),
			LogType:        aws.String("Tail"),
			Payload:        payload,
		})
		if err != nil {
			log.Fatalf("Could not invoke lambda: %s", err)
		}

		var reply map[string]interface{}
		err = json.Unmarshal(result.Payload, &reply)
		if err != nil {
			log.Fatalf("Could not unmarshal lambda response into map[string]interface{}: %s", err)
		}
		fmt.Println(reply["body"].(string))

		var parsedReply lambdaFunctionOutput
		err = json.Unmarshal([]byte(reply["body"].(string)), &parsedReply)
		if err != nil {
			log.Fatalf("Could not unmarshal lambda response body into functionOutput: %s", err)
		}

		timestampChain = parsedReply.TimestampChain
	}

	// ctx context.Context provides runtime Gateway information
	// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
	lc, _ := lambdacontext.FromContext(ctx)
	output, err := json.Marshal(lambdaFunctionOutput{
		AwsRequestID:   lc.AwsRequestID,
		TimestampChain: timestampChain,
		Payload:        generatePayload(payloadLengthBytes),
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

func generateTimestampMilli() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

type lambdaFunctionOutput struct {
	AwsRequestID   string   `json:"AwsRequestID"`
	TimestampChain []string `json:"TimestampChain"`
	Payload        []byte   `json:"Payload"`
}

func generatePayload(payloadLengthBytes int) []byte {
	log.Printf("Requested payload length: %d bytes.", payloadLengthBytes)

	randomPayload := make([]byte, payloadLengthBytes)
	rand.Read(randomPayload)
	return randomPayload
}

func stringToArrayOfString(str string) []string {
	str = strings.Split(str, "]")[0]
	str = strings.Split(str, "[")[1]
	return strings.Split(str, " ")
}
