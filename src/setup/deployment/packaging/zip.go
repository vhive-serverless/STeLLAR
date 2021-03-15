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

package packaging

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"vhive-bench/setup/deployment/connection/amazon"
	"vhive-bench/util"
)

//SetupZIPDeployment will package the function using ZIP
func SetupZIPDeployment(provider string, deploymentSizeBytes int64, zipPath string) {
	deploymentSizeMB := util.BytesToMB(deploymentSizeBytes)
	switch provider {
	case "aws":
		if deploymentSizeMB > 50. {
			amazon.UploadZIPToS3(zipPath, deploymentSizeMB)
		} else {
			amazon.SetLocalZip(zipPath)
		}
	default:
		log.Warnf("Provider %s does not support ZIP deployment, skipping ZIP generation...", provider)
	}

	log.Debugf("Cleaning up ZIP %q...", zipPath)
	util.RunCommandAndLog(exec.Command("rm", "-r", zipPath))
}

//GetZippedBinaryFileSize zips the binary and returns its size
func GetZippedBinaryFileSize(experimentID int, binaryPath string) int64 {
	log.Infof("[sub-experiment %d] Zipping binary file to find its size...", experimentID)

	util.RunCommandAndLog(exec.Command("zip", "zipped-binary", binaryPath))

	fi, err := os.Stat("zipped-binary.zip")
	if err != nil {
		log.Fatalf("Could not get size of zipped binary file: %s", err.Error())
	}
	zippedBinarySizeBytes := fi.Size()

	log.Debug("Cleaning up zipped binary...")
	util.RunCommandAndLog(exec.Command("rm", "-r", "zipped-binary.zip"))

	return zippedBinarySizeBytes
}

//GenerateZIP creates the zip file for deployment
func GenerateZIP(experimentID int, fillerFileName string, binaryPath string) string {
	log.Infof("[sub-experiment %d] Generating ZIP file to be deployed...", experimentID)
	const localZipName = "benchmarking.zip"

	util.RunCommandAndLog(exec.Command("zip", localZipName, binaryPath, fillerFileName))

	util.RunCommandAndLog(exec.Command("rm", "-r", fillerFileName))

	log.Infof("[sub-experiment %d] Successfully generated ZIP file.", experimentID)

	workingDirectory, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(workingDirectory, localZipName)
}
