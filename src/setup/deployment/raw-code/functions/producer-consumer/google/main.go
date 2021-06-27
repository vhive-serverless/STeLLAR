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

package p

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	log "github.com/sirupsen/logrus"
	"net/http"
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
func GenerateResponse(ctx context.Context, requestHTTP *http.Request, requestGRPC *InvokeChainRequest) ([]byte, []string) {
	dataTransferChainIDs, incrementLimit := extractChainIDsAndIncrementLimit(requestHTTP, requestGRPC)

	var updatedTimestampChain []string
	if isFirstFunctionInChain(requestGRPC, requestHTTP) {
		var payloadLengthBytesString string
		if requestHTTP != nil {
			payloadLengthBytesString = requestHTTP.URL.Query().Get("PayloadLengthBytes")
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
				log.Printf("TransferPayload was set to %s", stringPayload)
				requestHTTP.URL.RawQuery += fmt.Sprintf("&TransferPayload=%v", stringPayload)
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
			timestampChainStringForm = requestHTTP.URL.Query().Get("TimestampChain")
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
		reqId := "no-context"

		if ctx != nil {
			// ctx context.Context provides runtime Gateway information
			// (https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html)
			lc, _ := lambdacontext.FromContext(ctx)
			reqId = lc.AwsRequestID
		}

		httpOutput, err := json.Marshal(ProducerConsumerResponse{
			RequestID:      reqId,
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

func extractChainIDsAndIncrementLimit(requestHTTP *http.Request, requestGRPC *InvokeChainRequest) ([]string, string) {
	var dataTransferChainIDsString, incrementLimit string
	if requestHTTP != nil {
		incrementLimit = requestHTTP.URL.Query().Get("IncrementLimit")
		dataTransferChainIDsString = requestHTTP.URL.Query().Get("DataTransferChainIDs")
	} else {
		incrementLimit = requestGRPC.IncrementLimit
		dataTransferChainIDsString = fmt.Sprintf("%v", requestGRPC.DataTransferChainIDs)
	}
	return StringArrayToArrayOfString(dataTransferChainIDsString), incrementLimit
}

func isFirstFunctionInChain(requestGRPC *InvokeChainRequest, requestHTTP *http.Request) bool {
	if requestHTTP != nil {
		timestampChain := requestHTTP.URL.Query().Get("TimestampChain")
		return len(timestampChain) == 0
	}

	// gRPC
	return requestGRPC.TimestampChain == ""
}

//functionsLeftInChain checks if there are functions left in the chain
func functionsLeftInChain(dataTransferChainIDs []string) bool {
	return len(dataTransferChainIDs) > 0 && dataTransferChainIDs[0] != ""
}

func invokeNextFunction(requestHTTP *http.Request, updatedTimestampChain []string, dataTransferChainIDs []string, requestGRPC *InvokeChainRequest) []string {
	if requestHTTP != nil {
		result := invokeNextFunctionGoogle(map[string]string{
			"IncrementLimit":       requestHTTP.URL.Query().Get("IncrementLimit"),
			"TimestampChain":       fmt.Sprintf("%v", updatedTimestampChain),
			"TransferPayload":      requestHTTP.URL.Query().Get("TransferPayload"),
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
