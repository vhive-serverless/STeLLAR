package aws

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

// Logging
func runCommandAndLog(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s: %s", fmt.Sprint(err), stderr.String())
	}
	log.Printf("Result: %s", out.String())
}

func runCommandAndReturnOutput(cmd *exec.Cmd) string {
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(stdout)
}
