package packaging

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"stellar/setup/building"
	"stellar/setup/deployment/packaging"
	"testing"
)

type ZipTestSuite struct {
	suite.Suite
}

func (s *ZipTestSuite) SetupSuite() {
	if err := os.Chdir("../../../.."); err != nil {
		log.Fatal("Failed to change to /src directory ")
	}
}

func (s *ZipTestSuite) TestGenerateServerlessZipArtifactsPython() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellopy", "python3.9")
	packaging.GenerateServerlessZIPArtifacts(1, "aws", "python3.9", "hellopy", 50)
	assert.FileExists(s.T(), "setup/deployment/raw-code/serverless/aws/artifacts/hellopy/hellopy.zip")
}

func (s *ZipTestSuite) TestGenerateServerlessZipArtifactsGolang() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellogo", "go1.x")
	packaging.GenerateServerlessZIPArtifacts(2, "aws", "go1.x", "hellogo", 50)
	assert.FileExists(s.T(), "setup/deployment/raw-code/serverless/aws/artifacts/hellogo/hellogo.zip")
}

func TestZipTestSuite(t *testing.T) {
	suite.Run(t, new(ZipTestSuite))
}
