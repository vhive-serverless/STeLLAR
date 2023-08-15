package setup

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"stellar/setup"
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
		Package: setup.Package{Patterns: []string{"pattern1"}},
		Functions: map[string]*setup.Function{
			"test1_2_0": {
				Name:    "test1_2_0",
				Handler: "hellopy/lambda_function.lambda_handler",
				Runtime: "Python3.8",
				Events: []setup.Event{
					{HttpApi: setup.HttpApi{Path: "/test1_2_0", Method: "GET"}}}},
			"test1_2_1": {
				Name:    "test1_2_1",
				Handler: "hellopy/lambda_function.lambda_handler",
				Runtime: "Python3.8",
				Events: []setup.Event{
					{HttpApi: setup.HttpApi{Path: "/test1_2_1", Method: "GET"}}}},
		}}

	actual := &setup.Serverless{}

	subEx := &setup.SubExperiment{Title: "test1", Parallelism: 2, Runtime: "Python3.8", Handler: "hellopy/lambda_function.lambda_handler", PackagePattern: "pattern1"}
	actual.AddFunctionConfig(subEx, 2)

	require.Equal(t, expected, actual)
}

func TestCreateServerlessConfigFile(t *testing.T) {
	assert := require.New(t)

	// Define the expected Serverless struct
	serverless := &setup.Serverless{
		Service:          "TestService",
		FrameworkVersion: "3",
		Provider: setup.Provider{
			Name:    "aws",
			Runtime: "go1.x",
			Region:  "us-east-1",
		},
		Package: setup.Package{
			Patterns: []string{"pattern1", "pattern2"},
		},
		Functions: map[string]*setup.Function{
			"testFunction1": {
				Handler: "handler1",
				Runtime: "go1.16",
				Name:    "testFunction1",
				Events: []setup.Event{
					{
						HttpApi: setup.HttpApi{
							Path:   "/test1",
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
	actualData, err := ioutil.ReadFile("serverless.yml")
	assert.NoError(err, "Error reading actual data")

	// Generate YAML content from the expected Serverless struct
	expectedData, err := ioutil.ReadFile("test.yml")
	assert.NoError(err, "Error marshaling expected data")

	// Compare the contents byte by byte
	assert.True(bytes.Equal(expectedData, actualData), "YAML content mismatch")

}

func TestRemoveService(t *testing.T) {
	msg := setup.RemoveService()

	require.Equal(t, msg, "")
}

func TestDeployService(t *testing.T) {
	msg := setup.DeployService()
	require.Equal(t, msg, "")
}

func TestAddPackagePattern(t *testing.T) {
	assert := require.New(t)

	// Create a sample Serverless instance
	serverless := &setup.Serverless{
		Package: setup.Package{
			Patterns: []string{"pattern1", "pattern2"},
		},
	}

	// Call the AddPackagePattern function with a new pattern
	newPattern := "pattern3"
	serverless.AddPackagePattern(newPattern)

	// Verify that the new pattern has been added
	assert.Contains(serverless.Package.Patterns, newPattern, "New pattern not added")

	// Call the AddPackagePattern function with an existing pattern
	existingPattern := "pattern1"
	serverless.AddPackagePattern(existingPattern)

	// Verify that the existing pattern is not duplicated
	count := 0
	for _, p := range serverless.Package.Patterns {
		if p == existingPattern {
			count++
		}
	}
	assert.Equal(1, count, "Existing pattern duplicated")
}
