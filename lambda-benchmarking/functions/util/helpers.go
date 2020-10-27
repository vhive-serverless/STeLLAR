package util

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func CheckAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Errorf("Environment variable %s is not set.", key)
	}
	return envVar
}

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
