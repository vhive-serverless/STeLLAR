package benchmarking

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/configuration"
	"strconv"
	"sync"
	"time"
)

// TriggerRelativeAsyncBursts
func RunProfiler(config configuration.ExperimentConfig, deltas []time.Duration, safeExperimentWriter *SafeWriter) {
	estimateTime := estimateTotalDuration(config, deltas)

	log.Infof("Experiment %d: scheduling %d bursts with freq %vs and %d gateways (bursts/gateways*freq=%v), estimated to complete on %v",
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
			serviceLoad := config.FunctionIncrementLimits[deltaIndex%len(config.FunctionIncrementLimits)]
			sendBurst(config, burstId, burstSize, config.GatewayEndpoints[gatewayId], serviceLoad, safeExperimentWriter)
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

func sendBurst(config configuration.ExperimentConfig, burstId int, requests int, gatewayEndpointID string,
	functionIncrementLimit int64, safeExperimentWriter *SafeWriter) {
	gatewayEndpointURL := fmt.Sprintf("https://%s.execute-api.us-west-1.amazonaws.com/prod", gatewayEndpointID)
	log.Infof("Experiment %d: starting burst %d, making %d requests with service load %dms to API Gateway (%s).",
		config.Id,
		burstId,
		requests,
		functionIncrementLimit,
		gatewayEndpointURL,
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go safeExperimentWriter.GenerateLatencyRecord(gatewayEndpointURL, &requestsWaitGroup, functionIncrementLimit,
			config.PayloadLengthBytes, burstId)
	}
	requestsWaitGroup.Wait()
	log.Infof("Experiment %d: received all responses for burst %d.", config.Id, burstId)
}
