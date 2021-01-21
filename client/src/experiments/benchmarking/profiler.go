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
	"os"
	"time"
	"vhive-bench/client/setup"
)

//RunProfiler will trigger bursts sequentially to each available gateway for a given experiment, then sleep for the
//selected interval and start the process all over again.
func RunProfiler(provider string, experiment setup.SubExperiment, deltas []time.Duration,
	latenciesFile *os.File, dataTransfersFile *os.File) {
	log.Infof("[sub-experiment %d] Running profiler, scheduling %d bursts with freq ~%vs and %d gateways (bursts/gateways*freq=%v)",
		experiment.ID, experiment.Bursts, experiment.IATSeconds, len(experiment.GatewayEndpoints),
		float64(experiment.Bursts)/float64(len(experiment.GatewayEndpoints))*experiment.IATSeconds)

	latenciesWriter := NewLatenciesWriter(latenciesFile)
	dataTransferWriter := NewDataTransferWriter(dataTransfersFile, experiment.DataTransferChainLength)

	burstID := 0
	deltaIndex := 0
	for burstID < experiment.Bursts {
		time.Sleep(deltas[deltaIndex])

		// Send one burst to each available gateway (the more gateways used, the faster the experiment)
		for gatewayID := 0; gatewayID < len(experiment.GatewayEndpoints) && burstID < experiment.Bursts; gatewayID++ {
			// Every refresh period, we cycle through burst sizes if they're dynamic i.e. more than 1 element
			serviceLoad := experiment.FunctionIncrementLimits[min(deltaIndex, len(experiment.FunctionIncrementLimits)-1)]
			burstSize := experiment.BurstSizes[min(deltaIndex, len(experiment.BurstSizes)-1)]
			sendBurst(provider, experiment, burstID, burstSize, experiment.GatewayEndpoints[gatewayID], serviceLoad, latenciesWriter, dataTransferWriter)
			burstID++
		}

		deltaIndex++
		log.Debugf("[sub-experiment %d] All %d gateways have been used for bursts, flushing and sleeping for %v...", experiment.ID, len(experiment.GatewayEndpoints), deltas[deltaIndex-1])
		latenciesWriter.Writer.Flush()
		if dataTransferWriter != nil {
			dataTransferWriter.Writer.Flush()
		}
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
