package building

import (
	"log"
	"os"
	"stellar/setup/building"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BuildingTestSuite struct {
	suite.Suite
}

func (s *BuildingTestSuite) SetupSuite() {
	if err := os.Chdir("../../.."); err != nil { // so that BuildFunction generates binaries in the correct path relative to the /src directory
		log.Fatal("Failed to change to /src directory ")
	}
}

func (s *BuildingTestSuite) TestBuildFunctionJava() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellojava", "java11")
	assert.FileExists(s.T(), "setup/deployment/raw-code/serverless/aws/artifacts/hellojava/hellojava.zip")
}

func (s *BuildingTestSuite) TestBuildFunctionGolang() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellogo", "go1.x")
	assert.FileExists(s.T(), "setup/deployment/raw-code/serverless/aws/artifacts/hellogo/bootstrap")
}

func (s *BuildingTestSuite) TestBuildFunctionUnsupported() {
	b := &building.Builder{}
	b.BuildFunction("mockProvider", "mockFunctionName", "unsupported")
}

func TestBuildingTestSuite(t *testing.T) {
	suite.Run(t, new(BuildingTestSuite))
}
