package building

import (
	"github.com/stretchr/testify/mock"
	"os/exec"
	"stellar/setup/building"
	"stellar/util"
	"testing"
)

type MockCommandRunner struct {
	mock.Mock
}

func (m *MockCommandRunner) RunCommandAndLog(cmd *exec.Cmd) string {
	args := m.Called(cmd)
	return args.String(0)
}

func TestBuildFunctionJava(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction(util.RunCommandAndLog, "test/function/path", "java")
}

func TestBuildFunctionGolang(t *testing.T) {
	mockCommandRunner := MockCommandRunner{}
	mockCommandRunner.On("RunCommandAndLog", mock.Anything).Return("")

	b := &building.Builder{}
	b.BuildFunction(mockCommandRunner.RunCommandAndLog, "resources/hellogo", "go1.x")

	expectedArgument := exec.Command("env", "GOOS=linux", "GOARCH=amd64", "go", "build", "-C", "resources/hellogo")
	mockCommandRunner.AssertCalled(t, "RunCommandAndLog", expectedArgument)
}

func TestBuildFunctionUnsupported(t *testing.T) {
	b := &building.Builder{}
	b.BuildFunction(util.RunCommandAndLog, "test/function/path", "unsupported")
}
