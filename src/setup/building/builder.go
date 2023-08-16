package building

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"stellar/util"
)

const artifactDir = "setup/artifacts"

// Builder struct keeps track of the functions built by stellar
type Builder struct {
	functionsBuilt map[string]bool
}

func (b *Builder) BuildFunction(provider string, functionName string, runtime string) string {
	// First we check whether the function has not been built already
	if b.functionsBuilt[functionName] {
		log.Warnf("Function %s already built. Skipping.", functionName)
		return fmt.Sprintf("%s/%s/%s.zip", artifactDir, functionName, functionName)
	}

	functionPath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/%s", provider, functionName)
	switch runtime {
	case "java":
		buildJava(functionPath)
	case "go1.x":
		buildGolang(functionPath)
	default:
		// building not supported
		log.Warnf("Building runtime %s is not necessary, or not supported. Continuing without building.", runtime)
		return ""
	}
	b.functionsBuilt[functionName] = true
	return fmt.Sprintf("%s/%s/%s.zip", artifactDir, functionName, functionName)
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
