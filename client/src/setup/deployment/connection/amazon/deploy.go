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
	"vhive-bench/client/util"
)

const (
	maxFunctionTimeout  = 900
	lambdaExecutionRole = "arn:aws:iam::335329526041:role/AWSLambdaBasicExectionRole"
)

func (amazon instance) DeployFunction(language string, memoryAssigned int64) string {
	apiConfig := amazon.createRESTAPI()

	functionName := fmt.Sprintf("%s%s", amazon.NamePrefix, *apiConfig.Id)
	functionConfig := amazon.createFunction(functionName, language, memoryAssigned)

	resourceID := amazon.getResourceID(*apiConfig.Name, *apiConfig.Id)
	amazon.createAPIFunctionIntegration(*apiConfig.Name, functionName, *apiConfig.Id, resourceID, *functionConfig.FunctionArn)
	amazon.createAPIDeployment(*apiConfig.Name, *apiConfig.Id)
	amazon.addExecutionPermissions(functionName)

	return *apiConfig.Id
}

func (amazon instance) createRESTAPI() *apigateway.RestApi {
	log.Info("Creating REST API...")

	createArgs := &apigateway.CreateRestApiInput{
		Name:                  aws.String("vHive-API"),
		EndpointConfiguration: &apigateway.EndpointConfiguration{Types: aws.StringSlice([]string{"REGIONAL"})},
	}

	result, err := amazon.apiGatewaySvc.CreateRestApi(createArgs)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.createRESTAPI()
		}

		log.Fatalf("Cannot create REST API: %s", err.Error())
	}
	log.Debugf("Create REST API result: %s", result.String())

	result, _ = amazon.updateAPIWithTemplate(*result.Id)

	return result
}

func (amazon instance) updateAPIWithTemplate(apiID string) (*apigateway.RestApi, error) {
	putAPIArgs := &apigateway.PutRestApiInput{
		Body:      amazon.apiTemplate,
		Mode:      aws.String("merge"),
		RestApiId: aws.String(apiID),
	}

	result, err := amazon.apiGatewaySvc.PutRestApi(putAPIArgs)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.updateAPIWithTemplate(apiID)
		}

		log.Fatalf("Cannot update REST API with template: %s", err.Error())
	}
	log.Debugf("Update REST API with template result: %s", result.String())
	return result, nil
}

func (amazon instance) getResourceID(APIName string, apiID string) string {
	args := &apigateway.GetResourcesInput{
		RestApiId: aws.String(apiID),
	}

	result, err := amazon.apiGatewaySvc.GetResources(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.getResourceID(APIName, apiID)
		}

		log.Fatalf("Cannot get API resources: %s", err.Error())
	}
	log.Debugf("Get API resources result: %s", result.String())

	for _, resource := range (*result).Items {
		if resource.ParentId != nil {
			log.Infof("RESOURCEID of %s is %s", APIName, *resource.Id)
			return *resource.Id
		}
	}

	log.Infof("Could not find RESOURCEID of %s", APIName)
	return ""
}

func (amazon instance) createFunction(functionName string, language string, memoryAssigned int64) *lambda.FunctionConfiguration {
	log.Infof("Creating producer function %s", functionName)

	var createCode *lambda.FunctionCode
	if amazon.S3Key != "" {
		createCode = &lambda.FunctionCode{
			S3Bucket: aws.String(s3Bucket),
			S3Key:    aws.String(amazon.S3Key),
		}
	} else {
		createCode = &lambda.FunctionCode{
			ZipFile: amazon.localZip,
		}
	}

	// Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray.PassThrough otherwise.
	createArgs := &lambda.CreateFunctionInput{
		Code:          createCode,
		Description:   aws.String("Benchmarking function managed and used by vHive-bench."),
		Role:          aws.String(lambdaExecutionRole),
		FunctionName:  aws.String(functionName),
		Handler:       aws.String(util.BinaryName),
		Runtime:       aws.String(language),
		TracingConfig: &lambda.TracingConfig{Mode: aws.String("PassThrough")},
		Timeout:       aws.Int64(maxFunctionTimeout),
		MemorySize:    aws.Int64(memoryAssigned),
	}

	result, err := amazon.lambdaSvc.CreateFunction(createArgs)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.createFunction(functionName, language, memoryAssigned)
		}

		log.Fatalf("Cannot create function: %s", err.Error())
	}
	log.Debugf("Create function result: %s", result.String())

	return result
}

func (amazon instance) createAPIFunctionIntegration(APIName string, functionName string, apiID string, resourceID string, arn string) *apigateway.Integration {
	log.Infof("Creating integration between lambda %s and API %s", APIName, functionName)

	args := &apigateway.PutIntegrationInput{
		HttpMethod:            aws.String("ANY"),
		IntegrationHttpMethod: aws.String("ANY"),
		RequestTemplates: aws.StringMap(map[string]string{
			"application/x-www-form-urlencoded": `{\"body\": $input.json(\"$\`,
		}),
		ResourceId: aws.String(resourceID),
		RestApiId:  aws.String(apiID),
		Type:       aws.String("AWS_PROXY"),
		Uri: aws.String(fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations",
			amazon.region, arn)),
	}

	result, err := amazon.apiGatewaySvc.PutIntegration(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.createAPIFunctionIntegration(APIName, functionName, apiID, resourceID, arn)
		}

		log.Fatalf("Cannot put rest API - lambda function integration: %s", err.Error())
	}
	log.Debugf("Put rest API - lambda function integration result: %s", result.String())

	return result
}

func (amazon instance) createAPIDeployment(APIName string, apiID string) *apigateway.Deployment {
	log.Infof("Creating deployment for API %s (stage %s)", APIName, amazon.stage)

	args := &apigateway.CreateDeploymentInput{
		RestApiId: aws.String(apiID),
		StageName: aws.String(amazon.stage),
	}

	result, err := amazon.apiGatewaySvc.CreateDeployment(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.createAPIDeployment(APIName, apiID)
		}

		log.Fatalf("Cannot create API deployment: %s", err.Error())
	}
	log.Debugf("Create API deployment result: %s", result.String())

	return result
}

func (amazon instance) addExecutionPermissions(functionName string) *lambda.AddPermissionOutput {
	log.Infof("Adding permissions to execute lambda function %s", functionName)

	args := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: aws.String(functionName),
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String("apigateway-benchmarking"),
	}

	result, err := amazon.lambdaSvc.AddPermission(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.addExecutionPermissions(functionName)
		}

		log.Fatalf("Cannot add permission: %s", err.Error())
	}
	log.Debugf("Add permission result: %s", result.String())

	return result
}
