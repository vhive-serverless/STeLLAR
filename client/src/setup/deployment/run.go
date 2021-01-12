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
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"vhive-bench/client/setup/deployment/connection/amazon"
	"vhive-bench/client/util"
)

const (
	localZipName   = "benchmarking.zip"
	randomFileName = "random.file"
)

//SetupDeployment will create the serverless function zip deployment for the given provider,
//in the given language and of the given size in bytes. Returns size of deployment in MB.
func SetupDeployment(rawCodePath string, provider string, language string, sizeBytes int64) float64 {
	zippedBinarySizeBytes := createBinary(rawCodePath, language)

	if sizeBytes == 0 {
		log.Infof("Desired image size is set to default (0MB), assigning size of zipped binary (%vMB)...",
			util.BytesToMB(zippedBinarySizeBytes))
		sizeBytes = zippedBinarySizeBytes
	}

	if sizeBytes < zippedBinarySizeBytes {
		log.Fatalf("Total size (~%vMB) cannot be smaller than zipped binary size (~%vMB).",
			util.BytesToMB(sizeBytes),
			util.BytesToMB(zippedBinarySizeBytes),
		)
	}

	generateRandomFile(sizeBytes - zippedBinarySizeBytes)
	zipPath := generateZIP()

	sizeMB := util.BytesToMB(sizeBytes)

	switch provider {
	case "aws":
		if sizeMB > 50. {
			amazon.UploadZIPToS3(zipPath, sizeMB)
		} else {
			amazon.SetLocalZip(zipPath)
		}
	default:
		log.Warnf("Provider %s does not support code deployment, skipping ZIP generation...", provider)
	}

	log.Debugf("Cleaning up ZIP %q...", zipPath)
	util.RunCommandAndLog(exec.Command("rm", "-r", zipPath))

	return util.BytesToMB(sizeBytes)
}

func createBinary(rawCodePath string, runtime string) int64 {
	log.Info("Building binary file for the function(s) to be deployed...")

	if !util.FileExists(rawCodePath) {
		log.Fatalf("Code path `%s` does not exist, cannot deploy/update raw code.", rawCodePath)
	}

	switch runtime {
	case "go1.x":
		util.RunCommandAndLog(exec.Command("go", "build", "-v", "-o", util.BinaryName, rawCodePath))
	//TODO: add python3 support
	//case "python3.8":
	//	runCommandAndLog(exec.Command("python", "build", "-v", "-race", "-o", util.BinaryName))
	default:
		log.Fatalf("Unrecognized runtime %s", runtime)
	}

	log.Info("Zipping binary file to find its size...")
	util.RunCommandAndLog(exec.Command("zip", "zipped-binary", util.BinaryName))
	fi, err := os.Stat("zipped-binary.zip")
	if err != nil {
		log.Fatalf("Could not get size of zipped binary file: %s", err.Error())
	}

	log.Debug("Cleaning up zipped binary...")
	util.RunCommandAndLog(exec.Command("rm", "-r", "zipped-binary.zip"))

	log.Info("Successfully built binary file for deployment...")
	return fi.Size()
}

func generateRandomFile(sizeBytes int64) {
	log.Info("Generating random file to be included in deployment...")

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("Failed to fill buffer with random bytes: `%s`", err.Error())
	}

	if err := ioutil.WriteFile(randomFileName, buffer, 0666); err != nil {
		log.Fatalf("Could not generate random file with size %d bytes", sizeBytes)
	}

	log.Info("Successfully generated random file.")
}

func generateZIP() string {
	log.Info("Generating ZIP file to be deployed...")

	util.RunCommandAndLog(exec.Command("zip", localZipName, util.BinaryName, randomFileName))

	log.Debugf("Cleaning up binary %q...", util.BinaryName)
	util.RunCommandAndLog(exec.Command("rm", "-r", util.BinaryName))

	log.Debugf("Cleaning up random file %q...", randomFileName)
	util.RunCommandAndLog(exec.Command("rm", "-r", randomFileName))

	log.Info("Successfully generated ZIP file.")

	workingDirectory, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(workingDirectory, localZipName)
}
