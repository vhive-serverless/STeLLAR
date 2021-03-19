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
	"github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen"
	log "github.com/sirupsen/logrus"
	"strconv"
)

//GlobalRandomPayload is a 1MB string used for quick random payload generation
var GlobalRandomPayload string

//ProducerConsumerResponse is the structure that we expect a consumer-producer function response to follow
type ProducerConsumerResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

//GenerateResponse creates the HTTP or gRPC producer-consumer response payload
func GenerateResponse(ctx context.Context, requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen.InvokeChainRequest) ([]byte, []string) {
	dataTransferChainIDs, incrementLimit := extractChainIDsAndIncrementLimit(requestHTTP, requestGRPC)

	var updatedTimestampChain []string
	if isFirstFunctionInChain(requestGRPC, requestHTTP) {
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
		stringPayload := GeneratePayloadFromGlobalRandom(payloadLengthBytes)

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

		if isUsingStorage(requestGRPC, requestHTTP) &&
			len(stringPayload) != 0 &&
			functionsLeftInChain(dataTransferChainIDs) {
			saveObjectToStorage(requestHTTP, stringPayload, requestGRPC) // save again for the next function in the chain
		}
	}

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

func extractChainIDsAndIncrementLimit(requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen.InvokeChainRequest) ([]string, string) {
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

func isFirstFunctionInChain(requestGRPC *proto_gen.InvokeChainRequest, requestHTTP *events.APIGatewayProxyRequest) bool {
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

func invokeNextFunction(requestHTTP *events.APIGatewayProxyRequest, updatedTimestampChain []string, dataTransferChainIDs []string, requestGRPC *proto_gen.InvokeChainRequest) []string {
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
