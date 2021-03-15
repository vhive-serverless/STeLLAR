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

package setup

import (
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"time"
)

var cachedServiceTimeIncrement map[string]int64

// findBusySpinIncrements transforms given service times (e.g., 10s) into busy-spin increments (e.g., 10,000,000)
func findBusySpinIncrements(config *Configuration) {
	standardIncrement := int64(1e10)
	standardDurationMs := timeSession(standardIncrement).Milliseconds()
	cachedServiceTimeIncrement = make(map[string]int64)

	for subExperimentIndex := range config.SubExperiments {
		findBusySpinIncrement(&config.SubExperiments[subExperimentIndex], standardIncrement, standardDurationMs)
	}
}

func findBusySpinIncrement(subExperiment *SubExperiment, standardIncrement int64, standardDurationMs int64) {
	for _, serviceTime := range subExperiment.DesiredServiceTimes {
		if cachedIncrement, ok := cachedServiceTimeIncrement[serviceTime]; ok {
			log.Debugf("Using cached increment %d for desired service time %v", cachedIncrement, serviceTime)
			subExperiment.BusySpinIncrements = append(subExperiment.BusySpinIncrements, cachedIncrement)
			continue
		}

		parsedDesiredDuration, err := time.ParseDuration(serviceTime)
		if err != nil {
			log.Fatalf("Could not parse desired function run duration %s from configuration file.", serviceTime)
		}

		desiredDurationMs := parsedDesiredDuration.Milliseconds()
		log.Infof("Determining function increment for a duration of %dms...", desiredDurationMs)

		ratio := big.NewRat(desiredDurationMs, standardDurationMs)
		currentIncrement := big.NewRat(standardIncrement, 1)
		currentIncrement.Mul(currentIncrement, ratio)

		suggestedIncrementFloat, _ := currentIncrement.Float64()
		suggestedIncrement := int64(suggestedIncrementFloat)
		suggestedDurationMs := timeSession(suggestedIncrement).Milliseconds()
		if math.Abs(float64(suggestedDurationMs)-float64(desiredDurationMs)) > 0.05 {
			log.Warnf("Suggested increment %d (duration %dms) is not within 5%% of desired duration %dms",
				suggestedIncrement, suggestedDurationMs, desiredDurationMs)

			promptedIncrement := promptForNumber("Please enter a better increment (leave empty for unchanged): ")
			if promptedIncrement != nil {
				suggestedIncrement = *promptedIncrement
			}
		}

		log.Infof("Using increment %d (timed ~%dms) for desired %dms", suggestedIncrement, suggestedDurationMs, desiredDurationMs)
		cachedServiceTimeIncrement[serviceTime] = suggestedIncrement
		subExperiment.BusySpinIncrements = append(subExperiment.BusySpinIncrements, suggestedIncrement)
	}
}

func timeSession(increment int64) time.Duration {
	start := time.Now()
	for i := int64(0); i < increment; i++ {
	}
	return time.Since(start)
}
