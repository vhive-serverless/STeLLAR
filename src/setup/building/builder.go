package building

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"stellar/util"
)

// Builder struct keeps track of the functions built by stellar
type Builder struct {
	functionsBuilt []string
}

func (b *Builder) BuildFunction(provider string, functionName string, runtime string) {
	// TODO: Implement function

	// First we check whether the function has not been built already
	// TODO: Check if function path is in functionsBuilt if yes, skip the build. If no, continue the building process and add the functionPath to the list.

	functionPath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/%s", provider, functionName)
	b.functionsBuilt = append(b.functionsBuilt, functionPath)
	switch runtime {
	case "java":
		buildJava(functionPath)
	case "go1.x":
		buildGolang(functionPath)
	default:
		// building not supported
		log.Warnf("Building runtime %s is not necessary, or not supported. Continuing without building.", runtime)
	}
}

// buildJava builds the java zip artifact for serverless deployment using Gradle
func buildJava(functionPath string) {
	// TODO: Implement function.
	log.Warn(functionPath)
}

// buildGolang builds the Golang binary for serverless deployment
func buildGolang(functionPath string) {
	log.Infof("Building Go binary at %s path", functionPath)
	command := exec.Command("env", "GOOS=linux", "GOARCH=amd64", "go", "build", "-C", functionPath)
	util.RunCommandAndLog(command)
}
