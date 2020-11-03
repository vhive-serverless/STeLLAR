package benchmarking

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/configuration"
	"strconv"
	"time"
)

func RunProfiler(config configuration.ExperimentConfig, deltas []time.Duration, safeExperimentWriter *SafeWriter) {
	log.Infof("Experiment %d: running profiler, scheduling %d bursts with freq ~%vs and %d gateways (bursts/gateways*freq=%v), estimated to complete on %v",
		config.Id, config.Bursts, config.FrequencySeconds, len(config.GatewayEndpoints),
		float64(config.Bursts)/float64(len(config.GatewayEndpoints))*config.FrequencySeconds,
		time.Now().Add(estimateTotalDuration(config, deltas)).UTC().Format(time.RFC3339))

	burstId := 0
	deltaIndex := 0
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

		deltaIndex++
		log.Debugf("Experiment %d: all %d gateways have been used for bursts, flushing and sleeping for %v...", config.Id, len(config.GatewayEndpoints), deltas[deltaIndex-1])
		safeExperimentWriter.Writer.Flush()
	}
}

func estimateTotalDuration(config configuration.ExperimentConfig, deltas []time.Duration) time.Duration {
	log.Debugf("Experiment %d: estimating total duration with deltas %v", config.Id, deltas)
	estimateTime := deltas[0]
	for _, burstDelta := range deltas[1 : config.Bursts/len(config.GatewayEndpoints)] {
		estimateTime += burstDelta
	}
	return estimateTime
}
