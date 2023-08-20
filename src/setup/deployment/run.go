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
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"stellar/setup/deployment/packaging"
	"stellar/util"
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
				util.BytesToMB(zippedBinaryFileSizeBytes),
			)
			deploymentSizeBytes = zippedBinaryFileSizeBytes
		}

		if deploymentSizeBytes < zippedBinaryFileSizeBytes {
			log.Fatalf("[sub-experiment %d] Total size (~%vMB) cannot be smaller than zipped binary size (~%vMB).",
				experimentID,
				util.BytesToMB(deploymentSizeBytes),
				util.BytesToMB(zippedBinaryFileSizeBytes),
			)
		}

		packaging.GenerateFillerFile(experimentID, fillerFilePath, deploymentSizeBytes-zippedBinaryFileSizeBytes)
		zipPath := packaging.GenerateZIP(experimentID, fillerFilePath, binaryPath, "benchmarking.zip")
		packaging.SetupZIPDeployment(provider, deploymentSizeBytes, zipPath)

		return util.BytesToMB(deploymentSizeBytes), handlerPath
	case "Image":
		log.Warn("Container image deployment does not support code size verification on AWS, making the image size benchmarks unreliable.")

		// TODO: Size of containerized binary should be subtracted, seems to be 134MB in Amazon ECR...
		packaging.GenerateFillerFile(experimentID, fillerFilePath, int64(math.Max(float64(deploymentSizeBytes)-134, 0)))
		packaging.SetupContainerImageDeployment(function, provider, rawCodePath)

	default:
		log.Fatalf("[sub-experiment %d] Unrecognized package type: %s", experimentID, packageType)
	}

	return util.BytesToMB(deploymentSizeBytes), ""
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
