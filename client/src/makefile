PROTOCOL_LOCATION := "experiment/networking/gRPC/proto"
GO_API_OUT := "experiment/networking/gRPC/proto_gen"

proto_gen_go:
	@mkdir -p $(GO_API_OUT)
	@protoc vHive.proto --proto_path=$(PROTOCOL_LOCATION) --go_out=plugins=grpc:$(GO_API_OUT)

build:
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@golangci-lint run ./...

build_client: proto_gen_go build

run_client:
	@mkdir -p "../latency-samples"
	@./client \
		-o "../latency-samples" \
		-c "../experiments/test.json"\
		-g "../endpoints"

build_and_run_client: build_client run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

workflows_edit:
	@gedit ../../.github/workflows/go.yml

replace_github_commit:
	@echo "WARNING: This will latest remote commit with a new one having the newest changes, risks are involved."
	@git add ../../.
	@git commit --no-edit --amend
	@git push --force

# Cloudlab deployment
remote=theodor@hp034.utah.cloudlab.us

send_to_cloudlab:
	@ssh $(remote) "mkdir -p ~/.aws/"
	@scp /home/theodor/.aws/credentials $(remote):~/.aws/
	@scp ../initialize_machine.sh $(remote):~/
	@scp ../endpoints/* $(remote):~/
	@ssh $(remote) "mkdir -p ~/experiments/"
	@scp ../experiments/* $(remote):~/experiments/
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@scp ./client $(remote):~/

connect_to_cloudlab:
	@ssh -X -p 22 $(remote)

download_cloudlab_results:
	@mkdir -p "../latency-samples/cloudlab"
	@scp -r $(remote):~/latency-samples/* "../latency-samples/cloudlab/"

.PHONY: proto_gen_go build build_client run_client build_and_run_client workflows_edit replace_github_commit send_code_to_cloudlab connect_to_cloudlab download_cloudlab_results
