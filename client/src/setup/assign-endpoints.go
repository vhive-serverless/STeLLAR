package setup

import (
	log "github.com/sirupsen/logrus"
	"vhive-bench/client/setup/deployment"
	"vhive-bench/client/setup/deployment/connection"
	"vhive-bench/client/util"
)

func assignEndpoints(availableEndpoints []connection.Endpoint, experiment *SubExperiment, provider string, runtime string) []connection.Endpoint {
	deploymentGeneratedForSubExperiment := false

	var assignedEndpoints []string
	for i := 0; i < experiment.GatewaysNumber; i++ {
		foundEndpoint := false
		for index, endpoint := range availableEndpoints {
			if specsMatch(endpoint, experiment) {
				assignedEndpoints = append(assignedEndpoints, endpoint.GatewayID)
				availableEndpoints = removeEndpointFromSlice(availableEndpoints, index)
				foundEndpoint = true
				break
			}
		}
		if foundEndpoint {
			continue
		}

		log.Infof("Searched %d endpoints, could not find a function to assign with memory %dMB and image size %vMB.",
			len(availableEndpoints),
			experiment.FunctionMemoryMB,
			experiment.FunctionImageSizeMB)

		if !deploymentGeneratedForSubExperiment {
			log.Info("Setting up deployment...")
			experiment.FunctionImageSizeMB = deployment.SetupDeployment(provider, runtime, util.MBToBytes(experiment.FunctionImageSizeMB))
			deploymentGeneratedForSubExperiment = true
		}

		repurposedEndpoint := false
		for index, endpoint := range availableEndpoints {
			log.Info("Repurposing an existing function...")
			connection.Singleton.UpdateFunction(endpoint.GatewayID, int(experiment.FunctionMemoryMB))
			assignedEndpoints = append(assignedEndpoints, endpoint.GatewayID)
			availableEndpoints = removeEndpointFromSlice(availableEndpoints, index)
			log.Infof("Successfully repurposed %q (memory %dMB -> %dMB, image size %vMB -> %vMB).",
				endpoint.GatewayID, endpoint.FunctionMemoryMB, experiment.FunctionMemoryMB,
				endpoint.ImageSizeMB, experiment.FunctionImageSizeMB)
			repurposedEndpoint = true
			break
		}
		if repurposedEndpoint {
			continue
		}

		log.Info("Could not find an existing function to repurpose, creating a new function...")
		assignedEndpoints = append(assignedEndpoints, connection.Singleton.DeployFunction(runtime, int(experiment.FunctionMemoryMB)))
	}

	log.Debugf("Assigning following endpoints to sub-experiment `%s`: %v", experiment.Title, assignedEndpoints)
	experiment.GatewayEndpoints = assignedEndpoints
	return availableEndpoints
}

func removeEndpointFromSlice(s []connection.Endpoint, i int) []connection.Endpoint {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func specsMatch(endpoint connection.Endpoint, experiment *SubExperiment) bool {
	return endpoint.FunctionMemoryMB == experiment.FunctionMemoryMB &&
		(experiment.FunctionImageSizeMB == 0 || util.AlmostEqualFloats(endpoint.ImageSizeMB, experiment.FunctionImageSizeMB, 0.5))
}
