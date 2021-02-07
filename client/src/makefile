# MIT License
#
# Copyright (c) 2021 Theodor Amariucai
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

# Cloudlab deployment
CLOUDLAB_REMOTE=theodor@hp062.utah.cloudlab.us
PROTOCOL_LOCATION = "setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto"

build_and_run_client: build_client run_client_aws

build_client: proto_gen_go build test

SERVER_API_OUT = "setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen"
CLIENT_API_OUT = "benchmarking/networking/benchgrpc/proto_gen"
proto_gen_go:
	@mkdir -p $(CLIENT_API_OUT)
	@mkdir -p $(SERVER_API_OUT)
	@protoc chainfunction.proto --proto_path=$(PROTOCOL_LOCATION) --go_out=$(CLIENT_API_OUT) --go-grpc_out=$(CLIENT_API_OUT)
	@protoc chainfunction.proto --proto_path=$(PROTOCOL_LOCATION) --go_out=$(SERVER_API_OUT) --go-grpc_out=$(SERVER_API_OUT)

VHIVE_LOCATION = "setup/deployment/raw-code/producer-consumer/go1.x/vhive"
AWS_LOCATION = "setup/deployment/raw-code/producer-consumer/go1.x/aws"
build:
	@CGO_ENABLED=0 GOOS=linux cd $(VHIVE_LOCATION) && go build -v -o "handler" "main.go"
	@CCGO_ENABLED=0 GOOS=linux cd $(AWS_LOCATION) && go build -v -o "handler" "main.go"
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@golangci-lint run ./...

test:
	@go test -cover -short -v ./...  -coverprofile=c.out
	@go tool cover -html=c.out
	@rm c.out

run_client_aws:
	@mkdir -p "../latency-samples"
	@./client -o "../latency-samples" -g "../endpoints" -c "../experiments/tests/AWS/burstiness-zip-package.json"

run_client_vhive:
	@mkdir -p "../latency-samples"
	@./client -o "../latency-samples" -g "../endpoints" -c "../experiments/tests/vhive/burstiness-image-package.json"

edit_local_env_vars:
	@sudo kate /etc/profile.d/benchmarking.sh

workflows_edit:
	@kate ../../.github/workflows/pipeline.yml

replace_github_commit: ## This will replace latest CLOUDLAB_REMOTE commit with a new one having the newest changes, risks are involved.
	@echo "WARNING: This will replace latest CLOUDLAB_REMOTE commit with a new one having the newest changes, risks are involved."
	@git add ../../.
	@git commit -S --no-edit --amend
	@git push --force

start_docker_daemon:
	@sudo systemctl start docker

vhive_local_docker_deployment:
	@docker system prune
	@docker build -t vhive-bench:latest setup/deployment/raw-code/producer-consumer/go1.x/vhive/
	@docker run -dp 50051:50051 vhive-bench:latest

send_to_cloudlab:
	@ssh $(CLOUDLAB_REMOTE) "mkdir -p ~/vhive-bench-deployment/"
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@scp ./client $(CLOUDLAB_REMOTE):~/vhive-bench-deployment/
	@scp ./setup/deployment/raw-code/producer-consumer/api-template.json $(CLOUDLAB_REMOTE):~/vhive-bench-deployment/
	@scp ../initialize_machine.sh $(CLOUDLAB_REMOTE):~/vhive-bench-deployment/
	@ssh $(CLOUDLAB_REMOTE) "mkdir -p ~/vhive-bench-deployment/endpoints/"
	@scp -r ../endpoints/* $(CLOUDLAB_REMOTE):~/vhive-bench-deployment/endpoints/
	@ssh $(CLOUDLAB_REMOTE) "mkdir -p ~/vhive-bench-deployment/experiments/"
	@scp -r ../experiments/* $(CLOUDLAB_REMOTE):~/vhive-bench-deployment/experiments/
	@ssh $(CLOUDLAB_REMOTE) "mkdir -p ~/vhive-bench-deployment/setup/deployment/"
	@scp -rp ./setup/deployment/raw-code/ $(CLOUDLAB_REMOTE):~/vhive-bench-deployment/setup/deployment/

connect_to_cloudlab:
	@ssh -X -p 22 $(CLOUDLAB_REMOTE)

download_cloudlab_results:
	@mkdir -p "../latency-samples/cloudlab"
	@scp -r $(CLOUDLAB_REMOTE):~/latency-samples/* "../latency-samples/cloudlab/"

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: proto_gen_go build test build_client run_client_vhive run_client_aws build_and_run_client workflows_edit replace_github_commit start_docker_daemon send_to_cloudlab connect_to_cloudlab download_cloudlab_results
