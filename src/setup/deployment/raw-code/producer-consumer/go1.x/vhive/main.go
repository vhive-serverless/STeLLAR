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
	common2 "common"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	proto_gen2 "proto_gen"
)

const (
	port = ":50051"
)

type server struct {
	proto_gen2.UnimplementedProducerConsumerServer
}

func main() {
	log.Printf("Started listening on port %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	log.Print("Created new server")

	proto_gen2.RegisterProducerConsumerServer(s, &server{})
	log.Print("Registered ProducerConsumerServer")

	common2.InitializeGlobalRandomPayload()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) InvokeNext(ctx context.Context, request *proto_gen2.InvokeChainRequest) (*proto_gen2.InvokeChainReply, error) {
	_, grpcOutput := common2.GenerateResponse(ctx, nil, request)

	return &proto_gen2.InvokeChainReply{
		TimestampChain: fmt.Sprintf("%v", grpcOutput),
	}, nil
}
