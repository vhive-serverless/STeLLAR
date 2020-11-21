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

//RemoveFunction will remove the serverless function with id `i`.
func (amazon Instance) RemoveFunction(i int) *lambda.DeleteFunctionOutput {
	log.Infof("Removing producer lambda %s-%v", amazon.appName, i)

	args := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(fmt.Sprintf("%s-%v", amazon.appName, i)),
	}

	result, err := amazon.lambdaSvc.DeleteFunction(args)
	if err != nil {
		log.Errorf("Cannot remove function: %s", err.Error())
	}
	log.Debugf("Remove function result: %s", result.String())

	return result
}

//RemoveAPI will remove the API corresponding to the serverless function with id `i`.
func (amazon Instance) RemoveAPI(i int) *apigateway.DeleteRestApiOutput {
	apiID := amazon.getAPIID(i)

	if apiID == "" {
		log.Warnf("API %s-API-%v cannot be removed because it was not found", amazon.appName, i)
		return nil
	}

	log.Infof("Removing API %s-API-%v", amazon.appName, i)

	args := &apigateway.DeleteRestApiInput{RestApiId: aws.String(apiID)}

	result, err := amazon.apiGatewaySvc.DeleteRestApi(args)
	if err != nil {
		log.Errorf("Cannot remove REST API: %s", err.Error())
	}
	log.Debugf("Remove REST API result: %s", result.String())

	return result
}

func (amazon Instance) getAPIID(i int) string {
	log.Infof("Getting ID of API %s-API-%v", amazon.appName, i)

	args := &apigateway.GetRestApisInput{
		Limit: aws.Int64(600),
	}

	result, err := amazon.apiGatewaySvc.GetRestApis(args)
	if err != nil {
		log.Fatalf("Cannot get REST APIs: %s", err.Error())
	}
	log.Debugf("Get REST APIs result: %s", result.String())

	for _, item := range result.Items {
		if strings.Compare(*item.Name, fmt.Sprintf("%s-API-%v", amazon.appName, i)) == 0 {
			apiID := *item.Id
			log.Infof("API ID of %s-API-%v is %s", amazon.appName, i, apiID)
			return apiID
		}
	}

	log.Warnf("Could not find API ID of %s-API-%v in any of the results.", amazon.appName, i)
	return ""
}
