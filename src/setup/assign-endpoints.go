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
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
	"math"
	"vhive-bench/setup/deployment"
	"vhive-bench/setup/deployment/connection"
	"vhive-bench/util"
)

func assignEndpoints(availableEndpoints []connection.Endpoint, experiment *SubExperiment, provider string, runtime string) []connection.Endpoint {
	log.Infof("[sub-experiment %d] Setting up deployment...", experiment.ID)
	var assignedBinaryPath string

	if provider != "vhive" { // cannot deploy to vhive
		experiment.FunctionImageSizeMB, assignedBinaryPath = deployment.SetupDeployment(
			fmt.Sprintf("setup/deployment/raw-code/%s/%s/%s/", experiment.Function, runtime, provider),
			provider,
			util.MBToBytes(experiment.FunctionImageSizeMB),
			experiment.PackageType,
			experiment.ID,
			experiment.Function,
		)
	}

	var assignedEndpoints []EndpointInfo
	for i := 0; i < experiment.Parallelism; i++ {
		foundEndpointID := findEndpointToAssign(&availableEndpoints, experiment, assignedBinaryPath, runtime)

		gatewayEndpoint := EndpointInfo{ID: foundEndpointID}

		for j := experiment.DataTransferChainLength; j > 1; j-- {
			gatewayEndpoint.DataTransferChainIDs = append(
				gatewayEndpoint.DataTransferChainIDs,
				findEndpointToAssign(&availableEndpoints, experiment, assignedBinaryPath, runtime),
			)
		}

		assignedEndpoints = append(assignedEndpoints, gatewayEndpoint)
	}

	log.Debugf("[sub-experiment %d] Assigning following endpoints: %v", experiment.ID, assignedEndpoints)
	experiment.Endpoints = assignedEndpoints
	return availableEndpoints
}

func findEndpointToAssign(availableEndpoints *[]connection.Endpoint, experiment *SubExperiment, binaryPath string, runtime string) string {
	for index, endpoint := range *availableEndpoints {
		if specsMatch(endpoint, experiment) {
			*availableEndpoints = removeEndpointFromSlice(*availableEndpoints, index)
			return endpoint.GatewayID
		}
	}

	log.Infof("[sub-experiment %d] Searched %d endpoints, could not find a function to assign with: memory %dMB, image size %vMB, package type %q.",
		experiment.ID,
		len(*availableEndpoints),
		experiment.FunctionMemoryMB,
		experiment.FunctionImageSizeMB,
		experiment.PackageType,
	)

	// Only attempt repurposing functions if they are ZIP-packaged:
	// https://github.com/motdotla/node-lambda/issues/535 (Image-packaged functions have update errors on AWS)
	if experiment.PackageType == "Zip" {
		for index, endpoint := range *availableEndpoints {
			if endpoint.PackageType == "Zip" {
				log.Infof("[sub-experiment %d] Repurposing an existing function...", experiment.ID)
				connection.Singleton.UpdateFunction(experiment.PackageType, endpoint.GatewayID, experiment.FunctionMemoryMB)

				*availableEndpoints = removeEndpointFromSlice(*availableEndpoints, index)

				log.Infof("[sub-experiment %d] Successfully repurposed %q (memory %dMB -> %dMB, image size %vMB -> %vMB).",
					experiment.ID,
					endpoint.GatewayID,
					endpoint.FunctionMemoryMB,
					experiment.FunctionMemoryMB,
					endpoint.ImageSizeMB,
					experiment.FunctionImageSizeMB,
				)
				return endpoint.GatewayID
			}
		}
	}

	log.Infof("[sub-experiment %d] Could not find an existing function to repurpose, creating a new function...", experiment.ID)
	return connection.Singleton.DeployFunction(binaryPath, experiment.PackageType, runtime, experiment.FunctionMemoryMB)
}

func removeEndpointFromSlice(s []connection.Endpoint, i int) []connection.Endpoint {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func specsMatch(endpoint connection.Endpoint, experiment *SubExperiment) bool {
	if experiment.PackageType != endpoint.PackageType {
		return false
	}

	if endpoint.FunctionMemoryMB != experiment.FunctionMemoryMB {
		return false
	}

	// Image sizes are ignored for PackageTypeImage because AWS does not reveal this information
	if experiment.PackageType == lambda.PackageTypeImage {
		return true
	}

	return math.Abs(endpoint.ImageSizeMB-experiment.FunctionImageSizeMB) <= 5
}
