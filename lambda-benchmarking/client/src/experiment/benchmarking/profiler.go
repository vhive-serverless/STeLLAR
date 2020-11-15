package benchmarking

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/configuration"
	"time"
)

//RunProfiler will trigger bursts sequentially to each available gateway for a given experiment, then sleep for the
//selected interval and start the process all over again.
func RunProfiler(config configuration.SubExperiment, deltas []time.Duration, safeExperimentWriter *SafeWriter) {
	log.Infof("SubExperiment %d: running profiler, scheduling %d bursts with freq ~%vs and %d gateways (bursts/gateways*freq=%v)",
		config.ID, config.Bursts, config.CooldownSeconds, len(config.GatewayEndpoints),
		float64(config.Bursts)/float64(len(config.GatewayEndpoints))*config.CooldownSeconds)

	burstID := 0
	deltaIndex := 0
	for burstID < config.Bursts {
		time.Sleep(deltas[deltaIndex])

		// Send one burst to each available gateway (the more gateways used, the faster the experiment)
		for gatewayID := 0; gatewayID < len(config.GatewayEndpoints) && burstID < config.Bursts; gatewayID++ {
			// Every refresh period, we cycle through burst sizes if they're dynamic i.e. more than 1 element
			serviceLoad := config.FunctionIncrementLimits[min(deltaIndex, len(config.FunctionIncrementLimits)-1)]
			burstSize := config.BurstSizes[min(deltaIndex, len(config.BurstSizes)-1)]
			sendBurst(config, burstID, burstSize, config.GatewayEndpoints[gatewayID], serviceLoad, safeExperimentWriter)
			burstID++
		}

		deltaIndex++
		log.Debugf("SubExperiment %d: all %d gateways have been used for bursts, flushing and sleeping for %v...", config.ID, len(config.GatewayEndpoints), deltas[deltaIndex-1])
		safeExperimentWriter.Writer.Flush()
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
