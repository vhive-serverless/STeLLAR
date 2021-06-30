// MIT License
//
// Copyright (c) 2021 Theodor Amariucai, Michal Baczun
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

package p

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func invokeNextFunctionGoogle(parameters map[string]string, functionID string) []byte {
	rawQuery := fmt.Sprintf("IncrementLimit=%s&TimestampChain=%v&DataTransferChainIDs=%v",
		parameters["IncrementLimit"],
		parameters["TimestampChain"],
		parameters["DataTransferChainIDs"],
	)

	finalURL := fmt.Sprintf("%s?%s", functionID, rawQuery)

	log.Printf("Invoking next function: %s", finalURL)

	responseBody := bytes.NewBuffer([]byte(parameters["TransferPayload"]))

	resp, err := http.Post(finalURL, "application/json", responseBody)
	if err != nil {
		log.Fatalf("Error while issuing http Post request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading http response body: %v", err)
	}
	return body
}

func authenticateCloudStorageClient() (*storage.Client, context.Context) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create cloud storage client: %v", err)
	}

	return client, ctx
}
