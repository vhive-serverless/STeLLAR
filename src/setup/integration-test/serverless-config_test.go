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
	msgRemove := setup.RemoveServerlessService("../deployment/raw-code/serverless/aws/")

	assert.True(strings.Contains(msgDeploy, "Service deployed"))
	assert.True(strings.Contains(msgRemove, "successfully removed"))
}

func TestDeployAndRemoveServiceGCR(t *testing.T) {
	assert := require.New(t)
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
		Title:       "hellopytest",
		Parallelism: 1,
	}

	s.DeployGCRContainerService(subex, 0, "abc12", "docker.io/kkmin/hellopy", "../deployment/raw-code/serverless/gcr/hellopy/", "us-west1")
	deleteMsg := setup.RemoveGCRSingleService("abc12-hellopytest-0-0")
	assert.True(strings.Contains(deleteMsg, "Deleted service [abc12-hellopytest-0-0]"))
}

func TestDeployAndRemoveServiceGCRCPUBoost(t *testing.T) {
	assert := require.New(t)
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
		Title:           "cpuboosttest",
		Parallelism:     1,
		CPUBoostEnabled: true,
	}

	s.DeployGCRContainerService(subex, 0, "def12", "docker.io/kkmin/hellopy", "../deployment/raw-code/serverless/gcr/hellopy/", "us-west1")
	deleteMsg := setup.RemoveGCRSingleService("def12-cpuboosttest-0-0")
	assert.True(strings.Contains(deleteMsg, "Deleted service [def12-cpuboosttest-0-0]"))
}

func TestDeployAndRemoveServiceAzure(t *testing.T) {
	assert := require.New(t)

	util.RunCommandAndLog(exec.Command("cp", "azure-integration-test-serverless.yml", "../deployment/raw-code/serverless/azure/hellopy/serverless.yml"))

	msgDeploy := setup.DeployService("../deployment/raw-code/serverless/azure/hellopy/")
	msgRemove := setup.RemoveServerlessServiceForcefully("../deployment/raw-code/serverless/azure/hellopy/")

	assert.True(strings.Contains(msgDeploy, "Deployed serverless functions"))
	assert.True(strings.Contains(msgRemove, "successfully removed"))
}

func TestDeployAndRemoveServiceCloudflare(t *testing.T) {
	assert := require.New(t)

	subex := &setup.SubExperiment{
		Title:       "cloudflaretest",
		Function:    "hellonode",
		Handler:     "index.js",
		Parallelism: 1,
	}

	setup.DeployCloudflareWorkers(subex, 0, "abc12", "../deployment/raw-code/serverless/cloudflare")
	msgRemove := setup.RemoveCloudflareSingleWorker("abc12-cloudflaretest-0-0")

	assert.True(strings.Contains(msgRemove, "Successfully deleted"))
}

/*
func TestDeployAndRemoveServiceAlibaba(t *testing.T) {
	assert := require.New(t)

	util.RunCommandAndLog(exec.Command("cp", "aliyun-integration-test-serverless.yml", "../deployment/raw-code/serverless/aliyun/hellopy/serverless.yml"))

	msgDeploy := setup.DeployService("../deployment/raw-code/serverless/aliyun/hellopy/")
	msgRemove := setup.RemoveServerlessService("../deployment/raw-code/serverless/aliyun/hellopy/")

	assert.True(strings.Contains(msgDeploy, "Deployed API"))
	assert.True(strings.Contains(msgRemove, "Removed service"))
}
*/
