package setup

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os/exec"
	"stellar/setup"
	"stellar/util"
	"strings"
	"testing"
)

// If this test is failing on your local machine, try running it with sudo.
func TestDeployAndRemoveService(t *testing.T) {
	// The two unit tests were merged together in order to make sure we are not left with a number of deployed test function on the cloud which are never used in.
	util.RunCommandAndLog(exec.Command("cp", "test.yml", "../deployment/raw-code/serverless/aws/serverless.yml"))

	msgDeploy := setup.DeployService("../deployment/raw-code/serverless/aws/")

	linesDeploy := len(strings.Split(msgDeploy, "\n"))

	msgRemove := setup.RemoveService("aws", "../deployment/raw-code/serverless/aws/")
	linesRemove := len(strings.Split(msgRemove, "\n"))
	log.Info(msgDeploy)
	log.Info(msgRemove)
	require.Equal(t, 5, linesDeploy)
	require.Equal(t, 4, linesRemove)
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

	s.DeployContainerService(subex, 0, "docker.io/kkmin/hellopy", "../deployment/raw-code/serverless/gcr/hellopy/", "us-west1")
	deleteMsg := setup.RemoveService("gcr", "../deployment/raw-code/serverless/gcr/hellopy/")
	require.Equal(t, "All GCR services deleted.", deleteMsg)

}
