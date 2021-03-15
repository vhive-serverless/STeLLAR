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
	"github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/common"
	"github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen"
	"google.golang.org/grpc"
	"log"
	"net"
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

	common.InitializeGlobalRandomPayload()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) InvokeNext(ctx context.Context, request *proto_gen.InvokeChainRequest) (*proto_gen.InvokeChainReply, error) {
	_, grpcOutput := common.GenerateResponse(ctx, nil, request)

	return &proto_gen.InvokeChainReply{
		TimestampChain: fmt.Sprintf("%v", grpcOutput),
	}, nil
}
