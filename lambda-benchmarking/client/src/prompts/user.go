package prompts

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func PromptForBool(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		log.Printf("%s [y/n]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func PromptForNumber(prompt string) *int64 {
	reader := bufio.NewReader(os.Stdin)

	log.Print(prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Could not read response: %s.", err.Error())
	}

	if response == "\n" {
		return nil
	} else {
		response = strings.ReplaceAll(response, "\n", "")
	}

	parsedNumber, err := strconv.ParseInt(response, 10, 64)
	if err != nil {
		log.Fatalf("Could not parse integer %s: %s.", response, err.Error())
	}
	return &parsedNumber
}
