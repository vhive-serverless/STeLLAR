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

package main

import (
	"flag"
	"functions/connection"
	"functions/generator"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var rangeFlag = flag.String("range", "1_5", "Action functions with IDs in the given interval.")
var actionFlag = flag.String("action", "update", "Desired interaction with the functions (deploy, remove, update).")
var logLevelFlag = flag.String("logLevel", "info", "Select logging level.")

// https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html
var languageFlag = flag.String("language", "go1.x", "Programming language to deploy in.")

func main() {
	startTime := time.Now()
	flag.Parse()

	interval := strings.Split(*rangeFlag, "_")
	start, _ := strconv.Atoi(interval[0])
	end, _ := strconv.Atoi(interval[1])

	outputDirectoryPath := filepath.Join("logs", time.Now().Format(time.RFC850))
	log.Infof("Creating directory for this run at `%s`", outputDirectoryPath)
	if err := os.MkdirAll(outputDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFile := setupLogging(outputDirectoryPath)
	defer logFile.Close()

	provider := "aws"

	connection.Initialize(provider)

	sizeBytes := int(45000000 * 1.05) // ~5% is lost when compressing...
	generator.SetupDeployment(*actionFlag, *languageFlag, provider, sizeBytes)

	for id := start; id < end; id++ {
		switch *actionFlag {
		case "deploy":
			connection.Singleton.DeployFunction(id, *languageFlag, 128)
		case "remove":
			connection.Singleton.RemoveFunction(id)
		case "update":
			connection.Singleton.UpdateFunction(id, 128)
		default:
			log.Fatalf("Unrecognized function action %s", *actionFlag)
		}

		// AWS doesn't support simultaneous requests, or requests issued too quickly
		time.Sleep(time.Millisecond * 150)
	}

	log.Infof("Done in %v, exiting...", time.Since(startTime))
}

func setupLogging(path string) *os.File {
	logFile, err := os.Create(filepath.Join(path, "run_logs.txt"))
	if err != nil {
		log.Fatal(err)
	}

	switch *logLevelFlag {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}

	stdoutFileMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(stdoutFileMultiWriter)

	return logFile
}
