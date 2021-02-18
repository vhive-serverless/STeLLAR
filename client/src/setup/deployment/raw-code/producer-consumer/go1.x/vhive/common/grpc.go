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
	"context"
	"fmt"
	"github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func invokeNextFunctionGRPC(request *proto_gen.InvokeChainRequest, updatedTimestampChainString string, dataTransferChainIDs []string) []string {
	address := fmt.Sprintf("%s:80", dataTransferChainIDs[0])

	log.Printf("Invoking next function: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client, err := proto_gen.NewProducerConsumerClient(conn).InvokeNext(ctx, &proto_gen.InvokeChainRequest{
		IncrementLimit:       request.IncrementLimit,
		DataTransferChainIDs: request.DataTransferChainIDs[1:],
		TransferPayload:      request.TransferPayload,
		TimestampChain:       updatedTimestampChainString,
		S3Bucket:             request.S3Bucket,
		S3Key:                request.S3Key,
	})
	if err != nil {
		log.Fatalf("could not create new producer consumer client: %v", err)
	}

	return StringArrayToArrayOfString(client.GetTimestampChain())
}
