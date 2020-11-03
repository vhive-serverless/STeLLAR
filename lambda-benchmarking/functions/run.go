package main

import (
	"flag"
	"functions/provider"
	"functions/util"
	"functions/writer"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var rangeFlag = flag.String("range", "0_5", "Action functions with IDs in the given interval.")
var actionFlag = flag.String("action", "deploy", "Desired interaction with the functions (deploy, "+
	"remove, update_func, update_func_conf).")
var providerFlag = flag.String("provider", "aws", "Provider to interact with.")
var sizeBytesFlag = flag.Int("sizeBytes", 0, "The size of the image to deploy, in bytes.")
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

	if deploymentFile := setupDeployment(outputDirectoryPath); deploymentFile != nil {
		defer deploymentFile.Close()
	}

	connection := &provider.Connection{ProviderName: *providerFlag}

	// Issuing simultaneous requests poses problems with AWS
	for i := start; i < end; i++ {
		switch *actionFlag {
		case "deploy":
			connection.DeployFunction(i, *languageFlag)
		case "remove":
			connection.RemoveFunction(i)
		case "update_func":
			connection.UpdateFunction(i)
		case "update_func_conf":
			connection.UpdateFunctionConfiguration(i)
		default:
			log.Fatalf("Unrecognized function action %s", *actionFlag)
		}
	}

	if *actionFlag == "deploy" {
		log.Info("Flushing gateways to CSV file.")
		writer.GatewaysWriterSingleton.Writer.Flush()
	}

	log.Infof("Done in %v, exiting...", time.Since(startTime))
}

func setupDeployment(outputDirectoryPath string) *os.File {
	switch *actionFlag {
	case "deploy":
		deploymentFile, err := os.Create(filepath.Join(outputDirectoryPath, "gateways.csv"))
		if err != nil {
			log.Fatal(err)
		}
		writer.InitializeGatewaysWriter(deploymentFile)

		util.GenerateDeploymentZIP(*providerFlag, *languageFlag, *sizeBytesFlag)
		return deploymentFile
	case "update_func":
		util.GenerateDeploymentZIP(*providerFlag, *languageFlag, *sizeBytesFlag)
	case "update_func_conf":
		util.GenerateDeploymentZIP(*providerFlag, *languageFlag, *sizeBytesFlag)
	case "remove":
		// No setup required for removing functions
	default:
		log.Fatalf("Unrecognized function action %s", *actionFlag)
	}
	return nil
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
