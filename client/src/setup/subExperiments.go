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
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"lambda-benchmarking/client/prompts"
	"lambda-benchmarking/client/setup/functions/connection"
	"lambda-benchmarking/client/setup/functions/generator"
	"lambda-benchmarking/client/setup/functions/util"
	"math"
	"math/big"
	"os"
	"time"
)

//Configuration is the schema for all experiment configurations.
type Configuration struct {
	Sequential     bool            `json:"Sequential"`
	Provider       string          `json:"Provider"`
	Runtime        string          `json:"Runtime"`
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
	GatewaysNumber          int      `json:"GatewaysNumber"`
	Visualization           string   `json:"Visualization"`
	FunctionMemoryMB        int64    `json:"FunctionMemoryMB"`
	FunctionImageSizeMB     int64    `json:"FunctionImageSizeMB"`
	GatewayEndpoints        []string
	ID                      int
}

func initializeSubExperiment(config Configuration, index int, availableEndpoints []connection.Endpoint) []connection.Endpoint {
	config.SubExperiments[index].ID = index

	for _, burstSize := range config.SubExperiments[index].BurstSizes {
		if burstSize > manyRequestsInBurstWarnThreshold {
			log.Warnf("Experiment %d has a burst of size %d, NIC (Network Interface Controller) contention may occur.",
				index, burstSize)
			if !prompts.PromptForBool("Do you wish to continue?") {
				os.Exit(0)
			}
		}
	}

	chosenVisualization := config.SubExperiments[index].Visualization
	burstsNumber := config.SubExperiments[index].Bursts
	if burstsNumber >= manyFilesWarnThreshold && (chosenVisualization == "all" || chosenVisualization == "histogram") {
		log.Warnf("SubExperiment %d is generating histograms for each burst, this will create a large number (%d) of new files.",
			index, burstsNumber)
		if !prompts.PromptForBool("Do you wish to continue?") {
			os.Exit(0)
		}
	}

	if availableEndpoints == nil { // Provider string itself must be a hostname
		config.SubExperiments[index].GatewayEndpoints = []string{config.Provider}
	} else {
		return assignEndpoints(availableEndpoints, &config.SubExperiments[index], config.Provider, config.Runtime)
	}
	return nil
}

func assignEndpoints(availableEndpoints []connection.Endpoint, experiment *SubExperiment, provider string, runtime string) []connection.Endpoint {
	deploymentGeneratedForSubExperiment := false

	var assignedEndpoints []string
	for i := 0; i < experiment.GatewaysNumber; i++ {
		success := false

		for index, endpoint := range availableEndpoints {
			if endpoint.FunctionMemoryMB == experiment.FunctionMemoryMB &&
				(experiment.FunctionImageSizeMB == 0 || almostEqual(endpoint.ImageSizeMB, float64(experiment.FunctionImageSizeMB), 0.5)) {

				assignedEndpoints = append(assignedEndpoints, endpoint.GatewayID)
				availableEndpoints = removeEndpoint(availableEndpoints, index)
				success = true
				break
			}
		}

		if !success {
			log.Infof("Could not find a function to assign with memory %dMB and image size %dMB, deploying...",
				experiment.FunctionMemoryMB,
				experiment.FunctionImageSizeMB,
			)

			if !deploymentGeneratedForSubExperiment {
				generator.SetupDeployment(provider, runtime, util.MBToBytes(float64(experiment.FunctionImageSizeMB)))
				deploymentGeneratedForSubExperiment = true
			}
			assignedEndpoints = append(assignedEndpoints, connection.Singleton.DeployFunction(runtime, 128))

			//TODO: intelligently leverage connection.Singleton.RemoveFunction(uniqueID) &
			//TODO: connection.Singleton.UpdateFunction(uniqueID, 128)
			//TODO: once over 600 deployed functions
		}
	}

	log.Debugf("Assigning following endpoints to sub-experiment `%s`: %v", experiment.Title, assignedEndpoints)
	experiment.GatewayEndpoints = assignedEndpoints
	return availableEndpoints
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

	setDefaults(&parsedConfiguration)
	log.Debugf("Extracted %d sub-experiments from given configuration file.", len(parsedConfiguration.SubExperiments))
	return parsedConfiguration
}

const defaultVisualization = "cdf"
const defaultIATType = "stochastic"
const defaultProvider = "aws"
const defaultRuntime = "go1.x"
const defaultFunctionMemoryMB = 128
const defaultGatewaysNumber = 1

func setDefaults(parsedConfig *Configuration) {
	if parsedConfig.Provider == "" {
		parsedConfig.Provider = defaultProvider
	}
	if parsedConfig.Runtime == "" {
		parsedConfig.Runtime = defaultRuntime
	}

	for index := range parsedConfig.SubExperiments {
		if parsedConfig.SubExperiments[index].Visualization == "" {
			parsedConfig.SubExperiments[index].Visualization = defaultVisualization
		}
		if parsedConfig.SubExperiments[index].IATType == "" {
			parsedConfig.SubExperiments[index].IATType = defaultIATType
		}
		if parsedConfig.SubExperiments[index].FunctionMemoryMB == 0 {
			parsedConfig.SubExperiments[index].FunctionMemoryMB = defaultFunctionMemoryMB
		}
		if parsedConfig.SubExperiments[index].GatewaysNumber == 0 {
			parsedConfig.SubExperiments[index].GatewaysNumber = defaultGatewaysNumber
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
