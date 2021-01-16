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

package deployment

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"vhive-bench/client/setup/deployment/connection/amazon"
	"vhive-bench/client/util"
)

const localZipName = "benchmarking.zip"

func setupZIPDeployment(provider string, deploymentSizeBytes int64) {
	zippedBinaryFileSizeBytes := getZippedBinaryFileSize()

	if deploymentSizeBytes == 0 {
		log.Infof("Desired image size is set to default (0MB), assigning size of zipped binary (%vMB)...",
			util.BytesToMB(zippedBinaryFileSizeBytes))
		deploymentSizeBytes = zippedBinaryFileSizeBytes
	}

	if deploymentSizeBytes < zippedBinaryFileSizeBytes {
		log.Fatalf("Total size (~%vMB) cannot be smaller than zipped binary size (~%vMB).",
			util.BytesToMB(deploymentSizeBytes),
			util.BytesToMB(zippedBinaryFileSizeBytes),
		)
	}

	zipPath := generateZIP(
		generateFillerFile("random.file", deploymentSizeBytes-zippedBinaryFileSizeBytes),
	)

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

func getZippedBinaryFileSize() int64 {
	log.Info("Zipping binary file to find its size...")

	util.RunCommandAndLog(exec.Command("zip", "zipped-binary", util.BinaryName))

	fi, err := os.Stat("zipped-binary.zip")
	if err != nil {
		log.Fatalf("Could not get size of zipped binary file: %s", err.Error())
	}
	zippedBinarySizeBytes := fi.Size()

	log.Debug("Cleaning up zipped binary...")
	util.RunCommandAndLog(exec.Command("rm", "-r", "zipped-binary.zip"))

	return zippedBinarySizeBytes
}

func generateZIP(fillerFileName string) string {
	log.Info("Generating ZIP file to be deployed...")

	util.RunCommandAndLog(exec.Command("zip", localZipName, util.BinaryName, fillerFileName))

	log.Debugf("Cleaning up binary %q...", util.BinaryName)
	util.RunCommandAndLog(exec.Command("rm", "-r", util.BinaryName))

	log.Debugf("Cleaning up random file %q...", fillerFileName)
	util.RunCommandAndLog(exec.Command("rm", "-r", fillerFileName))

	log.Info("Successfully generated ZIP file.")

	workingDirectory, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(workingDirectory, localZipName)
}
