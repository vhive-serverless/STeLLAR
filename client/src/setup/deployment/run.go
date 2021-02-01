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
	"os/exec"
	"strings"
	"vhive-bench/client/setup/deployment/packaging"
	"vhive-bench/client/util"
)

//SetupDeployment will create the serverless function zip deployment for the given provider,
//in the given language and of the given size in bytes. Returns size of deployment in MB.
func SetupDeployment(rawCodePath string, provider string, language string, deploymentSizeBytes int64, packageType string, experimentID int) (float64, string) {
	_, binaryPath := generateBinaryFile(rawCodePath, language, experimentID)
	fillerFilePath := strings.TrimSuffix(binaryPath, "/handler") + "/random.file"

	switch packageType {
	case "Zip":
		zippedBinaryFileSizeBytes := packaging.GetZippedBinaryFileSize(binaryPath)

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

		generateFillerFile(fillerFilePath, deploymentSizeBytes-zippedBinaryFileSizeBytes)
		zipPath := packaging.GenerateZIP(fillerFilePath, binaryPath)
		packaging.SetupZIPDeployment(provider, deploymentSizeBytes, zipPath)

	case "Image":
		log.Warn("Container image deployment does not support code size verification on AWS, making the image size benchmarks unreliable.")

		// TODO: Size of containerized binary should be subtracted, seems to be 134MB in Amazon ECR...
		generateFillerFile(fillerFilePath, int64(math.Max(float64(deploymentSizeBytes)-134, 0)))
		packaging.SetupContainerImageDeployment(provider, binaryPath)

	default:
		log.Fatalf("[sub-experiment %d] Unrecognized package type: %s", experimentID, packageType)
	}

	return util.BytesToMB(deploymentSizeBytes), binaryPath
}

func generateFillerFile(fillerFilePath string, sizeBytes int64) {
	log.Info("Generating filler file to be included in deployment...")

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("Failed to fill buffer with random bytes: `%s`", err.Error())
	}

	if err := ioutil.WriteFile(fillerFilePath, buffer, 0666); err != nil {
		log.Fatalf("Could not generate random file with size %d bytes", sizeBytes)
	}

	log.Info("Successfully generated the filler file.")
}

func generateBinaryFile(rawCodePath string, language string, experimentID int) (int64, string) {
	log.Infof("[sub-experiment %d] Building binary file for the function(s) to be deployed...", experimentID)

	if !util.FileExists(rawCodePath) {
		log.Fatalf("[sub-experiment %d] Code path %q does not exist, cannot deploy/update raw code.", experimentID, rawCodePath)
	}

	binaryPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(rawCodePath, "/main.go"), "handler")

	switch language {
	case "go1.x":
		cmd := exec.Command("go", "build", "-v", "-o", binaryPath, rawCodePath)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
		cmd.Env = append(cmd.Env, "GOOS=linux")
		util.RunCommandAndLog(cmd)
	default:
		log.Fatalf("[sub-experiment %d] Unrecognized language %s", experimentID, language)
	}

	fi, err := os.Stat(binaryPath)
	if err != nil {
		log.Fatalf("[sub-experiment %d] Could not get size of binary file: %s", experimentID, err.Error())
	}

	log.Infof("[sub-experiment %d] Successfully built the binary file for deployment.", experimentID)
	return fi.Size(), binaryPath
}
