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
	"os"
	"os/exec"
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

	availableEndpoints := connection.Singleton.ListAPIs()

	for index := 0; index < len(config.SubExperiments); index++ {
		subExperiment := &config.SubExperiments[index]
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

	slsConfig.CreateHeaderConfig(config, "STeLLAR")
	slsConfig.packageIndividually()

	for index := 0; index < len(config.SubExperiments); index++ {
		subExperiment := &config.SubExperiments[index]
		//TODO: generate the code
		code_generation.GenerateCode(subExperiment.Function, config.Provider)

		// TODO: build the functions (Java and Golang)
		artifactPathRelativeToServerlessConfigFile := builder.BuildFunction(config.Provider, subExperiment.Function, subExperiment.Runtime)
		slsConfig.AddFunctionConfigAWS(&config.SubExperiments[index], index, artifactPathRelativeToServerlessConfigFile)

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

func ProvisionFunctionsServerlessAzure(config *Configuration, serverlessDirPath string) {
	wg := sync.WaitGroup{}

	randomExperimentTag := util.GenerateRandomLowercaseLetters(5)

	for index := 0; index < len(config.SubExperiments); index++ {
		subExperiment := &config.SubExperiments[index]
		code_generation.GenerateCode(subExperiment.Function, config.Provider)

		builder := &building.Builder{}
		builder.BuildFunction(config.Provider, subExperiment.Function, subExperiment.Runtime)

		if config.SubExperiments[index].Endpoints == nil {
			config.SubExperiments[index].Endpoints = []EndpointInfo{}
		}

		for parallelism := 0; parallelism < subExperiment.Parallelism; parallelism++ {
			deploymentDir := fmt.Sprintf("%ssub-experiment-%d/parallelism-%d", serverlessDirPath, index, parallelism)
			if err := os.MkdirAll(deploymentDir, os.ModePerm); err != nil {
				log.Fatalf("Error creating pre-deployment directory for function %s: %s", subExperiment.Function, err.Error())
			}
			artifactsPath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/artifacts/%s/main.py", config.Provider, subExperiment.Function)
			util.RunCommandAndLog(exec.Command("cp", artifactsPath, deploymentDir))

			slsConfig := &Serverless{}
			slsConfig.CreateHeaderConfig(config, fmt.Sprintf("stellar-%s-subex%d-para%d", randomExperimentTag, index, parallelism))
			slsConfig.addPlugin("serverless-azure-functions")
			slsConfig.AddFunctionConfigAzure(subExperiment, index, parallelism)
			slsConfig.CreateServerlessConfigFile(fmt.Sprintf("%s/serverless.yml", deploymentDir))

			log.Infof("Starting functions deployment. Deploying %d functions to %s.", len(slsConfig.Functions), config.Provider)
			wg.Add(1)
			log.Infof("Deploying subex %d", index)
			go func(config *Configuration, index int, deploymentDir string) {
				defer wg.Done()
				deployServiceAndSaveEndpointId(config, index, deploymentDir)
			}(config, index, deploymentDir)
		}
	}

	wg.Wait()
}

func deployServiceAndSaveEndpointId(config *Configuration, subExperimentIndex int, deploymentDir string) {
	log.Infof("insinde goroutine Deploying subex %d", subExperimentIndex)
	slsDeployMessage := DeployService(deploymentDir)
	endpointID := GetAzureEndpointID(slsDeployMessage)
	config.SubExperiments[subExperimentIndex].mu.Lock()
	defer config.SubExperiments[subExperimentIndex].mu.Unlock()
	config.SubExperiments[subExperimentIndex].Endpoints = append(config.SubExperiments[subExperimentIndex].Endpoints, EndpointInfo{ID: endpointID})
}

func ProvisionFunctionsGCR(config *Configuration, serverlessDirPath string) {
	slsConfig := &Serverless{}
	slsConfig.CreateHeaderConfig(config, "STeLLAR-GCR")

	for index := 0; index < len(config.SubExperiments); index++ {
		subExperiment := &config.SubExperiments[index]
		switch subExperiment.PackageType {
		case "Container":
			imageLink := packaging.SetupContainerImageDeployment(subExperiment.Function, config.Provider)
			slsConfig.DeployGCRContainerService(&config.SubExperiments[index], index, imageLink, serverlessDirPath, slsConfig.Provider.Region)
		default:
			log.Fatalf("Package type %s is not supported", subExperiment.PackageType)
		}
	}
}

func ProvisionFunctionsCloudflare(config *Configuration, serverlessDirPath string) {
	for index := range config.SubExperiments {
		DeployCloudflareWorkers(&config.SubExperiments[index], index, serverlessDirPath)
	}
}

func ProvisionFunctionsServerlessAlibaba(config *Configuration, serverlessDirPath string) {
	for index := 0; index < len(config.SubExperiments); index++ {
		subExperiment := &config.SubExperiments[index]
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
