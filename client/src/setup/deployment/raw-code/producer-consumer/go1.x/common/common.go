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

package common

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//FunctionsLeftInChain checks if there are functions left in the chain
func FunctionsLeftInChain(dataTransferChainIDs []string) bool {
	return len(dataTransferChainIDs) > 0 && dataTransferChainIDs[0] != ""
}

//GeneratePayload creates a transfer payload for the producer-consumer chain
func GeneratePayload(payloadLengthBytesString string) []byte {
	payloadLengthBytes, err := strconv.Atoi(payloadLengthBytesString)
	if err != nil {
		log.Fatalf("Could not parse PayloadLengthBytes: %s", err)
	}

	log.Printf("Generating transfer payload for producer-consumer chain (length %d bytes)", payloadLengthBytes)
	generatedTransferPayload := make([]byte, payloadLengthBytes)
	rand.Read(generatedTransferPayload)

	return generatedTransferPayload
}

//ExtractJSONTimestampChain will process raw bytes into a string array of timestamps
func ExtractJSONTimestampChain(responsePayload []byte) []string {
	var reply map[string]interface{}
	err := json.Unmarshal(responsePayload, &reply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response into map[string]interface{}: %s", err)
	}

	var parsedReply ProducerConsumerResponse
	err = json.Unmarshal([]byte(reply["body"].(string)), &parsedReply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response body into producerConsumerResponse: %s", err)
	}

	return parsedReply.TimestampChain
}

//AppendTimestampToChain will add a new timestamp to the chain
func AppendTimestampToChain(timestampChain []string) []string {
	timestampMilliString := strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)
	return append(timestampChain, timestampMilliString)
}

//SimulateWork will keep the CPU busy-spinning
func SimulateWork(incrementLimitString string) {
	incrementLimit, err := strconv.Atoi(incrementLimitString)
	if err != nil {
		log.Fatalf("Could not parse IncrementLimit parameter: %s", err.Error())
	}

	log.Printf("Running function up to increment limit (%d)...", incrementLimit)
	for i := 0; i < incrementLimit; i++ {
	}
}

//ProducerConsumerResponse is the structure that we expect a consumer-producer function response to follow
type ProducerConsumerResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

//StringArrayToArrayOfString will process, e.g., "[14 35 8]" into []string{14, 35, 8}
func StringArrayToArrayOfString(str string) []string {
	str = strings.Split(str, "]")[0]
	str = strings.Split(str, "[")[1]
	return strings.Split(str, " ")
}
