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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//GenerateResponse creates the HTTP or gRPC producer-consumer response payload
func GenerateResponse(ctx context.Context, requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen.InvokeChainRequest) ([]byte, []string) {
	var updatedTimestampChain []string
	if firstFunctionInChain(requestGRPC, requestHTTP) {
		var payloadLengthBytes string
		if requestHTTP != nil {
			payloadLengthBytes = requestHTTP.QueryStringParameters["PayloadLengthBytes"]
		} else {
			payloadLengthBytes = requestGRPC.PayloadLengthBytes
		}

		stringPayload := generateStringPayload(payloadLengthBytes)

		updatedTimestampChain = AppendTimestampToChain([]string{})

		if usingS3(requestGRPC, requestHTTP) {
			if requestHTTP != nil {
				s3key := saveS3Payload(stringPayload, requestHTTP.QueryStringParameters["S3Bucket"])
				requestHTTP.QueryStringParameters["S3Key"] = s3key
			} else {
				s3key := saveS3Payload(stringPayload, requestGRPC.S3Bucket)
				requestGRPC.S3Key = s3key
			}
		} else {
			log.Info("Using direct JSON, setting TransferPayload field.")

			if requestHTTP != nil {
				requestHTTP.QueryStringParameters["TransferPayload"] = stringPayload
			} else {
				requestGRPC.TransferPayload = stringPayload
			}
		}
	} else { // not the first function in the chain
		var stringPayload string
		if usingS3(requestGRPC, requestHTTP) {
			stringPayload = loadS3Object(requestHTTP, requestGRPC)
		}

		var timestampChainStringForm string
		if requestHTTP != nil {
			timestampChainStringForm = requestHTTP.QueryStringParameters["TimestampChain"]
		} else {
			timestampChainStringForm = requestGRPC.TimestampChain
		}

		//log.Infof("Not the first function in the chain, TimestampChain field is %q.", timestampChainStringForm)
		updatedTimestampChain = AppendTimestampToChain(StringArrayToArrayOfString(timestampChainStringForm))

		if usingS3(requestGRPC, requestHTTP) {
			if requestHTTP != nil {
				s3key := saveS3Payload(stringPayload, requestHTTP.QueryStringParameters["S3Bucket"])
				requestHTTP.QueryStringParameters["S3Key"] = s3key
			} else {
				s3key := saveS3Payload(stringPayload, requestGRPC.S3Bucket)
				requestGRPC.S3Key = s3key
			}
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

func getChainIDsAndIncrementLimit(requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen.InvokeChainRequest) ([]string, string) {
	var dataTransferChainIDsString, incrementLimit string
	if requestHTTP != nil {
		incrementLimit = requestHTTP.QueryStringParameters["IncrementLimit"]
		dataTransferChainIDsString = requestHTTP.QueryStringParameters["DataTransferChainIDs"]
	} else {
		incrementLimit = requestGRPC.IncrementLimit
		dataTransferChainIDsString = requestGRPC.DataTransferChainIDs
	}
	return StringArrayToArrayOfString(dataTransferChainIDsString), incrementLimit
}

func loadS3Object(requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen.InvokeChainRequest) string {
	var s3key, s3bucket string
	if requestHTTP != nil {
		s3key = requestHTTP.QueryStringParameters["S3Key"]
		s3bucket = requestHTTP.QueryStringParameters["S3Bucket"]
	} else {
		s3key = requestGRPC.S3Key
		s3bucket = requestGRPC.S3Bucket
	}

	s3svc, _ := authenticateS3Client()
	object, err := s3svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3bucket),
		Key:    aws.String(s3key),
	})
	if err != nil {
		log.Infof("Object %q not found in S3 bucket %q: %s", s3key, s3bucket, err.Error())
	}

	payload, err := ioutil.ReadAll(object.Body)
	if err != nil {
		log.Infof("Error reading object body: %v", err)
		return ""
	}

	return string(payload)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func saveS3Payload(payload string, s3bucket string) string {
	log.Infof(`Using S3, saving transfer payload (~%d bytes) to AWS S3.`, len(payload))
	s3key := fmt.Sprintf("transfer-payload-%s", randStringBytes(20))

	_, s3uploader := authenticateS3Client()

	uploadOutput, err := s3uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3bucket),
		Key:    aws.String(s3key),
		Body:   strings.NewReader(payload),
	})
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", s3key, s3bucket, err.Error())
	}

	log.Infof("Successfully uploaded %q to bucket %q (%s)", s3key, s3bucket, uploadOutput.Location)
	return s3key
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
			fmt.Sprintf("%v", updatedTimestampChain),
			dataTransferChainIDs,
		)
	}
	return updatedTimestampChain
}

func firstFunctionInChain(requestGRPC *proto_gen.InvokeChainRequest, requestHTTP *events.APIGatewayProxyRequest) bool {
	if requestHTTP != nil {
		_, hasTimestampChain := requestHTTP.QueryStringParameters["TimestampChain"]
		return !hasTimestampChain
	}

	// gRPC
	return requestGRPC.TimestampChain == ""
}

func usingS3(requestGRPC *proto_gen.InvokeChainRequest, requestHTTP *events.APIGatewayProxyRequest) bool {
	if requestHTTP != nil {
		value, hasUseS3Field := requestHTTP.QueryStringParameters["UseS3"]

		if !hasUseS3Field {
			return false
		}

		useS3, err := strconv.ParseBool(value)
		if err != nil {
			log.Errorf("Could not parse UseS3: %s", err.Error())
		}

		return useS3
	}

	// gRPC
	return requestGRPC.UseS3 == true
}

//functionsLeftInChain checks if there are functions left in the chain
func functionsLeftInChain(dataTransferChainIDs []string) bool {
	return len(dataTransferChainIDs) > 0 && dataTransferChainIDs[0] != ""
}

//generateStringPayload creates a transfer payload for the producer-consumer chain
func generateStringPayload(payloadLengthBytesString string) string {
	payloadLengthBytes, err := strconv.Atoi(payloadLengthBytesString)
	if err != nil {
		log.Fatalf("Could not parse PayloadLengthBytes: %s", err)
	}

	log.Infof("Generating transfer payload for producer-consumer chain (length %d bytes)", payloadLengthBytes)
	generatedTransferPayload := make([]byte, payloadLengthBytes)
	rand.Read(generatedTransferPayload)

	return string(generatedTransferPayload)
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
