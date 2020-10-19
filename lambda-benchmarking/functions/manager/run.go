package main

import (
	"flag"
	"functions/manager/provider"
	"functions/manager/writer"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rangeFlag = flag.String("range", "0_300", "Action functions with IDs in the given interval.")
var actionFlag = flag.String("action", "deploy", "Desired interaction with the functions.")
var providerFlag = flag.String("provider", "aws", "ProviderName where functions are located.")

func main() {
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

	connection := &provider.Connection{ProviderName: *providerFlag}

	rateLimiter := 1
	var deployWaitGroup sync.WaitGroup
	for i := start; i < end; {
		for requests := 0; requests < rateLimiter && i < end; requests++ {
			deployWaitGroup.Add(1)
			go func(deployWaitGroup *sync.WaitGroup, i int) {
				defer deployWaitGroup.Done()

				switch *actionFlag {
				case "deploy":
					csvFile, err := os.Create(filepath.Join(outputDirectoryPath, "gateways.csv"))
					if err != nil {
						log.Fatal(err)
					}
					writer.InitializeGatewaysWriter(csvFile)
					defer csvFile.Close()

					connection.DeployFunction(i)
				case "remove":
					connection.RemoveFunction(i)
				case "update":
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
