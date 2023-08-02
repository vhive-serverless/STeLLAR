package setup

import (
	"github.com/stretchr/testify/require"
	"stellar/setup"
	"testing"
)

func TestAddFunctionConfig(t *testing.T) {
	result_s := &setup.Serverless{}
	require_s := &setup.Serverless{}
	result_s.AddFunctionConfig()

	require.Equal(t, result_s, require_s)
}

func TestCreateServerlessConfigFile(t *testing.T) {
	result_s := &setup.Serverless{}
	result_s.CreateServerlessConfigFile()
}

func TestRemoveService(t *testing.T) {
	msg := setup.RemoveService()

	require.Equal(t, msg, "")
}

func TestDeployService(t *testing.T) {
	msg := setup.DeployService()
	require.Equal(t, msg, "")
}
