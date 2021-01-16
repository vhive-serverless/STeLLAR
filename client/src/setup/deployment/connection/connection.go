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

// Package connection provides support and abstraction (in the form of an interface)
// in communicating with external providers against which benchmarking is run.
package connection

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"strings"
	"vhive-bench/client/setup/deployment/connection/amazon"
	"vhive-bench/client/util"
)

//Endpoint is the schema for the configuration of provider endpoints.
type Endpoint struct {
	GatewayID        string  `json:"GatewayID"`
	FunctionMemoryMB int64   `json:"FunctionMemoryMB"`
	ImageSizeMB      float64 `json:"ImageSizeMB"`
	PackageType      string  `json:"PackageType"`
}

//ServerlessInterface creates an interface through which to interact with various providers
type ServerlessInterface struct {
	//ListAPIs will list all endpoints corresponding to all serverless functions.
	ListAPIs func() []Endpoint

	//DeployFunction will create a new serverless function in the specified language, with the specified amount of
	//memory. An API to access it will then be created, as well as corresponding permissions and integrations.
	DeployFunction func(packageType string, language string, memoryAssigned int64) string

	//RemoveFunction will remove the serverless function with given ID and its corresponding API.
	RemoveFunction func(uniqueID string)

	//UpdateFunction will update the source code of the serverless function with given ID to the specified
	//memory and to the most recently set code deployment settings (e.g., S3 key).
	UpdateFunction func(packageType string, uniqueID string, memoryAssigned int64)
}

//Singleton allows the client to interact with various serverless actions
var Singleton *ServerlessInterface

//Initialize will create a new provider connection to interact with
func Initialize(provider string, endpointsDirectoryPath string) {
	switch provider {
	case "aws":
		setupAWSConnection()
	case "vhive":
		setupFileConnection(path.Join(endpointsDirectoryPath, "vHive.json"))
	default:
		setupExternalConnection()
		log.Warnf("Provider %s does not support initialization with the client, setting to external URL.", provider)
	}
}

func setupAWSConnection() {
	amazon.InitializeSingleton()

	Singleton = &ServerlessInterface{
		ListAPIs: func() []Endpoint {
			result := amazon.AWSSingletonInstance.ListFunctions(nil)
			log.Infof("Found %d Lambda functions.", len(result))

			functions := make([]Endpoint, 0)
			for _, function := range result {
				functions = append(functions, Endpoint{
					GatewayID:        strings.Split(*function.FunctionName, "_")[1],
					FunctionMemoryMB: *function.MemorySize,
					ImageSizeMB:      util.BytesToMB(*function.CodeSize),
				})
			}

			return functions
		},
		DeployFunction: func(packageType string, language string, memoryAssigned int64) string {
			return amazon.AWSSingletonInstance.DeployFunction(packageType, language, memoryAssigned)
		},
		RemoveFunction: func(uniqueID string) {
			amazon.AWSSingletonInstance.RemoveFunction(uniqueID)
			amazon.AWSSingletonInstance.RemoveAPI(uniqueID)
		},
		UpdateFunction: func(packageType string, uniqueID string, memoryAssigned int64) {
			amazon.AWSSingletonInstance.UpdateFunction(packageType, uniqueID)
			amazon.AWSSingletonInstance.UpdateFunctionConfiguration(uniqueID, memoryAssigned)
		},
	}
}

func setupFileConnection(path string) {
	Singleton = &ServerlessInterface{
		ListAPIs: func() []Endpoint {
			endpointsFile := util.ReadFile(path)
			configByteValue, _ := ioutil.ReadAll(endpointsFile)

			var parsedEndpoints []Endpoint
			if err := json.Unmarshal(configByteValue, &parsedEndpoints); err != nil {
				log.Fatalf("Could not extract endpoints configuration from file: %s", err.Error())
			}

			return parsedEndpoints
		},
	}
}

func setupExternalConnection() {
	Singleton = &ServerlessInterface{
		ListAPIs: func() []Endpoint {
			return nil
		},
	}
}
