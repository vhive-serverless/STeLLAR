package generator

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"lambda-benchmarking/client/setup/functions/util"
	"math/rand"
	"os"
	"os/exec"
)

//SetupDeployment will create the serverless function zip deployment for the given provider,
//in the given language and of the given size in bytes.
func SetupDeployment(provider string, language string, sizeBytes int64) {
	zippedBinarySizeBytes := int64(createBinary(provider, language))

	if sizeBytes < zippedBinarySizeBytes {
		log.Fatalf("Total size (~%vMB) cannot be smaller than zipped binary size (~%vMB).",
			util.BytesToMB(sizeBytes),
			util.BytesToMB(zippedBinarySizeBytes))
	}

	randomFileName := "random.file"
	generateRandomFile(sizeBytes-zippedBinarySizeBytes, randomFileName)
	generateZIP(provider, randomFileName, util.BytesToMB(sizeBytes))
}

func generateRandomFile(sizeBytes int64, randomFileName string) {
	if fileExists(randomFileName) {
		log.Infof("Random file `%s` already exists, removing...", randomFileName)
		if err := os.Remove(randomFileName); err != nil {
			log.Fatalf("Failed to remove random file `%s`", randomFileName)
		}
	}

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("Failed to fill buffer with random bytes: `%s`", err.Error())
	}

	if err := ioutil.WriteFile(randomFileName, buffer, 0666); err != nil {
		log.Fatalf("Could not generate random file with size %d bytes", sizeBytes)
	}
}

func runCommandAndLog(cmd *exec.Cmd) string {
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

