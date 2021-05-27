// MIT License
//
// Copyright (c) 2021 Theodor Amariucai
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
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"vhive-bench/setup/deployment/packaging"
	"vhive-bench/util"
)

//SetupDeployment will create the serverless function zip deployment for the given provider,
//in the given language and of the given size in bytes. Returns size of deployment in MB.
func SetupDeployment(rawCodePath string, provider string, deploymentSizeBytes int64, packageType string, experimentID int, function string) (float64, string) {
	fillerFilePath := rawCodePath + "/filler.file"

	switch packageType {
	case "Zip":
		_, binaryPath := getBinaryInfo(rawCodePath, experimentID)

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

		generateFillerFile(experimentID, fillerFilePath, deploymentSizeBytes-zippedBinaryFileSizeBytes)
		zipPath := packaging.GenerateZIP(experimentID, fillerFilePath, binaryPath)
		packaging.SetupZIPDeployment(provider, deploymentSizeBytes, zipPath)

		return util.BytesToMB(deploymentSizeBytes), binaryPath
	case "Image":
		log.Warn("Container image deployment does not support code size verification on AWS, making the image size benchmarks unreliable.")

		// TODO: Size of containerized binary should be subtracted, seems to be 134MB in Amazon ECR...
		generateFillerFile(experimentID, fillerFilePath, int64(math.Max(float64(deploymentSizeBytes)-134, 0)))
		packaging.SetupContainerImageDeployment(function, provider, rawCodePath)

	default:
		log.Fatalf("[sub-experiment %d] Unrecognized package type: %s", experimentID, packageType)
	}

	return util.BytesToMB(deploymentSizeBytes), ""
}

func generateFillerFile(experimentID int, fillerFilePath string, sizeBytes int64) {
	log.Infof("[sub-experiment %d] Generating filler file to be included in deployment...", experimentID)

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("[sub-experiment %d] Failed to fill buffer with random bytes: `%s`", experimentID, err.Error())
	}

	if err := ioutil.WriteFile(fillerFilePath, buffer, 0666); err != nil {
		log.Fatalf("[sub-experiment %d] Could not generate random file with size %d bytes: %v", experimentID, sizeBytes, err)
	}

	log.Infof("[sub-experiment %d] Successfully generated the filler file.", experimentID)
}

func getBinaryInfo(rawCodePath string, experimentID int) (int64, string) {
	binaryPath := fmt.Sprintf("%s/%s", rawCodePath, "handler")

	log.Infof("[sub-experiment %d] Getting binary file size for the function(s) to be deployed, path is %q...", experimentID, binaryPath)

	fi, err := os.Stat(binaryPath)
	if err != nil {
		log.Fatalf("[sub-experiment %d] Could not get size of binary file: %s", experimentID, err.Error())
	}

	log.Infof("[sub-experiment %d] Successfully retrieved binary file size (%d bytes) for deployment.", experimentID, fi.Size())
	return fi.Size(), binaryPath
}
