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
	"os"
)

func ReadInstructions(endpointsPath string, configPath string) Configuration {
	log.Debugf("Reading endpoints file for this run from `%s`", endpointsPath)
	endpointsFile, err := os.Open(endpointsPath)
	if err != nil {
		log.Fatalf("Could not read endpoints file: %s", err.Error())
	}
	availableEndpoints := extractEndpoints(endpointsFile)
	// TODO: remove this I think
	memoryToGatewayIDs, memoryToLastAssignedIndex := mapMemoryToGateways(availableEndpoints)

	log.Debugf("Reading config file for this run from `%s`", configPath)
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Could not read config file: %s", err.Error())
	}
	config := extractSubExperiments(configFile)

	for index := range config.SubExperiments {
		config.SubExperiments[index].ID = index
		assignEndpoints(memoryToGatewayIDs, memoryToLastAssignedIndex, &config.SubExperiments[index])

		// Issue warning if sending too many requests in a single burst
		const manyRequestsInBurstThreshold = 2000
		for _, burstSize := range config.SubExperiments[index].BurstSizes {
			if burstSize > manyRequestsInBurstThreshold {
				log.Warnf("Experiment %d has a burst of size %d, NIC (Network Interface Controller) contention may occur.",
					index, burstSize)
			}
		}

		// Issue warning if generating too many files
		const manyFilesWarningThreshold = 500
		chosenVisualization := config.SubExperiments[index].Visualization
		burstsNumber := config.SubExperiments[index].Bursts
		if burstsNumber >= manyFilesWarningThreshold && (chosenVisualization == "all" || chosenVisualization == "histogram") {
			log.Warnf("Generating histograms for each burst, this will create a large number (%d) of new files.",
				burstsNumber)
		}
	}

	log.Debugf("Extracted %d sub-experiments from given configuration file.", len(config.SubExperiments))
	return config
}
