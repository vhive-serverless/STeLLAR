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
	ImageSizeMB      float64 `json:"ImageSizeMB"`
}

func extractProviderEndpoints(endpointsFile *os.File) []Endpoint {
	configByteValue, _ := ioutil.ReadAll(endpointsFile)

	var parsedEndpoints []Endpoint
	if err := json.Unmarshal(configByteValue, &parsedEndpoints); err != nil {
		log.Fatalf("Could not extract endpoints configuration from file: %s", err.Error())
	}

	return parsedEndpoints
}

func assignEndpoints(availableEndpoints []Endpoint, experiment *SubExperiment) {
	if experiment.GatewaysNumber > len(availableEndpoints) {
		log.Fatalf("Cannot assign %d endpoints to experiment %d: only %d available endpoints for provider %s.",
			experiment.GatewaysNumber,
			experiment.ID,
			len(availableEndpoints),
			experiment.Provider,
		)
	}

	var assignedEndpoints []string
	for i := 0; i < experiment.GatewaysNumber; i++ {
		for index, availableEndpoint := range availableEndpoints {
			if availableEndpoint.FunctionMemoryMB == experiment.FunctionMemoryMB &&
				(experiment.FunctionImageSizeMB == 0 || almostEqual(availableEndpoint.ImageSizeMB,
					float64(experiment.FunctionImageSizeMB),
					0.5)) {
				assignedEndpoints = append(assignedEndpoints, availableEndpoint.GatewayID)
				removeEndpoint(availableEndpoints, index)
			}
		}

		log.Warnf("Could not find a function to assign with %dMB memory and %dMB image size.",
			experiment.FunctionMemoryMB,
			experiment.FunctionImageSizeMB,
		)
		if !prompts.PromptForBool("Do you wish to continue?") {
			os.Exit(0)
		}
	}

	experiment.GatewayEndpoints = assignedEndpoints
}

func removeEndpoint(s []Endpoint, i int) []Endpoint {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
