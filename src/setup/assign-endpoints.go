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

package setup

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
	"math"
	"stellar/setup/deployment"
	"stellar/setup/deployment/connection"
	"stellar/util"
)

func assignEndpoints(availableEndpoints []connection.Endpoint, experiment *SubExperiment, provider string) []connection.Endpoint {
	log.Infof("[sub-experiment %d] Setting up deployment...", experiment.ID)
	log.Infof("[sub-experiment %d] Experiment configuration: %vMB memory, %vMB image size, %vs IAT, %q package.",
		experiment.ID, experiment.FunctionMemoryMB, experiment.FunctionImageSizeMB, experiment.IATSeconds,
		experiment.PackageType)
	var assignedHandler string

	if provider == "aws" { // deployment has only been automated for AWS so far
		experiment.FunctionImageSizeMB, assignedHandler = deployment.SetupDeployment(
			fmt.Sprintf("setup/deployment/raw-code/functions/%s/%s", experiment.Function, provider),
			provider,
			util.MebibyteToBytes(experiment.FunctionImageSizeMB),
			experiment.PackageType,
			experiment.ID,
			experiment.Function,
		)
	}

	var assignedEndpoints []EndpointInfo
	for i := 0; i < experiment.Parallelism; i++ {
		foundEndpointID := findEndpointToAssign(&availableEndpoints, experiment, assignedHandler)

		gatewayEndpoint := EndpointInfo{ID: foundEndpointID}

		for j := experiment.DataTransferChainLength; j > 1; j-- {
			gatewayEndpoint.DataTransferChainIDs = append(
				gatewayEndpoint.DataTransferChainIDs,
				findEndpointToAssign(&availableEndpoints, experiment, assignedHandler),
			)
		}

		assignedEndpoints = append(assignedEndpoints, gatewayEndpoint)
	}

	log.Debugf("[sub-experiment %d] Assigning following endpoints: %v", experiment.ID, assignedEndpoints)
	experiment.Endpoints = assignedEndpoints
	return availableEndpoints
}

func findEndpointToAssign(availableEndpoints *[]connection.Endpoint, experiment *SubExperiment, assignedHandler string) string {
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
	return connection.Singleton.DeployFunction(assignedHandler, experiment.PackageType, experiment.Function, experiment.FunctionMemoryMB)
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

// AssignEndpointIDs assigns a given endpoint to all deployed functions of the subexperiment.
func (s *SubExperiment) AssignEndpointIDs(endpointID string) {
	if s.Endpoints == nil {
		s.Endpoints = []EndpointInfo{}
	}
	for i := 0; i < s.Parallelism; i++ {
		s.Endpoints = append(s.Endpoints, EndpointInfo{ID: endpointID})
	}
}

func (s *SubExperiment) AddRoute(path string) {
	s.Routes = append(s.Routes, path)
}
