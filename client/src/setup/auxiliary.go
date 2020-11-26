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

package setup

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"strconv"
	"strings"
	"vhive-bench/client/setup/deployment/connection"
)

func removeEndpointFromSlice(s []connection.Endpoint, i int) []connection.Endpoint {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func almostEqualFloats(a, b float64, float64EqualityThreshold float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func promptForBool(prompt string) bool {
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

func promptForNumber(prompt string) *int64 {
	reader := bufio.NewReader(os.Stdin)

	log.Print(prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Could not read response: %s.", err.Error())
	}

	if response == "\n" {
		return nil
	}

	response = strings.ReplaceAll(response, "\n", "")

	parsedNumber, err := strconv.ParseInt(response, 10, 64)
	if err != nil {
		log.Fatalf("Could not parse integer %s: %s.", response, err.Error())
	}
	return &parsedNumber
}
