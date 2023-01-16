# MIT License
#
# Copyright (c) 2022 Dilina Dehigama and EASE Lab
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# Download protoc, protoc-gen-go, protoc-gen-go-grpc
sudo apt-get update && sudo apt-get install --no-install-recommends --assume-yes  protobuf-compiler \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
  && go install github.com/golang/protobuf/protoc-gen-go@latest

sudo cp ~/go/bin/protoc-gen-go /usr/bin/
sudo cp ~/go/bin/protoc-gen-go-grpc /usr/bin/


WORKING_DIR=$(pwd)
PROTOCOL_LOCATION="${WORKING_DIR}/setup/deployment/raw-code/proto"
GO_PROD_CONSUMER="${WORKING_DIR}/setup/deployment/raw-code/functions/producer-consumer"
SERVER_API_OUT1="${GO_PROD_CONSUMER}/setup/deployment/raw-code/functions/producer-consumer/vhive/proto_gen"
SERVER_API_OUT2="${GO_PROD_CONSUMER}/aws/proto_gen"
CLIENT_API_OUT="${WORKING_DIR}/benchmarking/networking/benchgrpc/proto_gen"

# Build the gRPC protocols
mkdir -p $CLIENT_API_OUT && mkdir -p $SERVER_API_OUT1 && mkdir -p $SERVER_API_OUT2
protoc chainfunction.proto --proto_path=$PROTOCOL_LOCATION --go_out=$SERVER_API_OUT1 --go-grpc_out=$SERVER_API_OUT1
protoc chainfunction.proto --proto_path=$PROTOCOL_LOCATION --go_out=$SERVER_API_OUT2 --go-grpc_out=$SERVER_API_OUT2
protoc chainfunction.proto --proto_path=$PROTOCOL_LOCATION --go_out=$CLIENT_API_OUT --go-grpc_out=$CLIENT_API_OUT

# Build binaries for vHive producer-consumer, AWS producer-consumer (ZIP deployment) and the client
cd "$GO_PROD_CONSUMER/vhive" \
   && go mod download \
   && go build -v -o handler main.go \
   && cd "$GO_PROD_CONSUMER/aws" \
   && go mod download \
   && go build -v -o handler main.go \
   && cd "$WORKING_DIR" \
   && echo $(pwd) \
   && go mod download \
   && go build -v -a
