package generator

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func createBinary(provider string, language string) int {
	log.Info("Building binary file for function to be deployed...")
	codePath := fmt.Sprintf("code/producer/%s/%s-handler.go", language, provider)
	if !fileExists(codePath) {
		log.Fatalf("Code path `%s` does not exist, cannot deploy/update code.", codePath)
	}

	switch language {
	case "go1.x":
		util.RunCommandAndLog(exec.Command("go", "build", "-v", "-o", util.BinaryName,
			"code/producer/go1.x/aws-handler.go"))
	//TODO: add python3 support
	//case "python3.8":
	//	RunCommandAndLog(exec.Command("python", "build", "-v", "-race", "-o", "producer-handler"))
	default:
		log.Fatalf("Unrecognized language %s", language)
	}

	log.Info("Zipping binary file to find its size...")
	util.RunCommandAndLog(exec.Command("zip", "zipped-binary", util.BinaryName))

	fi, err := os.Stat("zipped-binary.zip")
	if err != nil {
		log.Fatalf("Could not get size of zipped binary file: %s", err.Error())
	}
	util.RunCommandAndLog(exec.Command("rm", "-r", "zipped-binary.zip"))

	return int(fi.Size())
}