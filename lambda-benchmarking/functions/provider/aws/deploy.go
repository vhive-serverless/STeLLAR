package aws

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//DeployFunction will create a new serverless function in the specified language, with id `i`. An API for it will
//then be created, as well as corresponding interactions between them and specific permissions.
func (lambda Instance) DeployFunction(i int, language string) string {
	lambda.createFunction(i, language)
	arn := lambda.getFunctionARN(i)

	lambda.createRESTAPI(i)
	apiID := lambda.getAPIID(i)
	resourceID := lambda.getResourceID(i, apiID)

	lambda.createAPIFunctionIntegration(i, apiID, resourceID, arn)
	lambda.createAPIDeployment(i, apiID)
	lambda.addExecutionPermissions(i)
	return apiID
}

func (lambda Instance) getFunctionARN(i int) string {
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "list-functions", "--query",
		fmt.Sprintf("Functions[?FunctionName==`%s-%v`].FunctionArn", lambda.familiarName, i), "--output", "text",
		"--region", lambda.region)
	arn := util.RunCommandAndLog(cmd)
	arn, _ = strconv.Unquote(strings.ReplaceAll(strconv.Quote(arn), `\n`, ""))
	log.Infof("ARN of lambda %s-%v is %s", lambda.familiarName, i, arn)
	return arn
}

func (lambda Instance) createFunction(i int, language string) {
	log.Infof("Creating producer lambda %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "create-function", "--function-name",
		fmt.Sprintf("%s-%v", lambda.familiarName, i), "--runtime", language, "--role",
		checkAndReturnEnvVar("AWS_LAMBDA_ROLE"), "--handler", "producer-handler", "--zip-file",
		fmt.Sprintf("fileb://%s.zip", lambda.familiarName), "--tracing-config", "Mode=PassThrough")
	// Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray.PassThrough otherwise.
	util.RunCommandAndLog(cmd)
}

func (lambda Instance) createRESTAPI(i int) {
	log.Infof("Creating corresponding API %s-API-%v (clone of %s)", lambda.familiarName, i, lambda.cloneAPIID)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "create-rest-api", "--name",
		fmt.Sprintf("%s-API-%v", lambda.familiarName, i), "--description",
		fmt.Sprintf("The API used to access benchmarking Lambda function %v", i), "--endpoint-configuration",
		"types=REGIONAL", "--region", lambda.region, "--clone-from", lambda.cloneAPIID)
	util.RunCommandAndLog(cmd)
}

// Note: `items[1].id` for US, `items[0].id` for EU
func (lambda Instance) getResourceID(i int, apiID string) string {
	// items[0].id needed in eu-west-2
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "get-resources", "--rest-api-id",
		apiID, "--query", "items[1].id", "--output", "text", "--region", lambda.region)
	resourceID := util.RunCommandAndLog(cmd)
	resourceID, _ = strconv.Unquote(strings.ReplaceAll(strconv.Quote(resourceID), `\n`, ""))
	log.Infof("RESOURCEID of %s-API-%v is %s", lambda.familiarName, i, resourceID)
	return resourceID
}

func (lambda Instance) createAPIFunctionIntegration(i int, apiID string, resourceID string, arn string) {
	log.Infof("Creating integration between lambda %s-%v and API %s-API-%v", lambda.familiarName, i, lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "put-integration", "--rest-api-id",
		apiID, "--resource-id", resourceID, "--http-method", "ANY", "--type", "AWS_PROXY", "--integration-http-method",
		"ANY", "--uri", fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations",
			lambda.region, arn), "--request-templates",
		`{"application/x-www-form-urlencoded":"{\"body\": $input.json(\"$\")}"}`, "--region", lambda.region)
	util.RunCommandAndLog(cmd)
}

func (lambda Instance) createAPIDeployment(i int, apiID string) {
	log.Infof("Creating deployment for API %s-API-%v (stage %s)", lambda.familiarName, i, lambda.stage)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "create-deployment", "--rest-api-id",
		apiID, "--stage-name", lambda.stage, "--region", lambda.region)
	util.RunCommandAndLog(cmd)
}

func (lambda Instance) addExecutionPermissions(i int) {
	log.Infof("Adding permissions to execute lambda function %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "add-permission", "--function-name",
		fmt.Sprintf("%s-%v", lambda.familiarName, i), "--statement-id", "apigateway-benchmarking", "--action",
		"lambda:InvokeFunction", "--principal", "apigateway.amazonaws.com", "--region", lambda.region)
	util.RunCommandAndLog(cmd)
}

func checkAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Errorf("Environment variable %s is not set.", key)
	}
	return envVar
}
