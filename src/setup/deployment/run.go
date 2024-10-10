// MIT License
//
// Copyright (c) 2021 Theodor Amariucai and EASE Lab
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

package deployment

import (
	"fmt"
	"math"
	"os"
	"stellar/setup/deployment/packaging"
	"stellar/util"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"context"
	"stellar/setup/deployment/raw-code/serverless/azure"
)

// SetupDeployment will create the serverless function zip deployment for the given provider,
// in the given language and of the given size in bytes. Returns size of deployment in MB and the handler path for AWS automation.
func SetupDeployment(rawCodePath string, provider string, deploymentSizeBytes int64, packageType string, experimentID int, function string) (float64, string) {
	fillerFilePath := rawCodePath + "/filler.file"

	switch packageType {
	case "Zip":
		_, binaryPath, handlerPath := getExecutableInfo(rawCodePath, experimentID, function)

		zippedBinaryFileSizeBytes := packaging.GetZippedBinaryFileSize(experimentID, binaryPath)

		if deploymentSizeBytes == 0 {
			log.Infof("[sub-experiment %d] Desired image size is set to default (0MB), assigning size of zipped binary (%vMB)...",
				experimentID,
				util.BytesToMebibyte(zippedBinaryFileSizeBytes),
			)
			deploymentSizeBytes = zippedBinaryFileSizeBytes
		}

		if deploymentSizeBytes < zippedBinaryFileSizeBytes {
			log.Fatalf("[sub-experiment %d] Total size (~%vMB) cannot be smaller than zipped binary size (~%vMB).",
				experimentID,
				util.BytesToMebibyte(deploymentSizeBytes),
				util.BytesToMebibyte(zippedBinaryFileSizeBytes),
			)
		}

		packaging.GenerateFillerFile(experimentID, fillerFilePath, deploymentSizeBytes-zippedBinaryFileSizeBytes)
		zipPath := packaging.GenerateZIP(experimentID, fillerFilePath, binaryPath, "benchmarking.zip")
		packaging.SetupZIPDeployment(provider, deploymentSizeBytes, zipPath)

		return util.BytesToMebibyte(deploymentSizeBytes), handlerPath
	case "Image":
		log.Warn("Container image deployment does not support code size verification on AWS, making the image size benchmarks unreliable.")

		// TODO: Size of containerized binary should be subtracted, seems to be 134MB in Amazon ECR...
		packaging.GenerateFillerFile(experimentID, fillerFilePath, int64(math.Max(float64(deploymentSizeBytes)-134, 0)))
		//packaging.SetupContainerImageDeployment(function, provider, rawCodePath)

	default:
		log.Fatalf("[sub-experiment %d] Unrecognized package type: %s", experimentID, packageType)
	}

	return util.BytesToMebibyte(deploymentSizeBytes), ""
}

func getExecutableInfo(rawCodePath string, experimentID int, function string) (int64, string, string) {
	var binaryPath string
	var handlerPath string
	switch function {
	case "producer-consumer":
		binaryPath = fmt.Sprintf("%s/%s", rawCodePath, "handler")
		handlerPath = binaryPath
	case "hellopy":
		binaryPath = fmt.Sprintf("%s/%s", rawCodePath, "lambda_function.py")
		handlerPath = fmt.Sprintf("%s/%s", rawCodePath, "lambda_function.lambda_handler")
	default:
		log.Fatalf("[sub-experiment %d] Unrecognized or unimplemented function type for ZIP deployment: %s", experimentID, function)
	}

	log.Infof("[sub-experiment %d] Getting binary file size for the function(s) to be deployed, path is %q...", experimentID, binaryPath)

	fi, err := os.Stat(binaryPath)
	if err != nil {
		log.Fatalf("[sub-experiment %d] Could not get size of binary file: %s", experimentID, err.Error())
	}

	log.Infof("[sub-experiment %d] Successfully retrieved exec file size (%d bytes) for deployment.", experimentID, fi.Size())
	return fi.Size(), binaryPath, handlerPath
}

func RunDeployment() {
	// Step 1: Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	} else {
		log.Println(".env file loaded successfully.")
	}

	// Step 2: Load configuration into Config struct
	config := azure.LoadConfig()

	// Step 3: Validate required environment variables
	azure.ValidateConfig(config)

	// Step 4: Validate that required commands are available
	if !azure.IsCommandAvailable("az") {
		log.Fatal("'az' command is not available. Please install Azure CLI.")
	}

	if !azure.IsCommandAvailable("func") {
		log.Fatal("'func' command is not available. Please install Azure Functions Core Tools.")
	}

	// Step 5: Initialize Azure SDK credentials
	cred, err := azure.GetAzureCredential()
	if err != nil {
		log.Fatalf("Failed to obtain a credential: %v", err)
	}
	ctx := context.Background()

	// Step 6: Initialize Azure SDK clients
	err = azure.InitializeClients(ctx, cred, config)
	if err != nil {
		log.Fatalf("Failed to initialize Azure clients: %v", err)
	}

	// Step 7: Create Resource Group
	resourceGroup, err := azure.CreateResourceGroup(ctx, config)
	if err != nil {
		log.Fatalf("Failed to create resource group: %v", err)
	}
	log.Println("Resource Group Created:", *resourceGroup.ID)

	// Step 8: Check Storage Account Name Availability
	availability, err := azure.CheckNameAvailability(ctx, config)
	if err != nil {
		log.Fatalf("Failed to check storage account name availability: %v", err)
	}
	if !*availability.NameAvailable {
		log.Fatalf("Storage account name is not available: %s", *availability.Message)
	}

	// Step 9: Create Storage Account
	storageAccount, err := azure.CreateStorageAccount(ctx, config)
	if err != nil {
		log.Fatalf("Failed to create storage account: %v", err)
	}
	log.Println("Storage Account Created:", *storageAccount.ID)

	// Step 10: Get Storage Account Properties
	properties, err := azure.StorageAccountProperties(ctx, config)
	if err != nil {
		log.Fatalf("Failed to get storage account properties: %v", err)
	}
	log.Println("Storage Account Properties ID:", *properties.ID)

	// Step 11: Initialize Function App Project (if not already)
	err = azure.InitializeFunctionProject()
	if err != nil {
		log.Fatalf("Failed to initialize Function App project: %v", err)
	}
	log.Println("Function App Project Initialized Successfully.")

	// Step 12: Create New Function using `func new`
	err = azure.CreateNewFunction(config)
	if err != nil {
		log.Fatalf("Failed to create new Function: %v", err)
	}
	log.Println("New Function Created Successfully.")

	// Step 13: Execute Azure CLI Command to Create Function App
	err = azure.CreateFunctionApp(config)
	if err != nil {
		log.Fatalf("Failed to create Function App: %v", err)
	}
	log.Println("Function App Created Successfully.")

	// Step 14: Publish Function App
	err = azure.PublishFunctionApp(config)
	if err != nil {
		log.Fatalf("Failed to publish Function App: %v", err)
	}
	log.Println("Function App Published Successfully.")

	// Step 15: Cleanup Resources if KEEP_RESOURCE is not set
	if !azure.ShouldKeepResource(config.KeepResource) {
		err = azure.Cleanup(ctx, config)
		if err != nil {
			log.Fatalf("Failed to clean up resources: %v", err)
		}
		log.Println("Resources cleaned up successfully.")
	}
}
