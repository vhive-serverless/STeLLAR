package amazon

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

//UpdateFunction will update the source code of the serverless function with id `i`.
func (amazon Instance) UpdateFunction(i int, zipLocation string) {
	log.Infof("Updating producer lambda code %s-%v", amazon.appName, i)
	var cmd *exec.Cmd
	if strings.Contains(zipLocation, "fileb") {
		cmd = exec.Command("/usr/local/bin/aws", "lambda", "update-function-code", "--function-name",
			fmt.Sprintf("%s-%v", amazon.appName, i), "--zip-file", zipLocation)
	} else {
		cmd = exec.Command("/usr/local/bin/aws", "lambda", "update-function-code", "--function-name",
			fmt.Sprintf("%s-%v", amazon.appName, i), "--s3-bucket", util.S3Bucket, "--s3-key", util.S3ZipName)
	}
	util.RunCommandAndLog(cmd)
}

//UpdateFunctionConfiguration  will update the configuration (e.g. timeout) of the serverless function with id `i`.
func (amazon Instance) UpdateFunctionConfiguration(i int) (string, string) {
	log.Infof("Updating producer lambda configuration %s-%v", amazon.appName, i)
	assignedMemory := "128"
	cmd := exec.Command("/usr/local/bin/aws", "lambda", "update-function-configuration",
		"--function-name", fmt.Sprintf("%s-%v", amazon.appName, i), "--timeout", "900", "--memory-size", assignedMemory)
	util.RunCommandAndLog(cmd)
	return amazon.getAPIID(i), assignedMemory
}
