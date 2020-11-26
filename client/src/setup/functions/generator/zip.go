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

package generator

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/setup/functions/connection/amazon"
	"os"
	"os/exec"
)

func generateZIP(provider string, randomFileName string, sizeMB float64) {
	localZipPath := "benchmarking.zip"

	runCommandAndLog(exec.Command("zip", localZipPath, "producer-handler", randomFileName))

	switch provider {
	case "aws":
		if sizeMB > 50. {
			amazon.UploadZIPToS3(localZipPath, sizeMB)
		} else {
			amazon.SetLocalZip(localZipPath)
		}
	default:
		log.Warnf("Provider %s might not support code deployment, skipping ZIP generation...", provider)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
