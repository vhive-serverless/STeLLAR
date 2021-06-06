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

package p

import (
	"cloud.google.com/go/storage"
	"context"
	log "github.com/sirupsen/logrus"
)

func invokeNextFunctionGoogle(parameters map[string]string, functionID string) []byte {
	//type Payload struct {
	//	QueryStringParameters map[string]string `json:"queryStringParameters"`
	//}
	//nextFunctionPayload, err := json.Marshal(Payload{QueryStringParameters: parameters})
	//if err != nil {
	//	log.Fatalf("Could not marshal nextFunctionPayload: %s", err)
	//}
	//
	//log.Printf("Invoking next function: %s", functionID)
	//cloudFunctionsClient := authenticateCloudFunctionsClient()
	//result, err := cloudFunctionsClient. .Invoke(&lambdaSDK.InvokeInput{
	//	FunctionName:   aws.String(fmt.Sprintf("%s%s", namingPrefix, functionID)),
	//	InvocationType: aws.String("RequestResponse"),
	//	LogType:        aws.String("Tail"),
	//	Payload:        nextFunctionPayload,
	//})
	//if err != nil {
	//	log.Fatalf("Could not invoke lambda: %s", err)
	//}
	//
	//return result.Payload
	return nil
}

func authenticateCloudStorageClient() *storage.Client {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("authenticateCloudStorageClient threw %q", err)
	}

	return client
}