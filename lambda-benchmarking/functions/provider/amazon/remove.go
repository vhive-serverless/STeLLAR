package amazon

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

//RemoveFunction will remove the serverless function with id `i`.
func (amazon Instance) RemoveFunction(i int) {
	log.Infof("Removing producer lambda %s-%v", amazon.appName, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "delete-function", "--function-name",
		fmt.Sprintf("%s-%v", amazon.appName, i))
	util.RunCommandAndLog(cmd)
}

//RemoveAPI will remove the API corresponding to the serverless function with id `i`.
func (amazon Instance) RemoveAPI(i int) {
	log.Infof("Removing API %s-API-%v", amazon.appName, i)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "delete-rest-api", "--rest-api-id",
		amazon.getAPIID(i))
	util.RunCommandAndLog(cmd)
}
