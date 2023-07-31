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
	log "github.com/sirupsen/logrus"
	"os"
	"stellar/setup/deployment/connection/amazon"
	"strconv"
	"strings"
	"time"
)

//ProvisionFunctions will deploy, reconfigure, etc. functions to get ready for the sub-experiments.
func ProvisionFunctions(config *Configuration) {
	const (
		nicContentionWarnThreshold = 800 // Experimentally found
		storageSpaceWarnThreshold  = 500 // 500 * ~18KiB = 10MB just for 1 sub-experiment
	)

	//availableEndpoints := connection.Singleton.ListAPIs()

	slsConfig := &Serverless{}

	slsConfig.CreateHeader(*config)
	slsConfig.AddPackagePattern("!**")

	// Create a serverless.yml function configurations for all functions
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

		slsConfig.AddFunctionConfig(&config.SubExperiments[index], index)
		////  no clue what this does
		//if availableEndpoints == nil { // hostname must be the endpoint itself (external URL)
		//	config.SubExperiments[index].Endpoints = []EndpointInfo{{ID: config.Provider}}
		//	continue
		//}
		//
		// availableEndpoints = assignEndpoints(
		//	availableEndpoints,
		//	&config.SubExperiments[index],
		//	config.Provider,
		//)
	}

	log.Infof("Creating serverless.yml.")
	slsConfig.CreateServerlessConfigFile()

	log.Infof("Starting serverless.com deployment.")
	slsDeployMessage := deployService()
	log.Infof(slsDeployMessage)

	endpointID := getEndpointID(slsDeployMessage)

	// Assign Ednpoint ID to each deployed function
	for i := range config.SubExperiments {
		assignEndpointIDs(endpointID, &config.SubExperiments[i])
		log.Infof(strconv.Itoa(len(config.SubExperiments)))
	}

	if amazon.AWSSingletonInstance != nil && amazon.AWSSingletonInstance.ImageURI != "" {
		log.Info("A deployment was made using container images, waiting 10 seconds for changes to take effect with the provider...")
		time.Sleep(time.Second * 10)

	}
}

func assignEndpointIDs(endpointID string, subex *SubExperiment) {
	subex.Endpoints = []EndpointInfo{}
	for i := 0; i < subex.Parallelism; i++ {
		subex.Endpoints = append(subex.Endpoints, EndpointInfo{ID: endpointID})
	}
	log.Infof(strconv.Itoa(len(subex.Endpoints)))
}

// getEndpointID scrapes the serverless deploy message for the endpoint ID
func getEndpointID(slsDeployMessage string) string {
	lines := strings.Split(slsDeployMessage, "\n")
	if lines[1] == "endpoints:" {
		line := lines[2]
		link := strings.Split(line, " ")[4]
		httpId := strings.Split(link, ".")[0]
		endpointId := strings.Split(httpId, "//")[1]
		return endpointId
	}
	line := lines[1]
	link := strings.Split(line, " ")[3]
	httpId := strings.Split(link, ".")[0]
	endpointId := strings.Split(httpId, "//")[1]
	return endpointId
}
