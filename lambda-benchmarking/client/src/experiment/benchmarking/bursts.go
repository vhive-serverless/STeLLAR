package benchmarking

import (
	"lambda-benchmarking/client/experiment/configuration"
	"log"
	"strconv"
	"sync"
	"time"
)

// TriggerRelativeAsyncBursts
func RunProfiler(config configuration.ExperimentConfig, burstRelativeDeltas []time.Duration, safeExperimentWriter *SafeWriter) {
	var burstsWaitGroup sync.WaitGroup

	burstId := 0
	deltaIndex := 0
	// Schedule all bursts for this experiment
	for burstId < config.Bursts {
		relativeSleepTime := burstRelativeDeltas[deltaIndex]

		// Send one burst to each available gateway (the more gateways used, the faster the experiment)
		for gatewayId := 0; gatewayId < len(config.GatewayEndpoints) && burstId < config.Bursts; gatewayId++ {
			burstsWaitGroup.Add(1)

			// Every refresh period, we cycle through burst sizes if they're dynamic i.e. more than 1 element
			burstSize, _ := strconv.Atoi(config.BurstSizes[deltaIndex%len(config.BurstSizes)])
			go burst(&burstsWaitGroup, config, burstId, relativeSleepTime, burstSize, config.GatewayEndpoints[gatewayId],
				safeExperimentWriter)
			burstId++
		}

		// After all gateways have been used for bursts, a new refresh period starts
		deltaIndex++
	}

	log.Printf("Experiment %d: scheduled %d bursts, estimated to complete on %v", config.Id, config.Bursts,
		time.Now().Add(burstRelativeDeltas[deltaIndex-1]).Format(time.RFC3339))
	burstsWaitGroup.Wait()
}

func burst(burstsWaitGroup *sync.WaitGroup, config configuration.ExperimentConfig, burstId int, relativeDelta time.Duration,
	requests int, gatewayEndpoint string, safeExperimentWriter *SafeWriter) {
	defer burstsWaitGroup.Done()
	time.Sleep(relativeDelta)

	log.Printf("Experiment %d: starting burst %d (%v): making %d requests to API Gateway (%s).",
		config.Id,
		burstId,
		relativeDelta,
		requests,
		gatewayEndpoint,
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go safeExperimentWriter.GenerateLatencyRecord(gatewayEndpoint, &requestsWaitGroup, config.LambdaIncrementLimit,
			config.PayloadLengthBytes, burstId)
	}
	requestsWaitGroup.Wait()
	log.Printf("Experiment %d: received all responses for burst %d.", config.Id, burstId)
}
