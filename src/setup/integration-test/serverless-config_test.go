package setup

import (
	"github.com/stretchr/testify/require"
	"os/exec"
	"stellar/setup"
	"stellar/util"
	"strings"
	"testing"
)

// If this test is failing on your local machine, try running it with sudo.
func TestDeployAndRemoveServiceAWS(t *testing.T) {
	// The two unit tests were merged together in order to make sure we are not left with a number of deployed test function on the cloud which are never used in.
	assert := require.New(t)

	util.RunCommandAndLog(exec.Command("cp", "aws-integration-test-serverless.yml", "../deployment/raw-code/serverless/aws/serverless.yml"))

	msgDeploy := setup.DeployService("../deployment/raw-code/serverless/aws/")
	msgRemove := setup.RemoveAWSService("../deployment/raw-code/serverless/aws/")

	assert.True(strings.Contains(msgDeploy, "Service deployed"))
	assert.True(strings.Contains(msgRemove, "successfully removed"))
}

func TestDeployAndRemoveContainerService(t *testing.T) {
	s := &setup.Serverless{
		Service:          "STeLLAR",
		FrameworkVersion: "3",
		Provider: setup.Provider{
			Name:    "gcr",
			Runtime: "python3.9",
			Region:  "us-west1",
		},
	}

	subex := &setup.SubExperiment{
		Title:       "test_hellopy",
		Parallelism: 1,
	}

	s.DeployGCRContainerService(subex, 0, "docker.io/kkmin/hellopy", "../deployment/raw-code/serverless/gcr/hellopy/", "us-west1")
	deleteMsg := setup.RemoveService("gcr", "../deployment/raw-code/serverless/gcr/hellopy/", 1)
	require.Equal(t, "All GCR services deleted.", deleteMsg)
}

func TestDeployAndRemoveServiceAzure(t *testing.T) {
	assert := require.New(t)

	util.RunCommandAndLog(exec.Command("cp", "azure-integration-test-serverless.yml", "../deployment/raw-code/serverless/azure/hellopy/serverless.yml"))

	msgDeploy := setup.DeployService("../deployment/raw-code/serverless/azure/hellopy/")
	msgRemove := setup.RemoveAzureSingleService("../deployment/raw-code/serverless/azure/hellopy/")

	assert.True(strings.Contains(msgDeploy, "Deployed serverless functions"))
	assert.True(strings.Contains(msgRemove, "successfully removed"))
}
