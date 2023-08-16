package building

import (
	"bytes"
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

	actual, err := os.ReadFile("resources/hellogo/main")
	assert.NoError(t, err, "Failed to read actual Go binary built")

	expected, err := os.ReadFile("resources/hellogo/expectedBinary")
	assert.NoError(t, err, "Failed to read expected Go binary built")

	assert.True(t, bytes.Equal(expected, actual), "Binary content mismatch")
}

func TestBuildFunctionUnsupported(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("test/function/path", "unsupported")
}
