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

package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	log "github.com/sirupsen/logrus"
	"math/rand"
	proto_gen2 "proto_gen"
	"strconv"
	"strings"
	"time"
)

//GlobalRandomPayload is a 1MB string used for quick random payload generation
var GlobalRandomPayload string

//GenerateResponse creates the HTTP or gRPC producer-consumer response payload
func GenerateResponse(ctx context.Context, requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen2.InvokeChainRequest) ([]byte, []string) {
	var updatedTimestampChain []string
	if firstFunctionInChain(requestGRPC, requestHTTP) {
		var payloadLengthBytesString string
		if requestHTTP != nil {
			payloadLengthBytesString = requestHTTP.QueryStringParameters["PayloadLengthBytes"]
		} else {
			payloadLengthBytesString = requestGRPC.PayloadLengthBytes
		}

		payloadLengthBytes, err := strconv.Atoi(payloadLengthBytesString)
		if err != nil {
			log.Fatalf("Could not parse PayloadLengthBytes: %s", err)
		}

		log.Infof("Generating transfer payload for producer-consumer chain (length %d bytes)", payloadLengthBytes)
		stringPayload := GenerateStringPayload(payloadLengthBytes)

		updatedTimestampChain = AppendTimestampToChain([]string{})

		if isUsingStorage(requestGRPC, requestHTTP) && len(stringPayload) != 0 {
			saveObjectToStorage(requestHTTP, stringPayload, requestGRPC)
		} else {
			log.Info("Using inline JSON, setting the TransferPayload field.")

			if requestHTTP != nil {
				requestHTTP.QueryStringParameters["TransferPayload"] = stringPayload
			} else {
				requestGRPC.TransferPayload = stringPayload
			}
		}
	} else { // not the first function in the chain
		var stringPayload string
		if isUsingStorage(requestGRPC, requestHTTP) && len(stringPayload) != 0 {
			stringPayload = loadObjectFromStorage(requestHTTP, requestGRPC)
		}

		var timestampChainStringForm string
		if requestHTTP != nil {
			timestampChainStringForm = requestHTTP.QueryStringParameters["TimestampChain"]
		} else {
			timestampChainStringForm = requestGRPC.TimestampChain
		}

		//log.Infof("Not the first function in the chain, TimestampChain field is %q.", timestampChainStringForm)
		updatedTimestampChain = AppendTimestampToChain(StringArrayToArrayOfString(timestampChainStringForm))

		if isUsingStorage(requestGRPC, requestHTTP) && len(stringPayload) != 0 {
			saveObjectToStorage(requestHTTP, stringPayload, requestGRPC)
		}
	}

	dataTransferChainIDs, incrementLimit := getChainIDsAndIncrementLimit(requestHTTP, requestGRPC)

	simulateWork(incrementLimit)

	if functionsLeftInChain(dataTransferChainIDs) {
		log.Infof("There are %d functions left in the chain, invoking next one...", len(dataTransferChainIDs))

		updatedTimestampChain = invokeNextFunction(requestHTTP, updatedTimestampChain, dataTransferChainIDs, requestGRPC)
	}

	if requestHTTP != nil {
		// ctx context.Context provides runtime Gateway information
		// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
		lc, _ := lambdacontext.FromContext(ctx)
		httpOutput, err := json.Marshal(ProducerConsumerResponse{
			RequestID:      lc.AwsRequestID,
			TimestampChain: updatedTimestampChain,
		})
		if err != nil {
			log.Fatalf("Could not marshal function output: %s", err)
		}
		return httpOutput, nil
	}

	// gRPC
	return nil, updatedTimestampChain
}

func getChainIDsAndIncrementLimit(requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen2.InvokeChainRequest) ([]string, string) {
	var dataTransferChainIDsString, incrementLimit string
	if requestHTTP != nil {
		incrementLimit = requestHTTP.QueryStringParameters["IncrementLimit"]
		dataTransferChainIDsString = requestHTTP.QueryStringParameters["DataTransferChainIDs"]
	} else {
		incrementLimit = requestGRPC.IncrementLimit
		dataTransferChainIDsString = fmt.Sprintf("%v", requestGRPC.DataTransferChainIDs)
	}
	return StringArrayToArrayOfString(dataTransferChainIDsString), incrementLimit
}

func invokeNextFunction(requestHTTP *events.APIGatewayProxyRequest, updatedTimestampChain []string, dataTransferChainIDs []string, requestGRPC *proto_gen2.InvokeChainRequest) []string {
	if requestHTTP != nil {
		result := invokeNextFunctionAWS(map[string]string{
			"IncrementLimit":       requestHTTP.QueryStringParameters["IncrementLimit"],
			"TimestampChain":       fmt.Sprintf("%v", updatedTimestampChain),
			"TransferPayload":      requestHTTP.QueryStringParameters["TransferPayload"],
			"DataTransferChainIDs": fmt.Sprintf("%v", dataTransferChainIDs[1:]),
		},
			dataTransferChainIDs[0],
		)

		updatedTimestampChain = extractJSONTimestampChain(result)
	} else {
		updatedTimestampChain = invokeNextFunctionGRPC(
			requestGRPC,
			updatedTimestampChain,
			dataTransferChainIDs,
		)
	}
	return updatedTimestampChain
}

func firstFunctionInChain(requestGRPC *proto_gen2.InvokeChainRequest, requestHTTP *events.APIGatewayProxyRequest) bool {
	if requestHTTP != nil {
		_, hasTimestampChain := requestHTTP.QueryStringParameters["TimestampChain"]
		return !hasTimestampChain
	}

	// gRPC
	return requestGRPC.TimestampChain == ""
}

//functionsLeftInChain checks if there are functions left in the chain
func functionsLeftInChain(dataTransferChainIDs []string) bool {
	return len(dataTransferChainIDs) > 0 && dataTransferChainIDs[0] != ""
}

//GenerateStringPayload creates a transfer payload for the producer-consumer chain
func GenerateStringPayload(payloadLengthBytes int) string {
	repeatedRandomPayload := GlobalRandomPayload
	for len(repeatedRandomPayload) < payloadLengthBytes {
		repeatedRandomPayload += GlobalRandomPayload
	}
	return repeatedRandomPayload[:payloadLengthBytes]
}

//InitializeGlobalRandomPayload creates the initial transfer payload to be used for quicker random payload generation
func InitializeGlobalRandomPayload() {
	const (
		allowedChars                 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		globalRandomPayloadSizeBytes = 1024 * 1024
	)

	generatedTransferPayload := make([]byte, globalRandomPayloadSizeBytes)
	for i := range generatedTransferPayload {
		generatedTransferPayload[i] = allowedChars[rand.Intn(len(allowedChars))]
	}

	GlobalRandomPayload = string(generatedTransferPayload)
}

//extractJSONTimestampChain will process raw bytes into a string array of timestamps
func extractJSONTimestampChain(responsePayload []byte) []string {
	var reply map[string]interface{}
	err := json.Unmarshal(responsePayload, &reply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response into map[string]interface{}: %s", err)
	}

	var parsedReply ProducerConsumerResponse
	err = json.Unmarshal([]byte(reply["body"].(string)), &parsedReply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response body into producerConsumerResponse: %s", err)
	}

	return parsedReply.TimestampChain
}

//AppendTimestampToChain will add a new timestamp to the chain
func AppendTimestampToChain(timestampChain []string) []string {
	timestampMilliString := strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)
	return append(timestampChain, timestampMilliString)
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

//ProducerConsumerResponse is the structure that we expect a consumer-producer function response to follow
type ProducerConsumerResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

//StringArrayToArrayOfString will process, e.g., "[14 35 8]" into []string{14, 35, 8}
func StringArrayToArrayOfString(str string) []string {
	str = strings.Split(str, "]")[0]
	str = strings.Split(str, "[")[1]
	return strings.Split(str, " ")
}
