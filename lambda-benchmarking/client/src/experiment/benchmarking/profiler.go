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

package benchmarking

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/configuration"
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
