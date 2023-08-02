package building

import (
	"stellar/setup/building"
	"testing"
)

func TestBuildFunctionJava(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("test/function/path", "java")
}

func TestBuildFunctionGolang(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("test/function/path", "golang")
}

func TestBuildFunctionUnsupported(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction("test/function/path", "unsupported")
}
