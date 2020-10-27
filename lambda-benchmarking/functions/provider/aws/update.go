package aws

import (
	"fmt"
	"functions/util"
	"log"
	"os/exec"
)

func (lambda Interface) UpdateFunction(i int) {
	log.Printf("Updating producer lambda code %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "update-function-code", "--function-name",
		fmt.Sprintf("%s-%v", lambda.familiarName, i), "--zip-file", fmt.Sprintf("fileb://code/%s.zip", lambda.familiarName))
	util.RunCommandAndLog(cmd)
}

func (lambda Interface) UpdateFunctionConfiguration(i int) {
	log.Printf("Updating producer lambda configuration %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "update-function-configuration",
		"--function-name", fmt.Sprintf("%s-%v", lambda.familiarName, i), "--timeout", "900")
	util.RunCommandAndLog(cmd)
}
