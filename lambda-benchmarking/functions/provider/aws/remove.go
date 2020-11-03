package aws

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func (lambda Interface) RemoveFunction(i int) {
	log.Infof("Removing producer lambda %s-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "delete-function", "--function-name",
		fmt.Sprintf("%s-%v", lambda.familiarName, i))
	util.RunCommandAndLog(cmd)
}

func (lambda Interface) RemoveAPI(i int, apiID string) {
	log.Infof("Removing API %s-API-%v", lambda.familiarName, i)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "delete-rest-api", "--rest-api-id", apiID)
	util.RunCommandAndLog(cmd)
}
