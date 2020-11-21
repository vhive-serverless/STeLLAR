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

package amazon

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
)

//UpdateFunction will update the source code of the serverless function with id `i`.
func (amazon instance) UpdateFunction(id int) *lambda.FunctionConfiguration {
	log.Infof("Updating producer lambda code %s-%v", amazon.appName, id)

	var args *lambda.UpdateFunctionCodeInput
	if amazon.s3Key != "" {
		args = &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(fmt.Sprintf("%s-%v", amazon.appName, id)),
			S3Bucket:     aws.String(amazon.s3Bucket),
			S3Key:        aws.String(amazon.s3Key),
		}
	} else {
		args = &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(fmt.Sprintf("%s-%v", amazon.appName, id)),
			ZipFile:      aws.Uint8ValueSlice(aws.Uint8Slice(amazon.localZip)),
		}
	}

	result, err := amazon.lambdaSvc.UpdateFunctionCode(args)
	if err != nil {
		log.Fatalf("Cannot update function code: %s", err.Error())
	}
	log.Debugf("Update function code result: %s", result.String())

	return result
}

//UpdateFunctionConfiguration  will update the configuration (e.g. timeout) of the serverless function with id `i`.
func (amazon instance) UpdateFunctionConfiguration(id int, assignedMemory int64) {
	log.Infof("Updating producer lambda configuration %s-%v", amazon.appName, id)

	args := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(fmt.Sprintf("%s-%v", amazon.appName, id)),
		MemorySize:   aws.Int64(assignedMemory),
		Timeout:      aws.Int64(600),
	}

	result, err := amazon.lambdaSvc.UpdateFunctionConfiguration(args)
	if err != nil {
		log.Fatalf("Cannot update function configuration: %s", err.Error())
	}
	log.Debugf("Update function configuration result: %s", result.String())
}
