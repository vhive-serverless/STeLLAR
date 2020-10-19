package aws

import (
	"fmt"
	"log"
	"os/exec"
)

func (lambda Interface) RemoveFunction(i int) {
	log.Printf("Removing producer lambda %s-%v", name, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "delete-function", "--function-name",
		fmt.Sprintf("%s-%v", name, i))
	runCommandAndLog(cmd)
}

func (lambda Interface) RemoveAPI(i int, apiID string) {
	log.Printf("Removing API %s-API-%v", name, i)
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "delete-rest-api", "--rest-api-id", apiID)
	runCommandAndLog(cmd)
}
