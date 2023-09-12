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

package packaging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"stellar/setup/deployment/connection/amazon"
	"stellar/util"
)

var builtImages = make(map[string]bool)
var privateRepoURI string = ""
var loggedIn bool = false

// SetupContainerImageDeployment will package the function using container images and push to registry
func SetupContainerImageDeployment(function string, provider string) string {
	functionDir := fmt.Sprintf("setup/deployment/raw-code/serverless/%s/%s", provider, function)
	switch provider {
	case "aws":
		privateRepoURI = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", amazon.UserARNNumber, amazon.AWSRegion)

		log.Info("Authenticating Docker CLI to the Amazon ECR registry...")
		util.RunCommandAndLog(exec.Command("docker", "login", "-u", "AWS", "-p",
			amazon.GetECRAuthorizationToken(), privateRepoURI))
	case "gcr":
		fallthrough
	case "vhive":
		log.Info("Authenticating Docker CLI to the DockerHub registry...")

		if !loggedIn {
			privateRepoURI = os.Getenv("DOCKER_HUB_USERNAME")
			privateRepoToken := os.Getenv("DOCKER_HUB_ACCESS_TOKEN")
			util.RunCommandAndLog(exec.Command("docker", "login", "-u",
				privateRepoURI, "-p", privateRepoToken))
			loggedIn = true
		}
	default:
		log.Fatalf("Provider %s does not support container image deployment.", provider)
	}

	taggedImage := fmt.Sprintf("%s_stellar:latest", function)
	imageName := fmt.Sprintf("%s/%s", privateRepoURI, taggedImage)
	if builtImages[function] {
		log.Infof("Container image for function %q is already built. Skipping...", function)
		return imageName
	}

	util.RunCommandAndLog(exec.Command("docker", "build", "-t", taggedImage, functionDir))

	log.Infof("Pushing container image to %q...", imageName)

	util.RunCommandAndLog(exec.Command("docker", "tag", taggedImage, imageName))
	util.RunCommandAndLog(exec.Command("docker", "push", imageName))

	if provider == "aws" {
		amazon.AWSSingletonInstance.ImageURI = imageName
	}
	builtImages[function] = true
	return imageName
}
