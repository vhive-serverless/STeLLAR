package amazon

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const maxFunctionTimeout = "900"

//DeployFunction will create a new serverless function in the specified language, with id `i`. An API for it will
//then be created, as well as corresponding interactions between them and specific permissions.
func (amazon Instance) DeployFunction(i int, language string, zipLocation string) (string, string) {
	var memoryAssigned string
	if i < 300 {
		log.Info("Function has index < 300, assigning 128MB")
		memoryAssigned = "128"
	} else {
		log.Info("Function has index >= 300, assigning 1536MB")
		memoryAssigned = "1536"
	}

	amazon.createFunction(i, language, memoryAssigned, zipLocation)
	arn := amazon.getFunctionARN(i)

	amazon.createRESTAPI(i)
	apiID := amazon.getAPIID(i)
	resourceID := amazon.getResourceID(i, apiID)

	amazon.createAPIFunctionIntegration(i, apiID, resourceID, arn)
	amazon.createAPIDeployment(i, apiID)
	amazon.addExecutionPermissions(i)
	return apiID, memoryAssigned
}

// Functions
func (amazon Instance) createFunction(i int, language string, memoryAssigned string, zipLocation string) {
	log.Infof("Creating producer amazon %s-%v", amazon.appName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "create-function", "--function-name",
		fmt.Sprintf("%s-%v", amazon.appName, i), "--runtime", language, "--role",
		checkAndReturnEnvVar("AWS_LAMBDA_ROLE"), "--handler", "producer-handler", "--zip-file",
		zipLocation, "--tracing-config", "Mode=PassThrough",
		"--timeout", maxFunctionTimeout, "--memory-size", memoryAssigned)
	// Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray.PassThrough otherwise.
	util.RunCommandAndLog(cmd)
}

func (amazon Instance) getFunctionARN(i int) string {
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "list-functions", "--query",
		fmt.Sprintf("Functions[?FunctionName==`%s-%v`].FunctionArn", amazon.appName, i), "--output", "text",
		"--region", amazon.region)
	arn := util.RunCommandAndLog(cmd)
	arn, _ = strconv.Unquote(strings.ReplaceAll(strconv.Quote(arn), `\n`, ""))
	log.Infof("ARN of lambda %s-%v is %s", amazon.appName, i, arn)
	return arn
}

// APIs
func (amazon Instance) createRESTAPI(i int) {
	log.Infof("Creating corresponding API %s-API-%v (clone of %s)", amazon.appName, i, amazon.cloneAPIID)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "create-rest-api", "--name",
		fmt.Sprintf("%s-API-%v", amazon.appName, i), "--description",
		fmt.Sprintf("The API used to access benchmarking Lambda function %v", i), "--endpoint-configuration",
		"types=REGIONAL", "--region", amazon.region, "--clone-from", amazon.cloneAPIID)
	util.RunCommandAndLog(cmd)
}

// Note: `items[1].id` for US, `items[0].id` for EU
func (amazon Instance) getResourceID(i int, apiID string) string {
	// items[0].id needed in eu-west-2
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "get-resources", "--rest-api-id",
		apiID, "--query", "items[1].id", "--output", "text", "--region", amazon.region)
	resourceID := util.RunCommandAndLog(cmd)
	resourceID, _ = strconv.Unquote(strings.ReplaceAll(strconv.Quote(resourceID), `\n`, ""))
	log.Infof("RESOURCEID of %s-API-%v is %s", amazon.appName, i, resourceID)
	return resourceID
}

func (amazon Instance) createAPIFunctionIntegration(i int, apiID string, resourceID string, arn string) {
	log.Infof("Creating integration between lambda %s-%v and API %s-API-%v", amazon.appName, i, amazon.appName, i)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "put-integration", "--rest-api-id",
		apiID, "--resource-id", resourceID, "--http-method", "ANY", "--type", "AWS_PROXY", "--integration-http-method",
		"ANY", "--uri", fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations",
			amazon.region, arn), "--request-templates",
		`{"application/x-www-form-urlencoded":"{\"body\": $input.json(\"$\")}"}`, "--region", amazon.region)
	util.RunCommandAndLog(cmd)
}

func (amazon Instance) createAPIDeployment(i int, apiID string) {
	log.Infof("Creating deployment for API %s-API-%v (stage %s)", amazon.appName, i, amazon.stage)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "create-deployment", "--rest-api-id",
		apiID, "--stage-name", amazon.stage, "--region", amazon.region)
	util.RunCommandAndLog(cmd)
}

func (amazon Instance) addExecutionPermissions(i int) {
	log.Infof("Adding permissions to execute lambda function %s-%v", amazon.appName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "add-permission", "--function-name",
		fmt.Sprintf("%s-%v", amazon.appName, i), "--statement-id", "apigateway-benchmarking", "--action",
		"lambda:InvokeFunction", "--principal", "apigateway.amazonaws.com", "--region", amazon.region)
	util.RunCommandAndLog(cmd)
}

func checkAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Errorf("Environment variable %s is not set.", key)
	}
	return envVar
}
