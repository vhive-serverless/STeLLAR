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

package setup

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"vhive-bench/client/setup/deployment"
	"vhive-bench/client/setup/deployment/connection"
	"vhive-bench/client/util"
)

//Configuration is the schema for all experiment configurations.
type Configuration struct {
	Sequential     bool            `json:"Sequential"`
	Provider       string          `json:"Provider"`
	Runtime        string          `json:"Runtime"`
	SubExperiments []SubExperiment `json:"SubExperiments"`
}

//SubExperiment is the schema for sub-experiment configurations.
type SubExperiment struct {
	Title                   string   `json:"Title"`
	Bursts                  int      `json:"Bursts"`
	BurstSizes              []int    `json:"BurstSizes"`
	PayloadLengthBytes      int      `json:"PayloadLengthBytes"`
	CooldownSeconds         float64  `json:"CooldownSeconds"`
	FunctionIncrementLimits []int64  `json:"FunctionIncrementLimits"`
	DesiredServiceTimes     []string `json:"DesiredServiceTimes"`
	IATType                 string   `json:"IATType"`
	GatewaysNumber          int      `json:"GatewaysNumber"`
	Visualization           string   `json:"Visualization"`
	FunctionMemoryMB        int64    `json:"FunctionMemoryMB"`
	FunctionImageSizeMB     int64    `json:"FunctionImageSizeMB"`
	GatewayEndpoints        []string
	ID                      int
}

const (
	defaultVisualization             = "cdf"
	defaultIATType                   = "stochastic"
	defaultProvider                  = "aws"
	defaultRuntime                   = "go1.x"
	defaultGatewaysNumber            = 1
	defaultFunctionMemoryMB          = 128
	manyRequestsInBurstWarnThreshold = 2000
	manyFilesWarnThreshold           = 500
)

//PrepareSubExperiments will read any required files, deploy functions etc. to get ready for the sub-experiments.
func PrepareSubExperiments(endpointsDirectoryPath string, configPath string) Configuration {
	configFile := util.ReadFile(configPath)
	config := extractConfiguration(configFile)

	transformServiceTimesToFuncIncr(&config)

	connection.Initialize(config.Provider, endpointsDirectoryPath)

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
			config.SubExperiments[index].GatewayEndpoints = []string{config.Provider}
			continue
		}

		availableEndpoints = assignEndpoints(availableEndpoints, &config.SubExperiments[index], config.Provider, config.Runtime)
	}

	return config
}

func extractConfiguration(configFile *os.File) Configuration {
	configByteValue, _ := ioutil.ReadAll(configFile)

	var parsedConfig Configuration
	if err := json.Unmarshal(configByteValue, &parsedConfig); err != nil {
		log.Fatalf("Could not extract experiment configuration from file: %s", err.Error())
	}

	if parsedConfig.Provider == "" {
		parsedConfig.Provider = defaultProvider
	}
	if parsedConfig.Runtime == "" {
		parsedConfig.Runtime = defaultRuntime
	}

	for index := range parsedConfig.SubExperiments {
		if parsedConfig.SubExperiments[index].Visualization == "" {
			parsedConfig.SubExperiments[index].Visualization = defaultVisualization
		}
		if parsedConfig.SubExperiments[index].IATType == "" {
			parsedConfig.SubExperiments[index].IATType = defaultIATType
		}
		if parsedConfig.SubExperiments[index].FunctionMemoryMB == 0 {
			parsedConfig.SubExperiments[index].FunctionMemoryMB = defaultFunctionMemoryMB
		}
		if parsedConfig.SubExperiments[index].GatewaysNumber == 0 {
			parsedConfig.SubExperiments[index].GatewaysNumber = defaultGatewaysNumber
		}
	}

	log.Debugf("Extracted %d sub-experiments from given configuration file.", len(parsedConfig.SubExperiments))
	return parsedConfig
}

func assignEndpoints(availableEndpoints []connection.Endpoint, experiment *SubExperiment, provider string, runtime string) []connection.Endpoint {
	deploymentGeneratedForSubExperiment := false

	var assignedEndpoints []string
	for i := 0; i < experiment.GatewaysNumber; i++ {
		success := false

		for index, endpoint := range availableEndpoints {
			if endpoint.FunctionMemoryMB == experiment.FunctionMemoryMB &&
				(experiment.FunctionImageSizeMB == 0 || almostEqualFloats(endpoint.ImageSizeMB, float64(experiment.FunctionImageSizeMB), 0.5)) {

				assignedEndpoints = append(assignedEndpoints, endpoint.GatewayID)
				availableEndpoints = removeEndpointFromSlice(availableEndpoints, index)
				success = true
				break
			}
		}

		if !success {
			log.Infof("Could not find a function to assign with memory %dMB and image size %dMB, deploying...",
				experiment.FunctionMemoryMB,
				experiment.FunctionImageSizeMB,
			)

			if !deploymentGeneratedForSubExperiment {
				deployment.SetupDeployment(provider, runtime, util.MBToBytes(float64(experiment.FunctionImageSizeMB)))
				deploymentGeneratedForSubExperiment = true
			}
			assignedEndpoints = append(assignedEndpoints, connection.Singleton.DeployFunction(runtime, 128))

			//TODO: intelligently leverage connection.Singleton.RemoveFunction(uniqueID) &
			//TODO: connection.Singleton.UpdateFunction(uniqueID, 128)
			//TODO: once over 600 deployed functions
		}
	}

	log.Debugf("Assigning following endpoints to sub-experiment `%s`: %v", experiment.Title, assignedEndpoints)
	experiment.GatewayEndpoints = assignedEndpoints
	return availableEndpoints
}
