package building

import (
	"github.com/stretchr/testify/assert"
	"os"
	"stellar/setup/building"
	"testing"
)

func TestBuildFunctionJava(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("test/function/path", "java")
}

func TestBuildFunctionGolang(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("resources/hellogo", "go1.x")
	_, err := os.Stat("./resources/hellogo/main")
	assert.NoError(t, err)
}

func TestBuildFunctionUnsupported(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("test/function/path", "unsupported")
}
