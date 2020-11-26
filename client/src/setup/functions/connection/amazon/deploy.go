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
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/setup/functions/util"
	"strings"
)

const maxFunctionTimeout = 900

//DeployFunction will create a new serverless function in the specified language, with given ID. An API for it will
//then be created, as well as corresponding interactions between them and specific permissions.
func (amazon instance) DeployFunction(language string, memoryAssigned int64) string {
	uniqueID := uuid.New().String()

	APIName := fmt.Sprintf("%s%s", amazon.APINamePrefix, uniqueID)
	apiConfig := amazon.createRESTAPI(APIName)
	resourceID := amazon.getResourceID(APIName, *apiConfig.Id)

	functionName := fmt.Sprintf("%s%s", amazon.LambdaFunctionPrefix, *apiConfig.Id)
	functionConfig := amazon.createFunction(functionName, language, memoryAssigned)

	amazon.createAPIFunctionIntegration(APIName, functionName, *apiConfig.Id, resourceID, *functionConfig.FunctionArn)
	amazon.createAPIDeployment(APIName, *apiConfig.Id)
	amazon.addExecutionPermissions(functionName)

	return *apiConfig.Id
}

func (amazon instance) createRESTAPI(APIName string) *apigateway.RestApi {
	log.Infof("Creating API %s (clone of %s)", APIName, amazon.cloneAPIID)

	createArgs := &apigateway.CreateRestApiInput{
		CloneFrom:             aws.String(amazon.cloneAPIID),
		Description:           aws.String("The API used to access vHive-bench Lambda function with same unique ID."),
		EndpointConfiguration: &apigateway.EndpointConfiguration{Types: aws.StringSlice([]string{"REGIONAL"})},
		Name:                  aws.String(APIName),
	}

	result, err := amazon.apiGatewaySvc.CreateRestApi(createArgs)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.createRESTAPI(APIName)
		}

		log.Fatalf("Cannot create rest API: %s", err.Error())
	}
	log.Debugf("Create rest API result: %s", result.String())

	return result
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

	// Note: `items[1].id` for US, `items[0].id` for EU
	resourceID := *result.Items[1].Id

	log.Infof("RESOURCEID of %s is %s", APIName, resourceID)
	return resourceID
}

func (amazon instance) createFunction(functionName string, language string, memoryAssigned int64) *lambda.FunctionConfiguration {
	log.Infof("Creating producer function %s", functionName)

	var createCode *lambda.FunctionCode
	if amazon.s3Key != "" {
		createCode = &lambda.FunctionCode{
			S3Bucket: aws.String(amazon.s3Bucket),
			S3Key:    aws.String(amazon.s3Key),
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
		Role:          aws.String("arn:aws:iam::335329526041:role/service-role/basic_lambda"),
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
