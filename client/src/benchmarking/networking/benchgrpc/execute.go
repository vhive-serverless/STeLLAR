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

package benchgrpc

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
	"vhive-bench/client/benchmarking/networking/benchgrpc/proto_gen"
	"vhive-bench/client/setup"
	"vhive-bench/client/setup/deployment/connection/amazon"
)

const port = 80

//ExecuteRequest will send a gRPC request and return the timestamp chain (if any).
func ExecuteRequest(payloadLengthBytes int, gatewayEndpoint setup.GatewayEndpoint, incrementLimit int64, s3Transfer bool) (string, time.Time, time.Time) {
	address := fmt.Sprintf("%s:%d", gatewayEndpoint.ID, port)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client := proto_gen.NewProducerConsumerClient(conn)

	input := &proto_gen.InvokeChainRequest{
		IncrementLimit:       fmt.Sprintf("%d", incrementLimit),
		DataTransferChainIDs: fmt.Sprintf("%v", gatewayEndpoint.DataTransferChainIDs),
		PayloadLengthBytes:   fmt.Sprintf("%d", payloadLengthBytes),
	}

	if s3Transfer {
		input.S3Bucket = amazon.AWSBucketName
	}

	var reply *proto_gen.InvokeChainReply

	reqSentTime := time.Now()
	reply, err = client.InvokeNext(ctx, input)
	if err != nil {
		log.Fatalf("Could not invoke gRPC function: %v", err)
	}
	reqReceivedTime := time.Now()

	return reply.GetTimestampChain(), reqSentTime, reqReceivedTime
}
