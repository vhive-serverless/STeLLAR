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
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSDK "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"os"
)

var sessionInstance *session.Session

func invokeNextFunctionAWS(parameters map[string]string, functionID string) []byte {
	const namingPrefix = "vHive-bench_"

	type Payload struct {
		QueryStringParameters map[string]string `json:"queryStringParameters"`
	}
	nextFunctionPayload, err := json.Marshal(Payload{QueryStringParameters: parameters})
	if err != nil {
		log.Fatalf("Could not marshal nextFunctionPayload: %s", err)
	}

	log.Printf("Invoking next function: %s%s", namingPrefix, functionID)
	lambdaClient := authenticateLambdaClient()
	result, err := lambdaClient.Invoke(&lambdaSDK.InvokeInput{
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

func authenticateS3Client() (*s3.S3, *s3manager.Uploader) {
	if sessionInstance == nil {
		createSessionInstance()
	}

	return s3.New(sessionInstance), s3manager.NewUploader(sessionInstance)
}

func authenticateLambdaClient() *lambdaSDK.Lambda {
	if sessionInstance == nil {
		createSessionInstance()
	}

	return lambdaSDK.New(sessionInstance)
}

func createSessionInstance() {
	region := os.Getenv("AWS_REGION")
	createdSessionInstance, err := session.NewSession(&aws.Config{
		Region: &region,
	})
	if err != nil {
		log.Fatalf("Could not create a new session: %s", err)
	}
	sessionInstance = createdSessionInstance
}
