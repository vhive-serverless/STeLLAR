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

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
	"vhive-bench/client/setup/deployment/raw-code/producer-consumer/go1.x/common"
	"vhive-bench/client/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen"
)

const (
	port = ":50051"
)

type server struct {
	proto_gen.UnimplementedProducerConsumerServer
}

func main() {
	log.Printf("Started listening on port %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	log.Print("Created new server")

	proto_gen.RegisterProducerConsumerServer(s, &server{})
	log.Print("Registered ProducerConsumerServer")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) InvokeNext(_ context.Context, request *proto_gen.InvokeChainRequest) (*proto_gen.InvokeChainReply, error) {
	var updatedTimestampChain []string

	if firstFunctionInChain(request) {
		request.TransferPayload = string(common.GeneratePayload(request.PayloadLengthBytes))
		updatedTimestampChain = common.AppendTimestampToChain([]string{})
	} else {
		updatedTimestampChain = common.AppendTimestampToChain(common.StringArrayToArrayOfString(request.TimestampChain))
	}

	common.SimulateWork(request.IncrementLimit)

	dataTransferChainIDs := common.StringArrayToArrayOfString(request.DataTransferChainIDs)
	if common.FunctionsLeftInChain(dataTransferChainIDs) {
		log.Printf("There are %d functions left in the chain, invoking next one...", len(dataTransferChainIDs))
		updatedTimestampChain = invokeNextFunction(request, fmt.Sprintf("%v", updatedTimestampChain), dataTransferChainIDs)
	}

	return &proto_gen.InvokeChainReply{
		TimestampChain: fmt.Sprintf("%v", updatedTimestampChain),
	}, nil
}

func firstFunctionInChain(request *proto_gen.InvokeChainRequest) bool {
	return request.TimestampChain == ""
}

func invokeNextFunction(request *proto_gen.InvokeChainRequest, updatedTimestampChainString string, dataTransferChainIDs []string) []string {
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
		UseS3:                request.UseS3,
		IncrementLimit:       request.IncrementLimit,
		DataTransferChainIDs: request.DataTransferChainIDs[1:],
		TransferPayload:      request.TransferPayload,
		TimestampChain:       updatedTimestampChainString,
		S3Bucket:             request.S3Bucket,
		S3AccessKey:          request.S3AccessKey,
		S3SecretKey:          request.S3SecretKey,
	})
	if err != nil {
		log.Fatalf("could not create new producer consumer client: %v", err)
	}

	return common.StringArrayToArrayOfString(client.GetTimestampChain())
}
