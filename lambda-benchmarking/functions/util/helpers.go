package util

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	name = "benchmarking"
)

func CheckAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Fatalf("Environment variable %s is not set.", key)
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
	log.Printf("Result: %s", out.String())
	return out.String()
}
