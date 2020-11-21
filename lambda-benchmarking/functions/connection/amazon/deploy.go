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
	"functions/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
)

const maxFunctionTimeout = 900

//DeployFunction will create a new serverless function in the specified language, with id `i`. An API for it will
//then be created, as well as corresponding interactions between them and specific permissions.
func (amazon Instance) DeployFunction(id int, language string, memoryAssigned int64) {
	functionConfig := amazon.createFunction(id, language, memoryAssigned)

	apiConfig := amazon.createRESTAPI(id)
	resourceID := amazon.getResourceID(id, *apiConfig.Id)

	amazon.createAPIFunctionIntegration(id, *apiConfig.Id, resourceID, *functionConfig.FunctionArn)
	amazon.createAPIDeployment(id, *apiConfig.Id)
	amazon.addExecutionPermissions(id)
}

func (amazon Instance) createFunction(i int, language string, memoryAssigned int64) *lambda.FunctionConfiguration {
	log.Infof("Creating producer amazon %s-%v", amazon.appName, i)

	//cmd := exec.Command(  "--zip-file" zipLocation,
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
		Role:          aws.String("arn:aws:iam::335329526041:role/service-role/basic_lambda"),
		FunctionName:  aws.String(fmt.Sprintf("%s-%v", amazon.appName, i)),
		Handler:       aws.String(util.BinaryName),
		Runtime:       aws.String(language),
		TracingConfig: &lambda.TracingConfig{Mode: aws.String("PassThrough")},
		Timeout:       aws.Int64(maxFunctionTimeout),
		MemorySize:    aws.Int64(memoryAssigned),
	}

	result, err := amazon.lambdaSvc.CreateFunction(createArgs)
	if err != nil {
		log.Fatalf("Cannot create function: %s", err.Error())
	}
	log.Debugf("Create function result: %s", result.String())

	return result
}

func (amazon Instance) createRESTAPI(i int) *apigateway.RestApi {
	log.Infof("Creating corresponding API %s-API-%v (clone of %s)", amazon.appName, i, amazon.cloneAPIID)

	createArgs := &apigateway.CreateRestApiInput{
		CloneFrom:             aws.String(amazon.cloneAPIID),
		Description:           aws.String(fmt.Sprintf("The API used to access benchmarking Lambda function %v", i)),
		EndpointConfiguration: &apigateway.EndpointConfiguration{Types: aws.StringSlice([]string{"REGIONAL"})},
		Name:                  aws.String(fmt.Sprintf("%s-API-%v", amazon.appName, i)),
	}

	result, err := amazon.apiGatewaySvc.CreateRestApi(createArgs)
	if err != nil {
		log.Fatalf("Cannot create rest API: %s", err.Error())
	}
	log.Debugf("Create rest API result: %s", result.String())

	return result
}

func (amazon Instance) getResourceID(i int, apiID string) string {
	args := &apigateway.GetResourcesInput{
		Embed:     nil,
		Limit:     nil,
		Position:  nil,
		RestApiId: aws.String(apiID),
	}

	result, err := amazon.apiGatewaySvc.GetResources(args)
	if err != nil {
		log.Fatalf("Cannot get API resources: %s", err.Error())
	}
	log.Debugf("Get API resources result: %s", result.String())

	// Note: `items[1].id` for US, `items[0].id` for EU
	resourceID := *result.Items[1].Id

	log.Infof("RESOURCEID of %s-API-%v is %s", amazon.appName, i, resourceID)
	return resourceID
}

func (amazon Instance) createAPIFunctionIntegration(i int, apiID string, resourceID string, arn string) *apigateway.Integration {
	log.Infof("Creating integration between lambda %s-%v and API %s-API-%v", amazon.appName, i, amazon.appName, i)

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
		log.Fatalf("Cannot put rest API - lambda function integration: %s", err.Error())
	}
	log.Debugf("Put rest API - lambda function integration result: %s", result.String())

	return result
}

func (amazon Instance) createAPIDeployment(i int, apiID string) *apigateway.Deployment {
	log.Infof("Creating deployment for API %s-API-%v (stage %s)", amazon.appName, i, amazon.stage)

	args := &apigateway.CreateDeploymentInput{
		RestApiId: aws.String(apiID),
		StageName: aws.String(amazon.stage),
	}

	result, err := amazon.apiGatewaySvc.CreateDeployment(args)
	if err != nil {
		log.Fatalf("Cannot create API deployment: %s", err.Error())
	}
	log.Debugf("Create API deployment result: %s", result.String())

	return result
}

func (amazon Instance) addExecutionPermissions(i int) *lambda.AddPermissionOutput {
	log.Infof("Adding permissions to execute lambda function %s-%v", amazon.appName, i)

	args := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: aws.String(fmt.Sprintf("%s-%v", amazon.appName, i)),
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String("apigateway-benchmarking"),
	}

	result, err := amazon.lambdaSvc.AddPermission(args)
	if err != nil {
		log.Fatalf("Cannot add permission: %s", err.Error())
	}
	log.Debugf("Add permission result: %s", result.String())

	return result
}
