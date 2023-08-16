package building

import (
	"github.com/stretchr/testify/assert"
	"os"
	"stellar/setup/building"
	"testing"
)

func TestBuildFunctionJava(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("mockProvider", "mockFunctionName", "java")
}

func TestBuildFunctionGolang(t *testing.T) {
	b := &building.Builder{}
	err := os.Chdir("../../..")
	b.BuildFunction("aws", "hellogo", "go1.x")
	_, err = os.Stat("setup/deployment/raw-code/serverless/aws/hellogo")
	assert.NoError(t, err)
}

func TestBuildFunctionUnsupported(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("mockProvider", "mockFunctionName", "unsupported")
}
