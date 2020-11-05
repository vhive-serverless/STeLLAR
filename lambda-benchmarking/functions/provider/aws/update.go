package aws

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

//UpdateFunction will update the source code of the serverless function with id `i`.
func (lambda Instance) UpdateFunction(i int) {
	log.Infof("Updating producer lambda code %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "update-function-code", "--function-name",
		fmt.Sprintf("%s-%v", lambda.familiarName, i), "--zip-file", fmt.Sprintf("fileb://%s.zip", lambda.familiarName))
	util.RunCommandAndLog(cmd)
}

//UpdateFunctionConfiguration  will update the configuration (e.g. timeout) of the serverless function with id `i`.
func (lambda Instance) UpdateFunctionConfiguration(i int) {
	log.Infof("Updating producer lambda configuration %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "update-function-configuration",
		"--function-name", fmt.Sprintf("%s-%v", lambda.familiarName, i), "--timeout", "900")
	util.RunCommandAndLog(cmd)
}
