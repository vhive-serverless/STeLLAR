package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var rangeFlag = flag.String("range", "0 30", "Action lambdas with IDs in the given interval.")
var actionFlag = flag.String("action", "deploy", "deploy/remove/update the Lambda functions?")

func main() {
	flag.Parse()

	log.Println("This is unstable at the moment")
	return

	interval := strings.Split(*rangeFlag, " ")
	start, _ := strconv.Atoi(interval[0])
	end, _ := strconv.Atoi(interval[1])

	var deployWaitGroup sync.WaitGroup
	for i := start; i < end; i++ {
		deployWaitGroup.Add(1)
		go func(deployWaitGroup *sync.WaitGroup, i int) {
			defer deployWaitGroup.Done()

			log.Printf("Actioning %s Lambda %d", *actionFlag, i)
			deployLambda(i)
			log.Printf("Lambda %d actioned!", i)
		}(&deployWaitGroup, i)
	}

	deployWaitGroup.Wait()
	log.Println("Done, exiting...")
}

func deployLambda(i int) {
	cmd := exec.Command("/bin/sh", fmt.Sprintf("./deployment_scripts/%s-producers.sh", *actionFlag), strconv.Itoa(i))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
