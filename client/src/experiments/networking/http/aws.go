// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
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

package http

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	log "github.com/sirupsen/logrus"
)

const awsRegion = "us-west-1"

var signerSingleton *v4.Signer

func getAWSSignerSingleton() *v4.Signer {
	if signerSingleton != nil {
		return signerSingleton
	}

	sessionInstance := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))
	signerSingleton = v4.NewSigner(sessionInstance.Config.Credentials)
	return signerSingleton
}

//LambdaFunctionOutput is the structure holding the response from a Lambda function
type LambdaFunctionOutput struct {
	AwsRequestID   string   `json:"AwsRequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

//GetAWSRequestOutput will process an HTTP response body coming from an AWS integration, extracting its ID.
func GetAWSRequestOutput(respBody []byte) LambdaFunctionOutput {
	var lambdaFunctionResponse LambdaFunctionOutput
	if err := json.Unmarshal(respBody, &lambdaFunctionResponse); err != nil {
		log.Error(err)
	}
	return lambdaFunctionResponse
}
