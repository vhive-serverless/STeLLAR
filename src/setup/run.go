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
	log "github.com/sirupsen/logrus"
	"os"
	"stellar/setup/building"
	code_generation "stellar/setup/code-generation"
	"stellar/setup/deployment/connection"
	"stellar/setup/deployment/connection/amazon"
	"stellar/setup/deployment/packaging"
	"time"
)

// ProvisionFunctions will deploy, reconfigure, etc. functions to get ready for the sub-experiments.
func ProvisionFunctions(config Configuration) {
	const (
		nicContentionWarnThreshold = 800 // Experimentally found
		storageSpaceWarnThreshold  = 500 // 500 * ~18KiB = 10MB just for 1 sub-experiment
	)

	availableEndpoints := connection.Singleton.ListAPIs()

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
func ProvisionFunctionsServerless(config Configuration, serverlessDirPath string) {

	slsConfig := &Serverless{}
	builder := &building.Builder{}

	slsConfig.CreateHeaderConfig(&config)

	for index, subExperiment := range config.SubExperiments {
		//TODO: generate the code
		code_generation.GenerateCode(subExperiment.Function, config.Provider)

		// build the functions (Java and Golang)
		artifactPathRelativeToServerlessConfigFile := builder.BuildFunction(config.Provider, subExperiment.Function, subExperiment.Runtime)
		slsConfig.AddFunctionConfig(&subExperiment, index, artifactPathRelativeToServerlessConfigFile)

		// generate filler files and zip used as Serverless artifacts
		packaging.GenerateServerlessZIPArtifacts(subExperiment.ID, config.Provider, subExperiment.Runtime, subExperiment.Function, subExperiment.FunctionImageSizeMB)
	}

	slsConfig.CreateServerlessConfigFile(fmt.Sprintf("%sserverless.yml", serverlessDirPath))

	log.Infof("Starting functions deployment. Deploying %d functions to %s.", len(slsConfig.Functions), config.Provider)
	slsDeployMessage := DeployService(serverlessDirPath)
	log.Info(slsDeployMessage)

	// TODO: assign endpoints to subexperiments
	// Get the endpoints by scraping the serverless deploy message.

}
