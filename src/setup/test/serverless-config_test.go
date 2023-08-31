package setup

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"stellar/setup"
	"stellar/util"
	"strings"
	"testing"
)

func TestCreateHeaderConfig(t *testing.T) {

	// Define the expected Serverless struct and its associated values
	expected := &setup.Serverless{
		Service:          "STeLLAR",
		FrameworkVersion: "3",
		Provider: setup.Provider{
			Name:    "aws",
			Runtime: "go1.x",
			Region:  "us-west-1",
		},
		Package: setup.Package{Individually: true},
	}

	// Define the Configuration struct for testing
	config := &setup.Configuration{
		Provider: "aws",
		Runtime:  "go1.x",
	}

	actual := &setup.Serverless{}
	actual.CreateHeaderConfig(config)

	require.Equal(t, expected, actual)
}

func TestAddFunctionConfig(t *testing.T) {
	expected := &setup.Serverless{
		Package: setup.Package{Individually: true},
		Functions: map[string]*setup.Function{
			"test1_2_0": {
				Name:    "test1_2_0",
				Handler: "hellopy/lambda_function.lambda_handler",
				Runtime: "Python3.8",
				Package: setup.FunctionPackage{
					Patterns: []string{"pattern1"},
					Artifact: "",
				},
				Events: []setup.Event{
					{HttpApiAWS: setup.HttpApi{Path: "/test1_2_0", Method: "GET"}}}},
			"test1_2_1": {
				Name:    "test1_2_1",
				Handler: "hellopy/lambda_function.lambda_handler",
				Runtime: "Python3.8",
				Package: setup.FunctionPackage{
					Patterns: []string{"pattern1"},
					Artifact: "",
				},
				Events: []setup.Event{
					{HttpApiAWS: setup.HttpApi{Path: "/test1_2_1", Method: "GET"}}}},
		}}

	actual := &setup.Serverless{Package: setup.Package{Individually: true}}

	subEx := &setup.SubExperiment{Title: "test1", Parallelism: 2, Runtime: "Python3.8", Handler: "hellopy/lambda_function.lambda_handler", PackagePattern: "pattern1"}
	actual.AddFunctionConfig(subEx, 2, "", "aws")

	require.Equal(t, expected, actual)

	require.Equal(t, []string{"test1_2_0", "test1_2_1"}, subEx.Routes)
}

func TestCreateServerlessConfigFile(t *testing.T) {
	assert := require.New(t)

	// Define the expected Serverless struct
	serverless := &setup.Serverless{
		Service:          "TestService",
		FrameworkVersion: "3",
		Provider: setup.Provider{
			Name:    "aws",
			Runtime: "python3.9",
			Region:  "us-west-1",
		},
		Package: setup.Package{
			Individually: true,
		},
		Functions: map[string]*setup.Function{
			"testFunction1": {
				Handler: "hellopy/lambda_function.lambda_handler",
				Runtime: "python3.9",
				Name:    "parallelism1_0_0",
				Package: setup.FunctionPackage{
					Patterns: []string{"hellopy/lambda_function.py"},
				},
				Events: []setup.Event{
					{
						HttpApiAWS: setup.HttpApi{
							Path:   "/parallelism1_0_0",
							Method: "GET",
						},
					},
				},
			},
		},
	}

	// Call the CreateServerlessConfigFile function
	serverless.CreateServerlessConfigFile("serverless.yml")

	// Read the contents of the generated YAML file
	actualData, err := os.ReadFile("serverless.yml")
	assert.NoError(err, "Error reading actual data")

	// Generate YAML content from the expected Serverless struct
	expectedData, err := os.ReadFile("test_aws.yml")
	assert.NoError(err, "Error marshaling expected data")

	// Compare the contents byte by byte
	assert.True(bytes.Equal(expectedData, actualData), "YAML content mismatch")

}

// If this test is failing on your local machine, try running it with sudo.
func TestDeployAndRemoveService(t *testing.T) {
	// The two unit tests were merged together in order to make sure we are not left with a number of deployed test function on the cloud which are never used in.
	util.RunCommandAndLog(exec.Command("cp", "test_aws.yml", "../deployment/raw-code/serverless/aws/serverless.yml"))

	msgDeploy := setup.DeployService("../deployment/raw-code/serverless/aws/")

	linesDeploy := len(strings.Split(msgDeploy, "\n"))

	msgRemove := setup.RemoveService("../deployment/raw-code/serverless/aws/")
	linesRemove := len(strings.Split(msgRemove, "\n"))
	log.Info(msgDeploy)
	log.Info(msgRemove)
	require.Equal(t, 5, linesDeploy)
	require.Equal(t, 1, linesRemove)
}

func TestAddPackagePattern(t *testing.T) {
	assert := require.New(t)

	// Create a sample Serverless function instance
	function := &setup.Function{
		Package: setup.FunctionPackage{
			Patterns: []string{"pattern1", "pattern2"},
		},
	}

	// Call the AddPackagePattern function with a new pattern
	newPattern := "pattern3"
	function.AddPackagePattern(newPattern)

	// Verify that the new pattern has been added
	assert.Contains(function.Package.Patterns, newPattern, "New pattern not added")

	// Call the AddPackagePattern function with an existing pattern
	existingPattern := "pattern1"
	function.AddPackagePattern(existingPattern)

	// Verify that the existing pattern is not duplicated
	count := 0
	for _, p := range function.Package.Patterns {
		if p == existingPattern {
			count++
		}
	}
	assert.Equal(1, count, "Existing pattern duplicated")
}

func TestGetEndpointIdSingleFunction(t *testing.T) {
	testMsg := "\nendpoint: GET - https://7rhr5111eg.execute-api.us-west-1.amazonaws.com/parallelism1_0_0\nfunctions:\n  testFunction1: parallelism1_0_0 (3.5 kB)\n"
	actual := setup.GetEndpointID(testMsg)
	require.Equal(t, "7rhr5111eg", actual)
}

func TestGetEndpointIdMultipleFunctions(t *testing.T) {
	testMsg := "\nendpoints:\n  GET - https://z4a0lmtx64.execute-api.us-west-1.amazonaws.com/parallelism1_0_0\n  GET - https://z4a0lmtx64.execute-api.us-west-1.amazonaws.com/parallelism2_1_0\n  GET - https://z4a0lmtx64.execute-api.us-west-1.amazonaws.com/parallelism2_1_1\nfunctions:\n  parallelism1_0_0: parallelism1_0_0 (3.5 kB)\n  parallelism2_1_0: parallelism2_1_0 (3.5 kB)\n  parallelism2_1_1: parallelism2_1_1 (3.5 kB)\n"
	actual := setup.GetEndpointID(testMsg)
	require.Equal(t, "z4a0lmtx64", actual)
}
