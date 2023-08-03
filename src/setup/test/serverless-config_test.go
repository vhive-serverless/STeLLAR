package setup

import (
	"github.com/stretchr/testify/require"
	"stellar/setup"
	"testing"
)

func TestAddFunctionConfig(t *testing.T) {
	resultServerless := &setup.Serverless{}
	requireServerless := &setup.Serverless{}
	subEx := &setup.SubExperiment{Title: "test1"}
	resultServerless.AddFunctionConfig(subEx, 0)

	require.Equal(t, resultServerless, requireServerless)
}

func TestCreateServerlessConfigFile(t *testing.T) {
	resultServerless := &setup.Serverless{}
	resultServerless.CreateServerlessConfigFile()
}

func TestRemoveService(t *testing.T) {
	msg := setup.RemoveService()

	require.Equal(t, msg, "")
}

func TestDeployService(t *testing.T) {
	msg := setup.DeployService()
	require.Equal(t, msg, "")
}
