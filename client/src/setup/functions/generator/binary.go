package generator

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/setup/functions/util"
	"os"
	"os/exec"
)

func createBinary(provider string, language string) int {
	log.Info("Building binary file for function to be deployed...")

	codePath := fmt.Sprintf("setup/functions/code/producer/%s/%s-handler.go", language, provider)
	if !fileExists(codePath) {
		log.Fatalf("Code path `%s` does not exist, cannot deploy/update code.", codePath)
	}

	switch language {
	case "go1.x":
		runCommandAndLog(exec.Command("go", "build", "-v", "-o", util.BinaryName, codePath))
	//TODO: add python3 support
	//case "python3.8":
	//	runCommandAndLog(exec.Command("python", "build", "-v", "-race", "-o", "producer-handler"))
	default:
		log.Fatalf("Unrecognized language %s", language)
	}

	log.Info("Zipping binary file to find its size...")
	runCommandAndLog(exec.Command("zip", "zipped-binary", util.BinaryName))

	fi, err := os.Stat("zipped-binary.zip")
	if err != nil {
		log.Fatalf("Could not get size of zipped binary file: %s", err.Error())
	}
	runCommandAndLog(exec.Command("rm", "-r", "zipped-binary.zip"))

	return int(fi.Size())
}
