package aws

import (
	"fmt"
	"functions/util"
	"log"
	"os/exec"
)

func (lambda Interface) UpdateFunction(i int) {
	log.Printf("Updating producer lambda code %s-%v", name, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "update-function-code", "--function-name",
		fmt.Sprintf("%s-%v", name, i), "--zip-file", fmt.Sprintf("fileb://code/%s.zip", name))
	util.RunCommandAndLog(cmd)
}

func (lambda Interface) UpdateFunctionConfiguration(i int) {
	log.Printf("Updating producer lambda configuration %s-%v", name, i)
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "update-function-configuration",
		"--function-name", fmt.Sprintf("%s-%v", name, i), "--timeout", "900")
	util.RunCommandAndLog(cmd)
}
