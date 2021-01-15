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
	"os/exec"
	"vhive-bench/client/setup/deployment/connection/amazon"
	"vhive-bench/client/util"
)

func setupContainerImageDeployment(provider string, deploymentSizeBytes int64, language string) {
	var privateECRRepoURI string
	switch provider {
	case "aws":
		privateECRRepoURI = fmt.Sprintf("335329526041.dkr.ecr.%s.amazonaws.com", amazon.AWSRegion)
	default:
		log.Warnf("Provider %s does not support container image deployment, skipping...", provider)
		return
	}

	//deploymentSizeMB := util.BytesToMB(deploymentSizeBytes)
	// TODO: Add random file to Docker image

	log.Info("Building raw function code into a container image...")
	dockerfilePath := fmt.Sprintf("./setup/deployment/raw-code/%s/", language)
	util.RunCommandAndLog(exec.Command("docker", "build", "-t", "vhive-bench", dockerfilePath))

	log.Info("Authenticating Docker CLI to the Amazon ECR registry...")
	util.RunCommandAndLog(exec.Command("docker", "login", "-u", "AWS", "-p",
		amazon.GetECRAuthorizationToken(), privateECRRepoURI))

	log.Info("Pushing container image to the registry...")
	imageName := fmt.Sprintf("%s/%s", privateECRRepoURI, "vhive-bench:latest")
	util.RunCommandAndLog(exec.Command("docker", "tag", "vhive-bench:latest", imageName))
	util.RunCommandAndLog(exec.Command("docker", "push", imageName))

	amazon.AWSSingleton.DockerImageURI = imageName
}
