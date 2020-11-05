package util

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

//RunCommandAndLog will execute a bash command and log results to the standard logger, as well as return the contents.
func RunCommandAndLog(cmd *exec.Cmd) string {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s: %s", fmt.Sprint(err), stderr.String())
	}
	log.Debugf("Command result: %s", out.String())
	return out.String()
}