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

package benchhttp

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"vhive-bench/client/setup/deployment/connection/amazon"
)

//ProducerConsumerResponse is the structure holding the response from a producer-consumer function
type ProducerConsumerResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

//ExtractProducerConsumerResponse will process an HTTP response body coming from a producer-consumer function
func ExtractProducerConsumerResponse(respBody []byte) ProducerConsumerResponse {
	var response ProducerConsumerResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		log.Error(err)
	}
	return response
}

func appendProducerConsumerParameters(request *http.Request, payloadLengthBytes int,
	assignedFunctionIncrementLimit int64, dataTransferChainIDs []string, storageTransfer bool) *http.Request {
	request.URL.Path = "/prod/benchmarking"

	request.URL.RawQuery = fmt.Sprintf("IncrementLimit=%d&PayloadLengthBytes=%d&DataTransferChainIDs=%v",
		assignedFunctionIncrementLimit,
		payloadLengthBytes,
		dataTransferChainIDs,
	)

	if storageTransfer {
		request.URL.RawQuery += fmt.Sprintf("&Bucket=%v&StorageTransfer=true", amazon.AWSSingletonInstance.S3Bucket)
	}

	return request
}
