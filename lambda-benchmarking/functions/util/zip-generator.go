package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
)

const (
	randomFileName = "random.file"
	zipName        = "benchmarking"
)

func GenerateDeploymentZIP(provider string, language string, sizeBytes int) {
	zipPath := fmt.Sprintf("%s.zip", zipName)
	if fileExists(zipPath) {
		log.Infof("ZIP archive `%s` already exists, removing...", zipPath)
		if err := os.Remove(zipPath); err != nil {
			log.Fatalf("Failed to remove ZIP archive `%s`", zipPath)
		}
	}

	log.Infof("Building %s handler...", language)
	codePath := fmt.Sprintf("code/producer/%s/%s-handler.go", language, provider)
	if !fileExists(codePath) {
		log.Fatalf("Code path `%s` does not exist, cannot deploy/update code.", codePath)
	}

	switch language {
	case "go1.x":
		RunCommandAndLog(exec.Command("go", "build", "-v", "-o", "producer-handler",
			"code/producer/go1.x/aws-handler.go"))
	//case "python3.8":
	//	RunCommandAndLog(exec.Command("go", "build", "-v", "-race", "-o", "producer-handler"))
	default:
		log.Fatalf("Unrecognized language %s", language)
	}

	GenerateRandomFile(sizeBytes)
	RunCommandAndLog(exec.Command("zip", fmt.Sprintf("%s.zip", zipName), "producer-handler", randomFileName))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GenerateRandomFile(sizeBytes int) {
	if sizeBytes > 50*1000*1000 {
		log.Fatalf(`Deployment package is larger than 50 MB (~%dMB), you must use Amazon S3 (https://docs.aws.amazon.com/lambda/latest/dg/python-package.html).`,
			sizeBytes/1000000.0)
	}

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
