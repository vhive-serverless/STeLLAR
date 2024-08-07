package packaging

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"stellar/setup/building"
	"stellar/setup/deployment/packaging"
	"stellar/util"
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

func (s *ZipTestSuite) TestGenerateFillerFile() {
	expectedFillerFileSizeMB := 30.0

	packaging.GenerateFillerFile(1, "filler.file", util.MebibyteToBytes(expectedFillerFileSizeMB))

	fileInfo, err := os.Stat("filler.file")
	if err != nil {
		assert.Fail(s.T(), "Could not obtain file info of ZIP artifact")
	}
	actualFillerFileSize := util.BytesToMebibyte(fileInfo.Size())
	assert.InDelta(s.T(), expectedFillerFileSizeMB, actualFillerFileSize, 0.1)

	err = os.Remove("filler.file")
	if err != nil {
		log.Warn("Failed to clean up temporary filler files created by unit tests")
	}
}

func (s *ZipTestSuite) TestGenerateServerlessZipArtifactsPython() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellopy", "python3.9")
	packaging.GenerateServerlessZIPArtifacts(1, "aws", "python3.9", "hellopy", 50)
	fileInfo, err := os.Stat("setup/deployment/raw-code/serverless/aws/artifacts/hellopy/hellopy.zip")
	if err != nil {
		assert.Fail(s.T(), "Could not obtain file info of ZIP artifact")
	}
	assert.InDelta(s.T(), 50, util.BytesToMebibyte(fileInfo.Size()), 0.1)
}

func (s *ZipTestSuite) TestGenerateServerlessZipArtifactsGolang() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellogo", "go1.x")
	packaging.GenerateServerlessZIPArtifacts(2, "aws", "go1.x", "hellogo", 50)
	fileInfo, err := os.Stat("setup/deployment/raw-code/serverless/aws/artifacts/hellogo/hellogo.zip")
	if err != nil {
		assert.Fail(s.T(), "Could not obtain file info of ZIP artifact")
	}
	assert.InDelta(s.T(), 50, util.BytesToMebibyte(fileInfo.Size()), 0.1)
}

func (s *ZipTestSuite) TestGenerateServerlessZipArtifactsJava() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellojava", "java11")
	packaging.GenerateServerlessZIPArtifacts(3, "aws", "java11", "hellojava", 50)
	fileInfo, err := os.Stat("setup/deployment/raw-code/serverless/aws/artifacts/hellojava/hellojava.zip")
	if err != nil {
		assert.Fail(s.T(), "Could not obtain file info of ZIP artifact")
	}
	assert.InDelta(s.T(), 50, util.BytesToMebibyte(fileInfo.Size()), 0.1)
}

func (s *ZipTestSuite) TestGenerateServerlessZipArtifactsNode() {
	b := &building.Builder{}
	b.BuildFunction("aws", "hellonode", "nodejs18.x")
	packaging.GenerateServerlessZIPArtifacts(2, "aws", "nodejs18.x", "hellonode", 50)
	fileInfo, err := os.Stat("setup/deployment/raw-code/serverless/aws/artifacts/hellonode/hellonode.zip")
	if err != nil {
		assert.Fail(s.T(), "Could not obtain file info of ZIP artifact")
	}
	assert.InDelta(s.T(), 50, util.BytesToMebibyte(fileInfo.Size()), 0.1)
}

func (s *ZipTestSuite) TestGenerateServerlessZipArtifactsRuby() {
	b := &building.Builder{}
	b.BuildFunction("aws", "helloruby", "ruby3.2")
	packaging.GenerateServerlessZIPArtifacts(2, "aws", "ruby3.2", "helloruby", 50)
	fileInfo, err := os.Stat("setup/deployment/raw-code/serverless/aws/artifacts/helloruby/helloruby.zip")
	if err != nil {
		assert.Fail(s.T(), "Could not obtain file info of ZIP artifact")
	}
	assert.InDelta(s.T(), 50, util.BytesToMebibyte(fileInfo.Size()), 0.1)
}

func TestZipTestSuite(t *testing.T) {
	suite.Run(t, new(ZipTestSuite))
}
