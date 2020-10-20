package benchmarking

import (
	"fmt"
	"lambda-benchmarking/client/experiment/configuration"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	region = "us-west-1"
)

// TriggerRelativeAsyncBursts
func RunProfiler(config configuration.ExperimentConfig, deltas []time.Duration, safeExperimentWriter *SafeWriter) {
	estimateTime := estimateTotalDuration(config, deltas)

	log.Printf("Experiment %d: scheduling %d bursts with freq %vs and %d gateways (bursts/gateways*freq=%v), estimated to complete on %v",
		config.Id, config.Bursts, config.FrequencySeconds, len(config.GatewayEndpoints),
		float64(config.Bursts)/float64(len(config.GatewayEndpoints))*config.FrequencySeconds,
		time.Now().Add(estimateTime).UTC().Format(time.RFC3339))

	burstId := 0
	deltaIndex := 0
	// Schedule all bursts for this experiment
	for burstId < config.Bursts {
		time.Sleep(deltas[deltaIndex])

		// Send one burst to each available gateway (the more gateways used, the faster the experiment)
		for gatewayId := 0; gatewayId < len(config.GatewayEndpoints) && burstId < config.Bursts; gatewayId++ {
			// Every refresh period, we cycle through burst sizes if they're dynamic i.e. more than 1 element
			burstSize, _ := strconv.Atoi(config.BurstSizes[deltaIndex%len(config.BurstSizes)])
			serviceLoad, _ := strconv.Atoi(config.LambdaIncrementLimit[deltaIndex%len(config.LambdaIncrementLimit)])
			burst(config, burstId, burstSize, config.GatewayEndpoints[gatewayId], serviceLoad, safeExperimentWriter)
			burstId++
		}

		// After all gateways have been used for bursts, sleep again
		deltaIndex++
	}
}

func estimateTotalDuration(config configuration.ExperimentConfig, deltas []time.Duration) time.Duration {
	estimateTime := deltas[0]
	for _, burstDelta := range deltas[1 : config.Bursts/len(config.GatewayEndpoints)] {
		estimateTime += burstDelta
	}
	return estimateTime
}

func burst(config configuration.ExperimentConfig, burstId int, requests int, gatewayEndpointID string, serviceLoad int,
	safeExperimentWriter *SafeWriter) {
	gatewayEndpointURL := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod", gatewayEndpointID, region)
	log.Printf("Experiment %d: starting burst %d, making %d requests with service load %v to API Gateway (%s).",
		config.Id,
		burstId,
		requests,
		serviceLoad,
		gatewayEndpointURL,
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go safeExperimentWriter.GenerateLatencyRecord(gatewayEndpointURL, &requestsWaitGroup, serviceLoad,
			config.PayloadLengthBytes, burstId)
	}
	requestsWaitGroup.Wait()
	log.Printf("Experiment %d: received all responses for burst %d.", config.Id, burstId)
}
