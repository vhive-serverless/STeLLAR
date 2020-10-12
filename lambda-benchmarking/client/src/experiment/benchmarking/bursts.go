package benchmarking

import (
	"fmt"
	"lambda-benchmarking/client/experiment/configuration"
	"log"
	"strconv"
	"sync"
	"time"
)

// TriggerRelativeAsyncBursts
func RunProfiler(config configuration.ExperimentConfig, deltas []time.Duration, safeExperimentWriter *SafeWriter) {
	estimateTime := deltas[0]
	for _, burstDelta := range deltas[1 : config.Bursts/len(config.GatewayEndpoints)] {
		estimateTime += burstDelta
	}

	log.Printf("Experiment %d: scheduling %d bursts, estimated to complete on %v", config.Id, config.Bursts,
		time.Now().Add(estimateTime).Format(time.RFC3339))

	var burstsWaitGroup sync.WaitGroup
	burstId := 0
	deltaIndex := 0
	// Schedule all bursts for this experiment
	for burstId < config.Bursts {
		// Send one burst to each available gateway (the more gateways used, the faster the experiment)
		for gatewayId := 0; gatewayId < len(config.GatewayEndpoints) && burstId < config.Bursts; gatewayId++ {
			// Every refresh period, we cycle through burst sizes if they're dynamic i.e. more than 1 element
			burstSize, _ := strconv.Atoi(config.BurstSizes[deltaIndex%len(config.BurstSizes)])
			burstsWaitGroup.Add(1)
			go burst(&burstsWaitGroup, config, burstId, burstSize, config.GatewayEndpoints[gatewayId], safeExperimentWriter)
			burstId++
		}

		// After all gateways have been used for bursts, a new refresh period starts right away irrespective of burst statuses
		time.Sleep(deltas[deltaIndex])
		burstsWaitGroup.Wait()
		deltaIndex++
	}
}

func burst(burstsWaitGroup *sync.WaitGroup, config configuration.ExperimentConfig, burstId int, requests int, gatewayEndpointID string,
	safeExperimentWriter *SafeWriter) {
	gatewayEndpointURL := fmt.Sprintf("https://%s.execute-api.eu-west-2.amazonaws.com/prod", gatewayEndpointID)
	defer burstsWaitGroup.Done()
	log.Printf("Experiment %d: starting burst %d: making %d requests to API Gateway (%s).",
		config.Id,
		burstId,
		requests,
		gatewayEndpointURL,
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go safeExperimentWriter.GenerateLatencyRecord(gatewayEndpointURL, &requestsWaitGroup, config.LambdaIncrementLimit,
			config.PayloadLengthBytes, burstId)
	}
	requestsWaitGroup.Wait()
	log.Printf("Experiment %d: received all responses for burst %d.", config.Id, burstId)
}
