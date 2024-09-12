// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package setup provides support with loading the experiment configuration,
// preparing the sub-experiments and setting up the functions for benchmarking.
package setup

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"stellar/setup/building"
	code_generation "stellar/setup/code-generation"
	"stellar/setup/deployment/connection"
	"stellar/setup/deployment/connection/amazon"
	"stellar/setup/deployment/packaging"
	"stellar/util"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// ProvisionFunctions will deploy, reconfigure, etc. functions to get ready for the sub-experiments.
func ProvisionFunctions(config Configuration) {
	const (
		nicContentionWarnThreshold = 800 // Experimentally found
		storageSpaceWarnThreshold  = 500 // 500 * ~18KiB = 10MB just for 1 sub-experiment
	)

	// To filter out re-usable endpoints for continuous-benchmarking
	availableEndpoints := connection.Singleton.ListAPIs(config.SubExperiments[0].RepurposeIdentifier)

	for index, subExperiment := range config.SubExperiments {
		config.SubExperiments[index].ID = index

		for _, burstSize := range subExperiment.BurstSizes {
			if burstSize > nicContentionWarnThreshold {
				log.Warnf("Experiment %d has a burst of size %d, NIC (Network Interface Controller) contention may occur.",
					index, burstSize)
				if !promptForBool("Do you wish to continue?") {
					os.Exit(0)
				}
			}
		}

		if subExperiment.Bursts >= storageSpaceWarnThreshold &&
			(subExperiment.Visualization == "all" || subExperiment.Visualization == "histogram") {
			log.Warnf("SubExperiment %d is generating histograms for each burst, this will create a large number (%d) of new files (>10MB).",
				index, subExperiment.Bursts)
			if !promptForBool("Do you wish to continue?") {
				os.Exit(0)
			}
		}

		if availableEndpoints == nil { // hostname must be the endpoint itself (external URL)
			config.SubExperiments[index].Endpoints = []EndpointInfo{{ID: config.Provider}}
			continue
		}

		availableEndpoints = assignEndpoints(
			availableEndpoints,
			&config.SubExperiments[index],
			config.Provider,
		)
	}

	if amazon.AWSSingletonInstance != nil && amazon.AWSSingletonInstance.ImageURI != "" {
		log.Info("A deployment was made using container images, waiting 10 seconds for changes to take effect with the provider...")
		time.Sleep(time.Second * 10)
	}
}

// ProvisionFunctionsServerless will deploy, reconfigure, etc. functions to get ready for the sub-experiments.
func ProvisionFunctionsServerless(config *Configuration, serverlessDirPath string) {
	switch config.Provider {
	case "aws":
		ProvisionFunctionsServerlessAWS(config, serverlessDirPath)
	case "azure":
		ProvisionFunctionsServerlessAzure(config, serverlessDirPath)
	case "gcr":
		ProvisionFunctionsGCR(config, serverlessDirPath)
	case "cloudflare":
		ProvisionFunctionsCloudflare(config, serverlessDirPath)
	case "aliyun":
		ProvisionFunctionsServerlessAlibaba(config, serverlessDirPath)
	default:
		log.Fatalf("Provider %s not supported for deployment", config.Provider)
	}
}

// ProvisionFunctionsServerlessAWS will deploy, reconfigure, etc. functions to get ready for the sub-experiments.
func ProvisionFunctionsServerlessAWS(config *Configuration, serverlessDirPath string) {
	slsConfig := &Serverless{}
	builder := &building.Builder{}

	randomTag := util.GenerateRandLowercaseLetters(5)
	slsConfig.CreateHeaderConfig(config, fmt.Sprintf("STeLLAR-%s", randomTag))
	slsConfig.packageIndividually()

	for index, subExperiment := range config.SubExperiments {
		//TODO: generate the code
		code_generation.GenerateCode(subExperiment.Function, config.Provider)

		// TODO: build the functions (Java and Golang)
		artifactPathRelativeToServerlessConfigFile := builder.BuildFunction(config.Provider, subExperiment.Function, subExperiment.Runtime)
		slsConfig.AddFunctionConfigAWS(&config.SubExperiments[index], index, randomTag, artifactPathRelativeToServerlessConfigFile)

		// generate filler files and zip used as Serverless artifacts
		packaging.GenerateServerlessZIPArtifacts(subExperiment.ID, config.Provider, subExperiment.Runtime, subExperiment.Function, subExperiment.FunctionImageSizeMB)
	}

	slsConfig.CreateServerlessConfigFile(fmt.Sprintf("%sserverless.yml", serverlessDirPath))

	log.Infof("Starting functions deployment. Deploying %d functions to %s.", len(slsConfig.Functions), config.Provider)
	slsDeployMessage := DeployService(serverlessDirPath)
	log.Info(slsDeployMessage)

	// TODO: assign endpoints to subexperiments
	// Get the endpoints by scraping the serverless deploy message.

	endpointID := GetAWSEndpointID(slsDeployMessage)

	// Assign Endpoint ID to each deployed function
	for i := range config.SubExperiments {
		config.SubExperiments[i].AssignEndpointIDs(endpointID)
	}

}

// Code replace the logic that sets up Azure Functions deployment with Azure Databricks job submission
func ProvisionFunctionsServerlessAzure(config *Configuration, serverlessDirPath string) {
	randomExperimentTag := util.GenerateRandLowercaseLetters(5)

	for subExperimentIndex, subExperiment := range config.SubExperiments {
		// Generate the code for the sub-experiment
		code_generation.GenerateCode(subExperiment.Function, config.Provider)

		// Build the function (Java/Python)
		builder := &building.Builder{}
		builder.BuildFunction(config.Provider, subExperiment.Function, subExperiment.Runtime)

		// Initialize Databricks job submission
		submitDatabricksJob(config, subExperimentIndex, randomExperimentTag, subExperiment)
	}
}

// Add-on code: Handle submitting the Spark job to Azure Databricks
func submitDatabricksJob(config *Configuration, subExperimentIndex int, randomExperimentTag string, subExperiment SubExperiment) {
	// Define the Databricks Job Payload
	jobPayload := createDatabricksJobPayload(config, randomExperimentTag, subExperimentIndex, subExperiment)

	// Submit the Job to Databricks using the REST API
	jobID, err := sendDatabricksJobRequest(jobPayload)
	if err != nil {
		log.Fatalf("Failed to submit job to Azure Databricks: %v", err)
	}

	// Wait for the Job to complete
	jobStatus, err := waitForDatabricksJobCompletion(jobID)
	if err != nil || jobStatus != "SUCCESS" {
		log.Fatalf("Databricks job %s failed or did not complete successfully", jobID)
	}

	log.Infof("Databricks job completed successfully. Job ID: %s", jobID)

	// Assign the Job ID as an endpoint identifier (or other metadata)
	config.SubExperiments[subExperimentIndex].AssignEndpointIDs(jobID)
}

// Add-on code: Define the Databricks job submission payload
func createDatabricksJobPayload(config *Configuration, randomExperimentTag string, subExperimentIndex int, subExperiment SubExperiment) string {
	
	return fmt.Sprintf(`{
		"name": "job-%s-subex%d",
		"new_cluster": {
			"spark_version": "7.3.x-scala2.12",
			"node_type_id": "Standard_DS3_v2",
			"autoscale": {
				"min_workers": 1,
				"max_workers": 8
			}
		},
		"notebook_task": {
			"notebook_path": "/path/to/your/notebook",
			"base_parameters": {
				"param1": "value1"
			}
		}
	}`, randomExperimentTag, subExperimentIndex)
}

// Add-on code: Sends the job request to Azure Databricks using their REST API
func sendDatabricksJobRequest(jobPayload string) (string, error) {
	url := "https://<databricks-instance>/api/2.0/jobs/runs/submit"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jobPayload)))
	if err != nil {
		return "", err
	}

	// Set headers, including authentication token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer <your-databricks-token>")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result["run_id"].(string), nil
}

// Add-on code: Wait for job to complete
func waitForDatabricksJobCompletion(jobID string) (string, error) {
	for {
		status, err := getDatabricksJobStatus(jobID)
		if err != nil {
			return "", err
		}
		if status == "SUCCESS" || status == "FAILED" {
			return status, nil
		}
		time.Sleep(30 * time.Second) // Poll every 30 seconds
	}
}

func getDatabricksJobStatus(jobID string) (string, error) {
	url := fmt.Sprintf("https://<databricks-instance>/api/2.0/jobs/runs/get?run_id=%s", jobID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer <your-databricks-token>")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result["state"].(map[string]interface{})["life_cycle_state"].(string), nil
}

// Replace the deployment of these serverless functions with submitting Databricks jobs
func deploySubExperimentParallelismInBatches(config *Configuration, serverlessDirPath string, randomExperimentTag string, subExperimentIndex int, functionsPerBatch int) {
	subExperiment := config.SubExperiments[subExperimentIndex]

	numberOfBatches := int(math.Ceil(float64(subExperiment.Parallelism) / float64(functionsPerBatch)))

	for batchNumber := 0; batchNumber < numberOfBatches; batchNumber++ {
		wg := sync.WaitGroup{}

		for parallelism := batchNumber * functionsPerBatch; parallelism < (batchNumber+1)*functionsPerBatch && parallelism < subExperiment.Parallelism; parallelism++ {
			wg.Add(1)

			go func(parallelism int) {
				defer wg.Done()
				// Submit a Databricks job for each function in parallel
				submitDatabricksJob(config, subExperimentIndex, randomExperimentTag, subExperiment)
			}(parallelism)
		}

		wg.Wait() // Wait for all jobs in the batch to finish
	}
}


func ProvisionFunctionsGCR(config *Configuration, serverlessDirPath string) {
	slsConfig := &Serverless{}
	slsConfig.CreateHeaderConfig(config, "STeLLAR-GCR")

	for index, subExperiment := range config.SubExperiments {
		switch subExperiment.PackageType {
		case "Container":
			// size of compressed images of GCR functions on Docker Hub are experimentally found to be approximately 21.84 MiB
			currentSizeInBytes := util.MebibyteToBytes(21.84)
			targetSizeInBytes := util.MebibyteToBytes(subExperiment.FunctionImageSizeMB)
			fillerFileSize := packaging.CalculateFillerFileSizeInBytes(currentSizeInBytes, targetSizeInBytes)
			fillerFilePath := filepath.Join(serverlessDirPath, subExperiment.Function, "filler.file")
			packaging.GenerateFillerFile(subExperiment.ID, fillerFilePath, fillerFileSize)

			imageLink := packaging.SetupContainerImageDeployment(subExperiment.Function, config.Provider, subExperiment.FunctionImageSizeMB)
			randomTag := util.GenerateRandLowercaseLetters(5)
			slsConfig.DeployGCRContainerService(&config.SubExperiments[index], index, randomTag, imageLink, serverlessDirPath, slsConfig.Provider.Region)
		default:
			log.Fatalf("Package type %s is not supported", subExperiment.PackageType)
		}
	}
}

func ProvisionFunctionsCloudflare(config *Configuration, serverlessDirPath string) {
	for index := range config.SubExperiments {
		randomTag := util.GenerateRandLowercaseLetters(5)
		DeployCloudflareWorkers(&config.SubExperiments[index], index, randomTag, serverlessDirPath)
	}
}

func ProvisionFunctionsServerlessAlibaba(config *Configuration, serverlessDirPath string) {
	for index, subExperiment := range config.SubExperiments {
		code_generation.GenerateCode(subExperiment.Function, config.Provider)

		builder := &building.Builder{}
		builder.BuildFunction(config.Provider, subExperiment.Function, subExperiment.Runtime)

		preDeploymentDir := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/sub-experiment-%d", config.Provider, index)
		if err := os.MkdirAll(preDeploymentDir, os.ModePerm); err != nil {
			log.Fatalf("Error creating pre-deployment directory for function %s: %s", subExperiment.Function, err.Error())
		}
		artifactsPath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/artifacts/%s/main.py", config.Provider, subExperiment.Function)
		util.RunCommandAndLog(exec.Command("cp", artifactsPath, preDeploymentDir))

		slsConfig := &Serverless{}
		slsConfig.CreateHeaderConfig(config, fmt.Sprintf("stellar-aliyun-subex%d", index))
		slsConfig.addPlugin("serverless-aliyun-function-compute")
		slsConfig.AddFunctionConfigAlibaba(&config.SubExperiments[index], index, "")
		slsConfig.CreateServerlessConfigFile(fmt.Sprintf("%s/sub-experiment-%d/serverless.yml", serverlessDirPath, index))

		log.Infof("Starting functions deployment. Deploying %d functions to %s.", len(slsConfig.Functions), config.Provider)
		slsDeployMessage := DeployService(fmt.Sprintf("%ssub-experiment-%d", serverlessDirPath, index))

		endpointID := GetAlibabaEndpointID(slsDeployMessage)
		config.SubExperiments[index].AssignEndpointIDs(endpointID)
	}
}
