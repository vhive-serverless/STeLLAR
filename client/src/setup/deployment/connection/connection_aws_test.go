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
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
	"vhive-bench/client/setup/deployment"
	"vhive-bench/client/util"
)

const awsAPIsLimitIncl = 600
const aws = "aws"
const golang = "go1.x"

func TestAWSListAPIs(t *testing.T) {
	Initialize("aws", "")
	apis := Singleton.ListAPIs()

	require.True(t, 0 <= len(apis) && len(apis) <= awsAPIsLimitIncl)
}

func TestAWSRemoveFunction(t *testing.T) {
	Initialize("aws", "")
	apis := Singleton.ListAPIs()

	var removedAPIID string
	if len(apis) == 0 {
		removedAPIID, _, _ = deployRandomMemoryFunction()
	} else {
		removedAPIID = apis[0].GatewayID
	}

	Singleton.RemoveFunction(removedAPIID)

	// Check that removing succeeded
	apis = Singleton.ListAPIs()
	for _, api := range apis {
		if api.GatewayID == removedAPIID {
			require.FailNow(t, "Lambda function was in fact not removed: function still listed by AWS.")
		}
	}
}

func TestAWSDeployFunction(t *testing.T) {
	Initialize("aws", "")
	apis := Singleton.ListAPIs()

	if len(apis) >= awsAPIsLimitIncl {
		Singleton.RemoveFunction(apis[0].GatewayID)
	}

	deployedFunctionID, deployedImageSizeMB, desiredFunctionMemoryMB := deployRandomMemoryFunction()

	// Check that deployment succeeded
	apis = Singleton.ListAPIs()
	foundDeployedFunction := false
	for _, api := range apis {
		if api.GatewayID == deployedFunctionID &&
			int(api.ImageSizeMB) == int(deployedImageSizeMB) &&
			api.FunctionMemoryMB == desiredFunctionMemoryMB {
			foundDeployedFunction = true
		}
	}
	require.True(t, foundDeployedFunction)

	// Cleanup
	Singleton.RemoveFunction(deployedFunctionID)
}

func TestAWSUpdateFunction(t *testing.T) {
	Initialize("aws", "")
	apis := Singleton.ListAPIs()

	var repurposedAPIID string
	if len(apis) == 0 {
		repurposedAPIID, _, _ = deployRandomMemoryFunction()
	} else {
		repurposedAPIID = apis[0].GatewayID
		setupDeployment()
	}

	repurposedFunctionMemory := rand.Intn(1000-128) + 128
	Singleton.UpdateFunction(repurposedAPIID, int64(repurposedFunctionMemory))

	// Check that repurposing succeeded
	apis = Singleton.ListAPIs()
	foundRepurposedFunction := false
	for _, api := range apis {
		if api.GatewayID == repurposedAPIID &&
			api.FunctionMemoryMB == int64(repurposedFunctionMemory) {
			foundRepurposedFunction = true
		}
	}
	require.True(t, foundRepurposedFunction)

	// Cleanup
	Singleton.RemoveFunction(repurposedAPIID)
}

func deployRandomMemoryFunction() (string, float64, int64) {
	rand.Seed(time.Now().Unix())
	desiredFunctionMemoryMB := int64(rand.Intn(1000-128) + 128)

	deployedImageSizeMB := setupDeployment()

	return Singleton.DeployFunction(golang, desiredFunctionMemoryMB), deployedImageSizeMB, desiredFunctionMemoryMB
}

func setupDeployment() float64 {
	deployedImageSizeMB := deployment.SetupDeployment(
		fmt.Sprintf("../raw-code/%s/%s-handler.go", golang, aws),
		aws,
		golang,
		util.MBToBytes(60.),
	)
	return deployedImageSizeMB
}
