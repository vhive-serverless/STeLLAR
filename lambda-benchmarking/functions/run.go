package main

import (
	"flag"
	"functions/provider"
	"functions/util"
	"functions/writer"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rangeFlag = flag.String("range", "1_2", "Action functions with IDs in the given interval.")
var actionFlag = flag.String("action", "deploy", "Desired interaction with the functions.")
var providerFlag = flag.String("provider", "aws", "Provider to interact with.")
var sizeBytesFlag = flag.Int("sizeBytes", 0, "The size of the image to deploy, in bytes.")

// https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html
var languageFlag = flag.String("language", "go1.x", "Programming language to deploy in.")

func main() {
	rand.Seed(1896564)
	flag.Parse()

	interval := strings.Split(*rangeFlag, "_")
	start, _ := strconv.Atoi(interval[0])
	end, _ := strconv.Atoi(interval[1])

	outputDirectoryPath := filepath.Join("logs", time.Now().Format(time.RFC850))
	log.Printf("Creating directory for this run at `%s`", outputDirectoryPath)
	if err := os.MkdirAll(outputDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	logFile := setupManagerLogging(outputDirectoryPath)
	defer logFile.Close()

	if *actionFlag == "deploy" {
		csvFile, err := os.Create(filepath.Join(outputDirectoryPath, "gateways.csv"))
		if err != nil {
			log.Fatal(err)
		}
		writer.InitializeGatewaysWriter(csvFile)
		defer csvFile.Close()
	}

	connection := &provider.Connection{ProviderName: *providerFlag}

	rateLimiter := 15
	var deployWaitGroup sync.WaitGroup
	for i := start; i < end; {
		for requests := 0; requests < rateLimiter && i < end; requests++ {
			deployWaitGroup.Add(1)
			go func(deployWaitGroup *sync.WaitGroup, i int) {
				defer deployWaitGroup.Done()

				switch *actionFlag {
				case "deploy":
					util.GenerateDeploymentZIP(*providerFlag, *languageFlag, *sizeBytesFlag)
					connection.DeployFunction(i, *languageFlag)
				case "remove":
					connection.RemoveFunction(i)
				case "update":
					util.GenerateDeploymentZIP(*providerFlag, *languageFlag, *sizeBytesFlag)
					connection.UpdateFunction(i)
				default:
					log.Fatalf("Unrecognized function action %s", *actionFlag)
				}
			}(&deployWaitGroup, i)
			i++
		}
		deployWaitGroup.Wait()
	}

	if *actionFlag == "deploy" {
		log.Println("Flushing gateways to CSV file.")
		writer.GatewaysWriterSingleton.Writer.Flush()
	}

	log.Println("Done, exiting...")
}

func setupManagerLogging(path string) *os.File {
	logFile, err := os.Create(filepath.Join(path, "run_logs.txt"))
	if err != nil {
		log.Fatal(err)
	}
	stdoutFileMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(stdoutFileMultiWriter)
	return logFile
}
