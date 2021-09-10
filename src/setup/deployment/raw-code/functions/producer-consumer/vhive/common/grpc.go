// MIT License
//
// Copyright (c) 2021 Theodor Amariucai and EASE Lab
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
	"context"
	"fmt"
	protogen2 "github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/functions/producer-consumer/proto_gen"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func invokeNextFunctionGRPC(request *protogen2.InvokeChainRequest, updatedTimestampChain []string, dataTransferChainIDs []string) []string {
	log.Printf("Invoking next function: %s", dataTransferChainIDs[0])
	conn, err := grpc.Dial(dataTransferChainIDs[0], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client, err := protogen2.NewProducerConsumerClient(conn).InvokeNext(ctx, &protogen2.InvokeChainRequest{
		IncrementLimit:       request.IncrementLimit,
		DataTransferChainIDs: fmt.Sprintf("%v", dataTransferChainIDs[1:]),
		TransferPayload:      request.TransferPayload,
		TimestampChain:       fmt.Sprintf("%v", updatedTimestampChain),
		Bucket:               request.Bucket,
		Key:                  request.Key,
	})
	if err != nil {
		log.Fatalf("could not create new producer consumer client: %v", err)
	}

	return StringArrayToArrayOfString(client.GetTimestampChain())
}
