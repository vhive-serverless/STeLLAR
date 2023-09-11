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
	require.Equal(t, 1, linesRemove)
}
