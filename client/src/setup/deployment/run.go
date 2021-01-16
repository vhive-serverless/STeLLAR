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
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"vhive-bench/client/util"
)

//SetupDeployment will create the serverless function zip deployment for the given provider,
//in the given language and of the given size in bytes. Returns size of deployment in MB.
func SetupDeployment(rawCodePath string, provider string, language string, deploymentSizeBytes int64, packageType string) float64 {
	generateBinaryFile(rawCodePath, language)

	switch packageType {
	case "Zip":
		setupZIPDeployment(provider, deploymentSizeBytes)
	case "Image":
		log.Warn("Container image deployment does not support code size verification on AWS.")

		setupContainerImageDeployment(provider, deploymentSizeBytes)
	default:
		log.Fatalf("Unrecognized package type: %s", packageType)
	}

	return util.BytesToMB(deploymentSizeBytes)
}

func generateBinaryFile(rawCodePath string, runtime string) int64 {
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

	fi, err := os.Stat(util.BinaryName)
	if err != nil {
		log.Fatalf("Could not get size of binary file: %s", err.Error())
	}

	log.Info("Successfully built binary file for deployment...")
	return fi.Size()
}

func generateFillerFile(name string, sizeBytes int64) string {
	log.Info("Generating filler file to be included in deployment...")

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("Failed to fill buffer with random bytes: `%s`", err.Error())
	}

	if err := ioutil.WriteFile(name, buffer, 0666); err != nil {
		log.Fatalf("Could not generate random file with size %d bytes", sizeBytes)
	}

	log.Info("Successfully generated random file.")
	return name
}
