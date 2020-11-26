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

package connection

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/setup/functions/connection/amazon"
	"lambda-benchmarking/client/setup/functions/util"
	"strings"
)

//Endpoint is the schema for the configuration of provider endpoints.
type Endpoint struct {
	GatewayID        string  `json:"GatewayID"`
	FunctionMemoryMB int64   `json:"FunctionMemoryMB"`
	ImageSizeMB      float64 `json:"ImageSizeMB"`
}

//ServerlessInterface creates an interface through which to interact with various providers
type ServerlessInterface struct {
	//ListAPIs will list all endpoints corresponding to all serverless functions.
	ListAPIs func() []Endpoint

	//DeployFunction will create a new serverless function in the specified language, with id `i`. An API for it will
	//then be created, as well as corresponding interactions between them and specific permissions.
	DeployFunction func(language string, memoryAssigned int) string

	//RemoveFunction will remove the serverless function with given ID.
	RemoveFunction func(uniqueID string)

	//UpdateFunction will update the source code of the serverless function with given ID.
	UpdateFunction func(uniqueID string, memoryAssigned int)
}

//Singleton allows the client to interact with various serverless actions
var Singleton *ServerlessInterface

//Initialize will create a new provider connection to interact with
func Initialize(provider string) {
	switch provider {
	case "aws":
		setupAWSConnection()
	default:
		log.Warnf("Provider %s might not support initialization with the client.", provider)
	}
}

func setupAWSConnection() {
	amazon.InitializeSingleton()

	Singleton = &ServerlessInterface{
		ListAPIs: func() []Endpoint {
			result := amazon.AWSSingleton.ListFunctions()

			functions := make([]Endpoint, 0)
			for _, function := range result {
				functionGatewayID := strings.Split(*function.FunctionName, "_")[1]
				functions = append(functions, Endpoint{
					GatewayID:        functionGatewayID,
					FunctionMemoryMB: *function.MemorySize,
					ImageSizeMB:      util.BytesToMB(*function.CodeSize),
				})
			}

			return functions
		},
		DeployFunction: func(language string, memoryAssigned int) string {
			return amazon.AWSSingleton.DeployFunction(language, int64(memoryAssigned))
		},
		RemoveFunction: func(uniqueID string) {
			amazon.AWSSingleton.RemoveFunction(uniqueID)
			amazon.AWSSingleton.RemoveAPI(uniqueID)
		},
		UpdateFunction: func(uniqueID string, memoryAssigned int) {
			amazon.AWSSingleton.UpdateFunction(uniqueID)
			amazon.AWSSingleton.UpdateFunctionConfiguration(uniqueID, int64(memoryAssigned))
		},
	}
}
