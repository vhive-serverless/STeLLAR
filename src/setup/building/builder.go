package building

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"stellar/util"
)

// Builder struct keeps track of the functions built by stellar
type Builder struct {
	functionsBuilt map[string]bool
}

func (b *Builder) BuildFunction(provider string, functionName string, runtime string) string {
	artifactDir := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/artifacts/%s", provider, functionName)

	// First we check whether the function has not been built already
	if b.functionsBuilt == nil {
		b.functionsBuilt = make(map[string]bool)
	}

	if b.functionsBuilt[functionName] {
		log.Warnf("Function %s already built. Skipping.", functionName)
		return fmt.Sprintf("artifacts/%s/%s.zip", functionName, functionName)
	}

	// Create folder in artifacts for the function
	if err := os.MkdirAll(artifactDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory for function %s: %s", functionName, err.Error())
	}

	functionDir := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/%s", provider, functionName)

	switch runtime {
	case "java11":
		buildJava(functionName, functionDir, artifactDir)
	case "go1.x":
		buildGolang(functionDir)
	default:
		// building not supported
		log.Warnf("Building runtime %s is not necessary, or not supported. Continuing without building.", runtime)
		return ""
	}
	b.functionsBuilt[functionName] = true
	return fmt.Sprintf("artifacts/%s/%s.zip", functionName, functionName)
}

// buildJava builds the java zip artifact for serverless deployment using Gradle
func buildJava(functionName string, functionPath string, artifactDir string) string {
	log.Infof("Building Java from the source code at %s path", functionPath)
	util.RunCommandAndLog(exec.Command("gradle", "buildZip", "-p", functionPath))

	artifactPath := fmt.Sprintf("%s/%s.zip", artifactDir, functionName)
	util.RunCommandAndLog(exec.Command("mv", fmt.Sprintf("%s/build/distributions/%s.zip", functionPath, functionName), artifactPath))

	return artifactPath
}

// buildGolang builds the Golang binary for serverless deployment
func buildGolang(functionPath string) {
	log.Infof("Building Go binary from the source code at %s path", functionPath)
	command := exec.Command("env", "GOOS=linux", "GOARCH=amd64", "go", "build", "-C", functionPath)
	util.RunCommandAndLog(command)
}
