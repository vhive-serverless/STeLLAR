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
func TestDeployAndRemoveService(t *testing.T) {
	// The two unit tests were merged together in order to make sure we are not left with a number of deployed test function on the cloud which are never used in.
	assert := require.New(t)

	util.RunCommandAndLog(exec.Command("cp", "test.yml", "../deployment/raw-code/serverless/aws/serverless.yml"))

	msgDeploy := setup.DeployService("../deployment/raw-code/serverless/aws/")
	msgRemove := setup.RemoveServiceAWS("../deployment/raw-code/serverless/aws/")

	assert.True(strings.Contains(msgDeploy, "Service deployed"))
	assert.True(strings.Contains(msgRemove, "successfully removed"))
}
