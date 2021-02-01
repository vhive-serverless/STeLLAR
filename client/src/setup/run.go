// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
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
	"time"
	"vhive-bench/client/setup/deployment/connection"
	"vhive-bench/client/setup/deployment/connection/amazon"
	"vhive-bench/client/util"
)

//GatewayEndpoint represents the initial endpoint ID together with the IDs of lambda functions further in the data transfer chain
type GatewayEndpoint struct {
	ID                   string
	DataTransferChainIDs []string
}

//SubExperiment is the schema for sub-experiment configurations.
type SubExperiment struct {
	Title                   string   `json:"Title"`
	Bursts                  int      `json:"Bursts"`
	BurstSizes              []int    `json:"BurstSizes"`
	PayloadLengthBytes      int      `json:"PayloadLengthBytes"`
	IATSeconds              float64  `json:"IATSeconds"`
	FunctionIncrementLimits []int64  `json:"FunctionIncrementLimits"`
	DesiredServiceTimes     []string `json:"DesiredServiceTimes"`
	IATType                 string   `json:"IATType"`
	PackageType             string   `json:"PackageType"`
	GatewaysNumber          int      `json:"GatewaysNumber"`
	Visualization           string   `json:"Visualization"`
	FunctionMemoryMB        int64    `json:"FunctionMemoryMB"`
	FunctionImageSizeMB     float64  `json:"FunctionImageSizeMB"`
	DataTransferChainLength int      `json:"DataTransferChainLength"`
	GatewayEndpoints        []GatewayEndpoint
	ID                      int
}

const (
	defaultVisualization             = "cdf"
	defaultIATType                   = "stochastic"
	defaultProvider                  = "aws"
	defaultRuntime                   = "go1.x"
	defaultPackageType               = "Zip"
	defaultGatewaysNumber            = 1
	defaultDataTransferChainLength   = 1
	defaultFunctionMemoryMB          = 128
	manyRequestsInBurstWarnThreshold = 2000
	manyFilesWarnThreshold           = 500
)

//PrepareSubExperiments will read any required files, deploy functions etc. to get ready for the sub-experiments.
func PrepareSubExperiments(endpointsDirectoryPath string, configPath string) Configuration {
	configFile := util.ReadFile(configPath)
	config := extractConfiguration(configFile)

	transformServiceTimesToFunctionIncrements(&config)

	connection.Initialize(config.Provider, endpointsDirectoryPath, "./setup/deployment/raw-code/producer-consumer/api-template.json")

	availableEndpoints := connection.Singleton.ListAPIs()

	for index, subExperiment := range config.SubExperiments {
		config.SubExperiments[index].ID = index

		for _, burstSize := range subExperiment.BurstSizes {
			if burstSize > manyRequestsInBurstWarnThreshold {
				log.Warnf("Experiment %d has a burst of size %d, NIC (Network Interface Controller) contention may occur.",
					index, burstSize)
				if !promptForBool("Do you wish to continue?") {
					os.Exit(0)
				}
			}
		}

		if subExperiment.Bursts >= manyFilesWarnThreshold &&
			(subExperiment.Visualization == "all" || subExperiment.Visualization == "histogram") {
			log.Warnf("SubExperiment %d is generating histograms for each burst, this will create a large number (%d) of new files.",
				index, subExperiment.Bursts)
			if !promptForBool("Do you wish to continue?") {
				os.Exit(0)
			}
		}

		if availableEndpoints == nil { // hostname must be the endpoint itself (external URL)
			config.SubExperiments[index].GatewayEndpoints = []GatewayEndpoint{{ID: config.Provider}}
			continue
		}

		availableEndpoints = assignEndpoints(availableEndpoints, &config.SubExperiments[index], config.Provider, config.Runtime)
	}

	if amazon.AWSSingletonInstance != nil && amazon.AWSSingletonInstance.ImageURI != "" {
		log.Info("A deployment was made using container images, waiting 10 seconds for changes to take effect with the provider...")
		time.Sleep(time.Second * 10)
	}

	return config
}
