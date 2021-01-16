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
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"vhive-bench/client/setup/deployment"
	"vhive-bench/client/setup/deployment/connection"
	"vhive-bench/client/util"
)

func assignEndpoints(availableEndpoints []connection.Endpoint, experiment *SubExperiment, provider string,
	runtime string) []connection.Endpoint {
	deploymentGeneratedForSubExperiment := false

	var assignedEndpoints []string
	for i := 0; i < experiment.GatewaysNumber; i++ {
		if canAssignToExistingEndpoint(&availableEndpoints, experiment, &assignedEndpoints) {
			continue
		}

		log.Infof("Searched %d endpoints, could not find a function to assign with memory %dMB and image size %vMB.",
			len(availableEndpoints),
			experiment.FunctionMemoryMB,
			experiment.FunctionImageSizeMB)

		if !deploymentGeneratedForSubExperiment {
			log.Info("Setting up deployment...")
			experiment.FunctionImageSizeMB = deployment.SetupDeployment(
				fmt.Sprintf("setup/deployment/raw-code/%s/%s-handler/main.go", runtime, provider),
				provider,
				runtime,
				util.MBToBytes(experiment.FunctionImageSizeMB),
				experiment.PackageType,
			)
			deploymentGeneratedForSubExperiment = true
		}

		if canAssignToRepurposedEndpoint(&availableEndpoints, experiment, &assignedEndpoints, experiment.PackageType) {
			continue
		}

		log.Info("Could not find an existing function to repurpose, creating a new function...")
		assignedEndpoints = append(assignedEndpoints,
			connection.Singleton.DeployFunction(experiment.PackageType, runtime, experiment.FunctionMemoryMB))
	}

	log.Debugf("Assigning following endpoints to sub-experiment `%s`: %v", experiment.Title, assignedEndpoints)
	experiment.GatewayEndpoints = assignedEndpoints
	return availableEndpoints
}

func canAssignToRepurposedEndpoint(availableEndpoints *[]connection.Endpoint, experiment *SubExperiment,
	assignedEndpoints *[]string, packageType string) bool {
	for index, endpoint := range *availableEndpoints {
		// Can only repurpose function of same package type.
		if endpoint.PackageType == experiment.PackageType {
			log.Info("Repurposing an existing function...")
			connection.Singleton.UpdateFunction(packageType, endpoint.GatewayID, experiment.FunctionMemoryMB)
			*assignedEndpoints = append(*assignedEndpoints, endpoint.GatewayID)
			*availableEndpoints = removeEndpointFromSlice(*availableEndpoints, index)
			log.Infof("Successfully repurposed %q (memory %dMB -> %dMB, image size %vMB -> %vMB).",
				endpoint.GatewayID, endpoint.FunctionMemoryMB, experiment.FunctionMemoryMB,
				endpoint.ImageSizeMB, experiment.FunctionImageSizeMB)
			return true
		}
	}
	return false
}

func canAssignToExistingEndpoint(availableEndpoints *[]connection.Endpoint, experiment *SubExperiment, assignedEndpoints *[]string) bool {
	foundEndpoint := false
	for index, endpoint := range *availableEndpoints {
		if specsMatch(endpoint, experiment) {
			*assignedEndpoints = append(*assignedEndpoints, endpoint.GatewayID)
			*availableEndpoints = removeEndpointFromSlice(*availableEndpoints, index)
			foundEndpoint = true
			break
		}
	}
	return foundEndpoint
}

func removeEndpointFromSlice(s []connection.Endpoint, i int) []connection.Endpoint {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func specsMatch(endpoint connection.Endpoint, experiment *SubExperiment) bool {
	return endpoint.FunctionMemoryMB == experiment.FunctionMemoryMB &&
		endpoint.PackageType == experiment.PackageType &&
		(experiment.FunctionImageSizeMB == 0 ||
			math.Abs(endpoint.ImageSizeMB-experiment.FunctionImageSizeMB) <= 0.5)
}
