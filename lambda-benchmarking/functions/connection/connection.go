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
	"functions/connection/amazon"
	log "github.com/sirupsen/logrus"
)

//ServerlessInterface creates an interface through which to interact with various providers
type ServerlessInterface struct {
	//DeployFunction will create a new serverless function in the specified language, with id `i`. An API for it will
	//then be created, as well as corresponding interactions between them and specific permissions.
	DeployFunction func(id int, language string, memoryAssigned int)

	//RemoveFunction will remove the serverless function with id `i`.
	RemoveFunction func(id int)

	//UpdateFunction will update the source code of the serverless function with id `i`.
	UpdateFunction func(id int, memoryAssigned int)
}

//Singleton allows the client to interact with various serverless actions
var Singleton *ServerlessInterface

//Initialize will create a new provider connection to interact with
func Initialize(provider string) {
	switch provider {
	case "aws":
		setupAWSConnection()
	default:
		log.Fatalf("Unrecognized provider %s", provider)
	}
}

func setupAWSConnection() {
	amazon.InitializeSingleton()

	Singleton = &ServerlessInterface{
		DeployFunction: func(id int, language string, memoryAssigned int) {
			amazon.AWSSingleton.DeployFunction(id, language, int64(memoryAssigned))
		},
		RemoveFunction: func(id int) {
			amazon.AWSSingleton.RemoveFunction(id)
			amazon.AWSSingleton.RemoveAPI(id)
		},
		UpdateFunction: func(id int, memoryAssigned int) {
			amazon.AWSSingleton.UpdateFunction(id)
			amazon.AWSSingleton.UpdateFunctionConfiguration(id, int64(memoryAssigned))
		},
	}
}
