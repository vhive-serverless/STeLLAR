package setup

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
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
	actual.CreateHeaderConfig(config, "STeLLAR")

	require.Equal(t, expected, actual)
}

func TestAddFunctionConfigAWS(t *testing.T) {
	expected := &setup.Serverless{
		Package: setup.Package{Individually: true},
		Functions: map[string]*setup.Function{
			"test1-2-0": {
				Name:    "test1-2-0",
				Handler: "main.lambda_handler",
				Runtime: "Python3.8",
				Package: setup.FunctionPackage{
					Patterns: []string{"main.py"},
					Artifact: "",
				},
				Events: []setup.Event{
					{AWSHttpEvent: setup.AWSHttpEvent{Path: "/test1-2-0", Method: "GET"}}}},
			"test1-2-1": {
				Name:    "test1-2-1",
				Handler: "main.lambda_handler",
				Runtime: "Python3.8",
				Package: setup.FunctionPackage{
					Patterns: []string{"main.py"},
					Artifact: "",
				},
				Events: []setup.Event{
					{AWSHttpEvent: setup.AWSHttpEvent{Path: "/test1-2-1", Method: "GET"}}}},
		}}

	actual := &setup.Serverless{Package: setup.Package{Individually: true}}

	subEx := &setup.SubExperiment{Title: "test1", Parallelism: 2, Runtime: "Python3.8", Handler: "main.lambda_handler", PackagePattern: "main.py"}
	actual.AddFunctionConfigAWS(subEx, 2, "")

	require.Equal(t, expected, actual)

	require.Equal(t, []string{"test1-2-0", "test1-2-1"}, subEx.Routes)
}

func TestAddFunctionConfigAzure(t *testing.T) {
	expected := &setup.Serverless{
		Functions: map[string]*setup.Function{
			"test1-2-0": {
				Name:    "test1-2-0",
				Handler: "main.main",
				Runtime: "Python3.8",
				Package: setup.FunctionPackage{
					Patterns: []string{"main.py"},
					Artifact: "",
				},
				Events: []setup.Event{
					{
						AzureHttp:      true,
						AzureMethods:   []string{"GET"},
						AzureAuthLevel: "anonymous",
					},
				},
			},
			"test1-2-1": {
				Name:    "test1-2-1",
				Handler: "main.main",
				Runtime: "Python3.8",
				Package: setup.FunctionPackage{
					Patterns: []string{"main.py"},
					Artifact: "",
				},
				Events: []setup.Event{
					{
						AzureHttp:      true,
						AzureMethods:   []string{"GET"},
						AzureAuthLevel: "anonymous",
					},
				},
			},
		},
	}

	actual := &setup.Serverless{}

	subEx := &setup.SubExperiment{Title: "test1", Parallelism: 2, Runtime: "Python3.8", Handler: "main.main", PackagePattern: "main.py"}
	actual.AddFunctionConfigAzure(subEx, 2, "")

	require.Equal(t, expected, actual)

	require.Equal(t, []string{"test1-2-0", "test1-2-1"}, subEx.Routes)
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
				Name:    "parallelism1-0-0",
				Package: setup.FunctionPackage{
					Patterns: []string{"hellopy/lambda_function.py"},
				},
				Events: []setup.Event{
					{
						AWSHttpEvent: setup.AWSHttpEvent{
							Path:   "/parallelism1-0-0",
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
	expectedData, err := os.ReadFile("test.yml")
	assert.NoError(err, "Error marshaling expected data")

	// Compare the contents byte by byte
	assert.True(bytes.Equal(expectedData, actualData), "YAML content mismatch")

}

func TestCreateServerlessConfigFileSnapStart(t *testing.T) {
	assert := require.New(t)

	// Define the expected Serverless struct
	serverless := &setup.Serverless{
		Service:          "TestService",
		FrameworkVersion: "3",
		Provider: setup.Provider{
			Name:    "aws",
			Runtime: "java11",
			Region:  "us-west-1",
		},
		Package: setup.Package{
			Individually: true,
		},
		Functions: map[string]*setup.Function{
			"testFunction1": {
				Handler:   "org.hellojava.Handler",
				Runtime:   "java11",
				Name:      "parallelism1-0-0",
				SnapStart: true,
				Package: setup.FunctionPackage{
					Patterns: []string{},
				},
				Events: []setup.Event{
					{
						AWSHttpEvent: setup.AWSHttpEvent{
							Path:   "/parallelism1-0-0",
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
	expectedData, err := os.ReadFile("test_snapstart.yml")
	assert.NoError(err, "Error marshaling expected data")

	// Compare the contents byte by byte
	assert.True(bytes.Equal(expectedData, actualData), "YAML content mismatch")

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

func TestGetAWSEndpointIdSingleFunction(t *testing.T) {
	testMsg := "\nendpoint: GET - https://7rhr5111eg.execute-api.us-west-1.amazonaws.com/parallelism1_0_0\nfunctions:\n  testFunction1: parallelism1_0_0 (3.5 kB)\n"
	actual := setup.GetAWSEndpointID(testMsg)
	require.Equal(t, "7rhr5111eg", actual)
}

func TestGetAWSEndpointIdMultipleFunctions(t *testing.T) {
	testMsg := "\nendpoints:\n  GET - https://z4a0lmtx64.execute-api.us-west-1.amazonaws.com/parallelism1_0_0\n  GET - https://z4a0lmtx64.execute-api.us-west-1.amazonaws.com/parallelism2_1_0\n  GET - https://z4a0lmtx64.execute-api.us-west-1.amazonaws.com/parallelism2_1_1\nfunctions:\n  parallelism1_0_0: parallelism1_0_0 (3.5 kB)\n  parallelism2_1_0: parallelism2_1_0 (3.5 kB)\n  parallelism2_1_1: parallelism2_1_1 (3.5 kB)\n"
	actual := setup.GetAWSEndpointID(testMsg)
	require.Equal(t, "z4a0lmtx64", actual)
}

func TestGetAzureEndpointID(t *testing.T) {
	testMsg := "Deployed serverless functions:\n-> subexperiment2_1_0: [GET] sls-seasi-dev-stellar-sub-experiment-1.azurewebsites.net/api/subexperiment2_1_0\n-> subexperiment2_1_1: [GET] sls-seasi-dev-stellar-sub-experiment-1.azurewebsites.net/api/subexperiment2_1_1\n"
	actual := setup.GetAzureEndpointID(testMsg)
	require.Equal(t, "sls-seasi-dev-stellar-sub-experiment-1", actual)
}

func TestGetGCREndpointID(t *testing.T) {
	testMsg := "Service [test-function] revision [test-function-00001-cec] has been deployed and is serving 100 percent of traffic.\nService URL: https://test-function-nfjrndgaha-uw.a.run.app"
	actual := setup.GetGCREndpointID(testMsg)
	require.Equal(t, "test-function-nfjrndgaha-uw.a.run.app", actual)
}

func TestGetCloudflareEndpointID(t *testing.T) {
	testMsg := "Published hellojs_wrangler0 (3.48 sec)\nhttps://hellonode.stellarbench.workers.dev\nCurrent Deployment ID: 26923084-4e66-4b4b-b876-cb85341b75f6"
	actual := setup.GetCloudflareEndpointID(testMsg)
	require.Equal(t, "hellonode.stellarbench.workers.dev", actual)
}
