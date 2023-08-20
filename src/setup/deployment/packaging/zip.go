// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
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
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"stellar/setup/deployment/connection/amazon"
	"stellar/util"
)

// SetupZIPDeployment will package the function using ZIP
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

// GetZippedBinaryFileSize zips the binary and returns its size
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

func GenerateFillerFile(experimentID int, fillerFilePath string, sizeBytes int64) {
	log.Infof("[sub-experiment %d] Generating filler file to be included in deployment...", experimentID)

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("[sub-experiment %d] Failed to fill buffer with random bytes: `%s`", experimentID, err.Error())
	}

	if err := os.WriteFile(fillerFilePath, buffer, 0666); err != nil {
		log.Fatalf("[sub-experiment %d] Could not generate random file with size %d bytes: %v", experimentID, sizeBytes, err)
	}

	log.Infof("[sub-experiment %d] Successfully generated the filler file.", experimentID)
}

// GenerateZIP creates the zip file for deployment
func GenerateZIP(experimentID int, fillerFilePath string, binaryPath string, zipName string) string {
	log.Infof("[sub-experiment %d] Generating ZIP file to be deployed...", experimentID)

	util.RunCommandAndLog(exec.Command("zip", "-j", zipName, binaryPath, fillerFilePath))

	util.RunCommandAndLog(exec.Command("rm", "-r", fillerFilePath))

	log.Infof("[sub-experiment %d] Successfully generated ZIP file.", experimentID)

	workingDirectory, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(workingDirectory, zipName)
}

func GenerateServerlessZIPArtifacts(experimentID int, provider string, runtime string, functionName string, functionImageSizeMB float64) {
	switch runtime {
	case "python3.9":
		fallthrough
	case "go1.x":
		generateServerlessZIPArtifactsPythonGolang(experimentID, provider, runtime, functionName, functionImageSizeMB)
	case "java11":
		generateServerlessZIPArtifactsJava(experimentID, provider, runtime, functionName, functionImageSizeMB)
	}
}

func generateServerlessZIPArtifactsPythonGolang(experimentID int, provider string, runtime string, functionName string, functionImageSizeMB float64) {
	defaultBinaryName := map[string]string{
		"python3.9": "lambda_function.py",
		"go1.x":     "main",
	}
	binaryPath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/artifacts/%s/%s", provider, functionName, defaultBinaryName[runtime])
	currentSizeInBytes := GetZippedBinaryFileSize(experimentID, binaryPath)
	targetSizeInBytes := util.MBToBytes(functionImageSizeMB)

	if targetSizeInBytes == 0 {
		log.Infof("[sub-experiment %d] Desired image size is set to default (0MB), assigning size of zipped binary (%vMB)...",
			experimentID,
			util.BytesToMB(currentSizeInBytes),
		)
		targetSizeInBytes = currentSizeInBytes
	}
	if targetSizeInBytes < currentSizeInBytes {
		log.Fatalf("[sub-experiment %d] Total size (~%vMB) cannot be smaller than zipped binary size (~%vMB).",
			experimentID,
			util.BytesToMB(targetSizeInBytes),
			util.BytesToMB(currentSizeInBytes),
		)
	}

	fillerFilePath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/artifacts/%s/filler.file", provider, functionName)
	GenerateFillerFile(experimentID, fillerFilePath, targetSizeInBytes-currentSizeInBytes)
	zipPath := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/artifacts/%s/%s.zip", provider, functionName, functionName)
	GenerateZIP(experimentID, fillerFilePath, binaryPath, zipPath)
}

func generateServerlessZIPArtifactsJava(experimentID int, provider string, runtime string, functionName string, functionImageSizeMB float64) {
	// TODO generate filler file and add it to the existing ZIP archive built by Gradle
}
