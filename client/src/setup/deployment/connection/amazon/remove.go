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
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (instance awsSingleton) RemoveFunction(uniqueID string) *lambda.DeleteFunctionOutput {
	functionName := fmt.Sprintf("%s%s", instance.NamePrefix, uniqueID)
	log.Infof("Removing producer lambda %s", functionName)

	args := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(functionName),
	}

	result, err := instance.lambdaSvc.DeleteFunction(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.RemoveFunction(uniqueID)
		}

		log.Errorf("Cannot remove function: %s", err.Error())
	}
	log.Debugf("Remove function result: %s", result.String())

	return result
}

//RemoveAPI will remove the API corresponding to the serverless function given ID.
func (instance awsSingleton) RemoveAPI(uniqueID string) *apigateway.DeleteRestApiOutput {
	log.Infof("Removing API %s-API", instance.NamePrefix)

	args := &apigateway.DeleteRestApiInput{RestApiId: aws.String(uniqueID)}

	result, err := instance.apiGatewaySvc.DeleteRestApi(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.RemoveAPI(uniqueID)
		}

		log.Errorf("Cannot remove REST API: %s", err.Error())
	}
	log.Debugf("Remove REST API result: %s", result.String())

	return result
}
