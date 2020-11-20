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
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"lambda-benchmarking/client/prompts"
	"os"
)

//Endpoint is the schema for the configuration of provider endpoints.
type Endpoint struct {
	GatewayID        string  `json:"GatewayID"`
	FunctionMemoryMB int64   `json:"FunctionMemoryMB"`
}

func extractEndpoints(endpointsFile *os.File) []Endpoint {
	configByteValue, _ := ioutil.ReadAll(endpointsFile)

	var parsedEndpoints []Endpoint
	if err := json.Unmarshal(configByteValue, &parsedEndpoints); err != nil {
		log.Fatalf("Could not extract endpoints configuration from file: %s", err.Error())
	}

	return parsedEndpoints
}

func assignEndpoints(gateways map[int64][]string, memoryToLastAssignedIndex map[int64]int, experiment *SubExperiment) {
	lastAssignedIndexExcl := memoryToLastAssignedIndex[experiment.FunctionMemoryMB]
	newLastAssignedIndexExcl := lastAssignedIndexExcl + experiment.GatewaysNumber
	nrGatewaysWithDesiredMemory := len(gateways[experiment.FunctionMemoryMB])

	if newLastAssignedIndexExcl > nrGatewaysWithDesiredMemory {
		remainingGatewaysToAssign := nrGatewaysWithDesiredMemory - lastAssignedIndexExcl
		log.Errorf("Not enough remaining gateways were found in the given gateways file with requested memory %dMB, found %d but trying to assign from %d to %d. Experiment `%s` will be assigned %d gateways.",
			experiment.FunctionMemoryMB, remainingGatewaysToAssign,
			lastAssignedIndexExcl, newLastAssignedIndexExcl,
			experiment.Title, remainingGatewaysToAssign)

		if remainingGatewaysToAssign <= 0 {
			log.Fatalf("Cannot assign %d gateways to an experiment.", remainingGatewaysToAssign)
		}

		if !prompts.PromptForBool("Would you like to continue with this setting?") {
			os.Exit(0)
		}

		newLastAssignedIndexExcl = nrGatewaysWithDesiredMemory
	}
	memoryToLastAssignedIndex[experiment.FunctionMemoryMB] = newLastAssignedIndexExcl
	experiment.GatewayEndpoints = gateways[experiment.FunctionMemoryMB][lastAssignedIndexExcl:newLastAssignedIndexExcl]
}

func mapMemoryToGateways(parsedEndpoints []Endpoint) (map[int64][]string, map[int64]int) {
	memoryToListOfGatewayIDs := make(map[int64][]string)
	memoryToLastAssignedIndex := make(map[int64]int)
	for idx, endpoint := range parsedEndpoints {
		if idx == 0 {
			continue
		}

		memoryToLastAssignedIndex[endpoint.FunctionMemoryMB] = 0
		memoryToListOfGatewayIDs[endpoint.FunctionMemoryMB] =
			append(memoryToListOfGatewayIDs[endpoint.FunctionMemoryMB], endpoint.GatewayID)
	}
	return memoryToListOfGatewayIDs, memoryToLastAssignedIndex
}
