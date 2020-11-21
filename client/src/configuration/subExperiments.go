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

package configuration

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"lambda-benchmarking/client/prompts"
	"math"
	"math/big"
	"os"
	"time"
)

//Configuration is the schema for all experiment configurations.
type Configuration struct {
	Sequential     bool            `json:"Sequential"`
	SubExperiments []SubExperiment `json:"SubExperiments"`
}

//SubExperiment is the schema for sub-experiment configurations.
type SubExperiment struct {
	Title                   string   `json:"Title"`
	Bursts                  int      `json:"Bursts"`
	BurstSizes              []int    `json:"BurstSizes"`
	PayloadLengthBytes      int      `json:"PayloadLengthBytes"`
	CooldownSeconds         float64  `json:"CooldownSeconds"`
	FunctionIncrementLimits []int64  `json:"FunctionIncrementLimits"`
	DesiredServiceTimes     []string `json:"DesiredServiceTimes"`
	IATType                 string   `json:"IATType"`
	Provider                string   `json:"Provider"`
	GatewaysNumber          int      `json:"GatewaysNumber"`
	Visualization           string   `json:"Visualization"`
	FunctionMemoryMB        int64    `json:"FunctionMemoryMB"`
	FunctionImageSizeMB     int64    `json:"FunctionImageSizeMB"`
	GatewayEndpoints        []string
	ID                      int
}

//extractSubExperiments will read the given JSON configuration file and load it as an array of sub-experiment configurations.
func extractSubExperiments(configFile *os.File) Configuration {
	configByteValue, _ := ioutil.ReadAll(configFile)

	var parsedConfiguration Configuration
	if err := json.Unmarshal(configByteValue, &parsedConfiguration); err != nil {
		log.Fatalf("Could not extract experiment configuration from file: %s", err.Error())
	}

	standardIncrement := int64(1e10)
	standardDurationMs := timeSession(standardIncrement).Milliseconds()
	cachedServiceTimeIncrement = make(map[string]int64)
	for subExperimentIndex := range parsedConfiguration.SubExperiments {
		determineFunctionIncrementLimits(&parsedConfiguration.SubExperiments[subExperimentIndex],
			standardIncrement, standardDurationMs)
	}

	setDefaults(parsedConfiguration.SubExperiments)
	return parsedConfiguration
}

const defaultVisualization = "cdf"
const defaultIATType = "stochastic"
const defaultProvider = "aws"
const defaultFunctionMemoryMB = 128
const defaultGatewaysNumber = 1

func setDefaults(parsedSubExps []SubExperiment) {
	for index := range parsedSubExps {
		if parsedSubExps[index].Visualization == "" {
			parsedSubExps[index].Visualization = defaultVisualization
		}
		if parsedSubExps[index].IATType == "" {
			parsedSubExps[index].IATType = defaultIATType
		}
		if parsedSubExps[index].Provider == "" {
			parsedSubExps[index].Provider = defaultProvider
		}
		if parsedSubExps[index].FunctionMemoryMB == 0 {
			parsedSubExps[index].FunctionMemoryMB = defaultFunctionMemoryMB
		}
		if parsedSubExps[index].GatewaysNumber == 0 {
			parsedSubExps[index].GatewaysNumber = defaultGatewaysNumber
		}
	}
}

var cachedServiceTimeIncrement map[string]int64

func determineFunctionIncrementLimits(subExperiment *SubExperiment, standardIncrement int64, standardDurationMs int64) {
	for _, serviceTime := range subExperiment.DesiredServiceTimes {
		if cachedIncrement, ok := cachedServiceTimeIncrement[serviceTime]; ok {
			log.Infof("Using cached increment %d for desired %v", cachedIncrement, serviceTime)
			subExperiment.FunctionIncrementLimits = append(subExperiment.FunctionIncrementLimits, cachedIncrement)
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
		if !almostEqual(float64(suggestedDurationMs), float64(desiredDurationMs), float64(desiredDurationMs)*0.05) {
			log.Warnf("Suggested increment %d (duration %dms) is not within 5%% of desired duration %dms",
				suggestedIncrement, suggestedDurationMs, desiredDurationMs)

			promptedIncrement := prompts.PromptForNumber("Please enter a better increment (leave empty for unchanged): ")
			if promptedIncrement != nil {
				suggestedIncrement = *promptedIncrement
			}
		}

		log.Infof("Using increment %d (timed ~%dms) for desired %dms", suggestedIncrement, suggestedDurationMs, desiredDurationMs)
		cachedServiceTimeIncrement[serviceTime] = suggestedIncrement
		subExperiment.FunctionIncrementLimits = append(subExperiment.FunctionIncrementLimits, suggestedIncrement)
	}
}

func timeSession(increment int64) time.Duration {
	start := time.Now()
	for i := int64(0); i < increment; i++ {
	}
	return time.Since(start)
}

func almostEqual(a, b float64, float64EqualityThreshold float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
