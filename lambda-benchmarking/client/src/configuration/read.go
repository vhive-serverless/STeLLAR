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

package configuration

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/prompts"
	"os"
	"path"
)

// HashMap pointing from provider (e.g., "aws") to a list of existing endpoints for that provider
var availableEndpoints map[string][]Endpoint

//ReadInstructions will read any required files, such as experiment configurations and endpoint files.
func ReadInstructions(endpointsDirectoryPath string, configPath string) Configuration {
	supportedProviders := []string{"aws", "vHive"}

	configFile := readFile(configPath)

	config := extractSubExperiments(configFile)
	log.Debugf("Extracted %d sub-experiments from given configuration file.", len(config.SubExperiments))

	availableEndpoints = make(map[string][]Endpoint)
	for index, subExperiment := range config.SubExperiments {
		_, alreadyExtractedProvider := availableEndpoints[subExperiment.Provider]
		if !stringInSlice(subExperiment.Provider, supportedProviders) {
			config.SubExperiments[index].GatewayEndpoints = []string{subExperiment.Provider}
		} else if !alreadyExtractedProvider {
			endpointsFile := readFile(path.Join(endpointsDirectoryPath, subExperiment.Provider + ".json"))
			availableEndpoints[subExperiment.Provider] = extractProviderEndpoints(endpointsFile)
			assignEndpoints(availableEndpoints[subExperiment.Provider], &config.SubExperiments[index])
		}

		config.SubExperiments[index].ID = index

		const manyRequestsInBurstThreshold = 2000
		for _, burstSize := range config.SubExperiments[index].BurstSizes {
			if burstSize > manyRequestsInBurstThreshold {
				log.Warnf("Experiment %d has a burst of size %d, NIC (Network Interface Controller) contention may occur.",
					index, burstSize)
				if !prompts.PromptForBool("Do you wish to continue?") {
					os.Exit(0)
				}
			}
		}

		const manyFilesWarningThreshold = 500
		chosenVisualization := config.SubExperiments[index].Visualization
		burstsNumber := config.SubExperiments[index].Bursts
		if burstsNumber >= manyFilesWarningThreshold && (chosenVisualization == "all" || chosenVisualization == "histogram") {
			log.Warnf("SubExperiment %d is generating histograms for each burst, this will create a large number (%d) of new files.",
				index, burstsNumber)
			if !prompts.PromptForBool("Do you wish to continue?") {
				os.Exit(0)
			}
		}
	}

	return config
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func readFile(path string) *os.File {
	log.Debugf("Reading file for this run from `%s`", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not read file: %s", err.Error())
	}
	return file
}
