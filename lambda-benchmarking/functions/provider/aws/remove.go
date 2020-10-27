package aws

import (
	"fmt"
	"functions/util"
	"log"
	"os/exec"
)

func (lambda Interface) RemoveFunction(i int) {
	log.Printf("Removing producer lambda %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "delete-function", "--function-name",
		fmt.Sprintf("%s-%v", lambda.familiarName, i))
	util.RunCommandAndLog(cmd)
}

func (lambda Interface) RemoveAPI(i int, apiID string) {
	log.Printf("Removing API %s-API-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "delete-rest-api", "--rest-api-id", apiID)
	util.RunCommandAndLog(cmd)
}
