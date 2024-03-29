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
	"os/exec"
	"stellar/setup/deployment/connection/amazon"
	"stellar/util"
)

//SetupContainerImageDeployment will package the function using container images
func SetupContainerImageDeployment(function string, provider string, rawCodePath string) {
	var privateRepoURI string
	switch provider {
	case "aws":
		privateRepoURI = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", amazon.UserARNNumber, amazon.AWSRegion)

		log.Info("Authenticating Docker CLI to the Amazon ECR registry...")
		util.RunCommandAndLog(exec.Command("docker", "login", "-u", "AWS", "-p",
			amazon.GetECRAuthorizationToken(), privateRepoURI))
	case "vhive":
		log.Info("Authenticating Docker CLI to the DockerHub registry...")

		privateRepoURI = *promptForString("Please enter your DockerHub username: ")
		util.RunCommandAndLog(exec.Command("docker", "login", "-u",
			privateRepoURI, "-p", *promptForString("Please enter your DockerHub password: ")))
	default:
		log.Warnf("Provider %s does not support container image deployment, skipping...", provider)
		return
	}

	taggedImage := fmt.Sprintf("%s:latest", function)

	util.RunCommandAndLog(exec.Command("docker", "build", "-t", taggedImage, rawCodePath))

	imageName := fmt.Sprintf("%s/%s", privateRepoURI, taggedImage)
	log.Infof("Pushing container image to %q...", imageName)

	util.RunCommandAndLog(exec.Command("docker", "tag", taggedImage, imageName))
	util.RunCommandAndLog(exec.Command("docker", "push", imageName))

	if provider == "aws" {
		amazon.AWSSingletonInstance.ImageURI = imageName
	}
}
