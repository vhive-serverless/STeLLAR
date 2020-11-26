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
	"lambda-benchmarking/client/setup/functions/connection"
	"os"
)

const (
	manyRequestsInBurstWarnThreshold = 2000
	manyFilesWarnThreshold           = 500
)

//PrepareSubExperiments will read any required files, deploy functions etc. to get ready for the sub-experiments.
func PrepareSubExperiments(endpointsDirectoryPath string, configPath string) Configuration {
	configFile := readFile(configPath)

	config := extractSubExperiments(configFile)

	connection.Initialize(config.Provider)

	availableEndpoints := getAvailableEndpoints(endpointsDirectoryPath, config)

	for index := range config.SubExperiments {
		availableEndpoints = initializeSubExperiment(config, index, availableEndpoints)
	}

	return config
}

func extractProviderEndpoints(endpointsFile *os.File) []connection.Endpoint {
	configByteValue, _ := ioutil.ReadAll(endpointsFile)

	var parsedEndpoints []connection.Endpoint
	if err := json.Unmarshal(configByteValue, &parsedEndpoints); err != nil {
		log.Fatalf("Could not extract endpoints configuration from file: %s", err.Error())
	}

	return parsedEndpoints
}

func isStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func readFile(path string) *os.File {
	log.Debugf("Reading file for this run from `%s`", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not read file: %s", err.Error())
	}
	return file
}
