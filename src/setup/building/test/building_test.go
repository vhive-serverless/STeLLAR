package building

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"stellar/setup/building"
	"testing"
)

func TestBuildFunctionJava(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("aws", "test/function/path", "java11")
}

func TestBuildFunctionGolang(t *testing.T) {
	b := &building.Builder{}
	err := os.Chdir("../../..") // so that BuildFunction generates binaries in the correct path relative to the /src directory
	if err != nil {
		log.Fatal("Failed to change to /src directory ")
	}
	b.BuildFunction("aws", "hellogo", "go1.x")
	assert.FileExists(t, "setup/deployment/raw-code/serverless/aws/hellogo/main")
}

func TestBuildFunctionUnsupported(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("mockProvider", "mockFunctionName", "unsupported")
}
