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
	"strings"
)

func (instance awsSingleton) UpdateFunction(packageType string, uniqueID string) *lambda.FunctionConfiguration {
	functionName := fmt.Sprintf("%s%s", instance.NamePrefix, uniqueID)
	log.Infof("Updating producer lambda code %s", functionName)

	var args *lambda.UpdateFunctionCodeInput
	switch packageType {
	case "Zip":
		if instance.S3Key != "" {
			args = &lambda.UpdateFunctionCodeInput{
				FunctionName: aws.String(functionName),
				S3Bucket:     aws.String(s3Bucket),
				S3Key:        aws.String(instance.S3Key),
			}
		} else {
			args = &lambda.UpdateFunctionCodeInput{
				FunctionName: aws.String(functionName),
				ZipFile:      aws.Uint8ValueSlice(aws.Uint8Slice(instance.localZip)),
			}
		}
	case "Image":
		args = &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(functionName),
			ImageUri:     aws.String(instance.ImageURI),
		}
	default:
		log.Fatalf("Package type %s not supported for function update.", packageType)
	}

	result, err := instance.lambdaSvc.UpdateFunctionCode(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.UpdateFunction(packageType, uniqueID)
		}

		log.Fatalf("Cannot update function code: %s", err.Error())
	}
	log.Debugf("Update function code result: %s", result.String())

	return result
}

//UpdateFunctionConfiguration  will update the configuration (e.g. timeout) of the serverless function with id `i`.
func (instance awsSingleton) UpdateFunctionConfiguration(uniqueID string, assignedMemory int64) {
	functionName := fmt.Sprintf("%s%s", instance.NamePrefix, uniqueID)
	log.Infof("Updating producer lambda configuration %s", functionName)

	args := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
		MemorySize:   aws.Int64(assignedMemory),
		Timeout:      aws.Int64(600),
	}

	result, err := instance.lambdaSvc.UpdateFunctionConfiguration(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			instance.UpdateFunctionConfiguration(uniqueID, assignedMemory)
		}

		log.Fatalf("Cannot update function configuration: %s", err.Error())
	}
	log.Debugf("Update function configuration result: %s", result.String())
}
