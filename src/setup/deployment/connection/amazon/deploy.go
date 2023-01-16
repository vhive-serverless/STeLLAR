// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
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

func (instance awsSingleton) DeployFunction(binaryPath string, packageType string, language string, memoryAssigned int64, repurposeIdentifier string) string {
	apiConfig := instance.createRESTAPI()

	functionName := fmt.Sprintf("%s%s_%s", namingPrefix, repurposeIdentifier, *apiConfig.Id)

	functionConfig := instance.createFunction(binaryPath, packageType, functionName, language, memoryAssigned)

	resourceID := instance.getResourceID(*apiConfig.Name, *apiConfig.Id)
	instance.createAPIFunctionIntegration(*apiConfig.Name, functionName, *apiConfig.Id, resourceID, *functionConfig.FunctionArn)
	instance.createAPIDeployment(*apiConfig.Name, *apiConfig.Id)
	instance.addExecutionPermissions(functionName)

	return *apiConfig.Id
}

func (instance awsSingleton) createRESTAPI() *apigateway.RestApi {
	log.Info("Creating REST API...")

	createArgs := &apigateway.CreateRestApiInput{
		Name:                  aws.String("vHive-API"),
		EndpointConfiguration: &apigateway.EndpointConfiguration{Types: aws.StringSlice([]string{"REGIONAL"})},
	}

	result, err := instance.apiGatewaySvc.CreateRestApi(createArgs)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.createRESTAPI()
		}

		log.Fatalf("Cannot create REST API: %s", err.Error())
	}
	log.Debugf("Create REST API result: %s", result.String())

	result, _ = instance.updateAPIWithTemplate(*result.Id)

	return result
}

func (instance awsSingleton) updateAPIWithTemplate(apiID string) (*apigateway.RestApi, error) {
	putAPIArgs := &apigateway.PutRestApiInput{
		Body:      instance.apiTemplateFileContents,
		Mode:      aws.String("merge"),
		RestApiId: aws.String(apiID),
	}

	result, err := instance.apiGatewaySvc.PutRestApi(putAPIArgs)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.updateAPIWithTemplate(apiID)
		}

		log.Fatalf("Cannot update REST API with template: %s", err.Error())
	}
	log.Debugf("Update REST API with template result: %s", result.String())
	return result, nil
}

func (instance awsSingleton) getResourceID(APIName string, apiID string) string {
	args := &apigateway.GetResourcesInput{
		RestApiId: aws.String(apiID),
	}

	result, err := instance.apiGatewaySvc.GetResources(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.getResourceID(APIName, apiID)
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

func (instance awsSingleton) createFunction(binaryPath string, packageType string, functionName string, language string, memoryAssigned int64) *lambda.FunctionConfiguration {
	var lambdaExecutionRole = fmt.Sprintf("arn:aws:iam::%s:role/LambdaProducerConsumer", UserARNNumber)
	log.Infof("Creating producer function %s with role ARN %s", functionName, lambdaExecutionRole)

	var createArgs *lambda.CreateFunctionInput
	switch packageType {
	case "Zip":
		var createCode *lambda.FunctionCode
		if instance.S3Key != "" {
			createCode = &lambda.FunctionCode{
				S3Bucket: aws.String(AWSSingletonInstance.S3Bucket),
				S3Key:    aws.String(instance.S3Key),
			}
		} else {
			createCode = &lambda.FunctionCode{
				ZipFile: instance.localZipFileContents,
			}
		}

		createArgs = &lambda.CreateFunctionInput{
			PackageType:   aws.String(lambda.PackageTypeZip),
			Code:          createCode,
			Description:   aws.String("[Do not modify] Continuous benchmarking function used by STeLLAR."),
			Role:          aws.String(lambdaExecutionRole),
			FunctionName:  aws.String(functionName),
			Handler:       aws.String(binaryPath),
			Runtime:       aws.String(language),
			TracingConfig: &lambda.TracingConfig{Mode: aws.String("PassThrough")},
			Timeout:       aws.Int64(maxFunctionTimeout),
			MemorySize:    aws.Int64(memoryAssigned),
		}
	case "Image":
		createArgs = &lambda.CreateFunctionInput{
			PackageType: aws.String(lambda.PackageTypeImage),
			Code: &lambda.FunctionCode{
				ImageUri: aws.String(instance.ImageURI),
			},
			Description:   aws.String("[Do not modify] Continuous benchmarking function used by STeLLAR."),
			Role:          aws.String(lambdaExecutionRole),
			FunctionName:  aws.String(functionName),
			TracingConfig: &lambda.TracingConfig{Mode: aws.String("PassThrough")},
			Timeout:       aws.Int64(maxFunctionTimeout),
			MemorySize:    aws.Int64(memoryAssigned),
		}
	default:
		log.Fatalf("Package type %s not supported for function creation.", packageType)
	}

	result, err := instance.lambdaSvc.CreateFunction(createArgs)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.createFunction(binaryPath, packageType, functionName, language, memoryAssigned)
		}

		log.Fatalf("Cannot create function: %s", err.Error())
	}
	log.Debugf("Create function result: %s", result.String())

	return result
}

func (instance awsSingleton) createAPIFunctionIntegration(APIName string, functionName string, apiID string, resourceID string, arn string) *apigateway.Integration {
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
			AWSRegion, arn)),
	}

	result, err := instance.apiGatewaySvc.PutIntegration(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.createAPIFunctionIntegration(APIName, functionName, apiID, resourceID, arn)
		}

		log.Fatalf("Cannot put rest API - lambda function integration: %s", err.Error())
	}
	log.Debugf("Put rest API - lambda function integration result: %s", result.String())

	return result
}

func (instance awsSingleton) createAPIDeployment(APIName string, apiID string) *apigateway.Deployment {
	log.Infof("Creating deployment for API %s (stage %s)", APIName, deploymentStage)

	args := &apigateway.CreateDeploymentInput{
		RestApiId: aws.String(apiID),
		StageName: aws.String(deploymentStage),
	}

	result, err := instance.apiGatewaySvc.CreateDeployment(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.createAPIDeployment(APIName, apiID)
		}

		log.Fatalf("Cannot create API deployment: %s", err.Error())
	}
	log.Debugf("Create API deployment result: %s", result.String())

	return result
}

func (instance awsSingleton) addExecutionPermissions(functionName string) *lambda.AddPermissionOutput {
	log.Infof("Adding permissions to execute lambda function %s", functionName)

	args := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: aws.String(functionName),
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String("apigateway-benchmarking"),
	}

	result, err := instance.lambdaSvc.AddPermission(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.addExecutionPermissions(functionName)
		}

		log.Fatalf("Cannot add permission: %s", err.Error())
	}
	log.Debugf("Add permission result: %s", result.String())

	return result
}
