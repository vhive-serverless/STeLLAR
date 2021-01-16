PROTOCOL_LOCATION := "experiments/networking/gRPC/proto"
GO_API_OUT := "experiments/networking/gRPC/proto_gen"

proto_gen_go:
	@mkdir -p $(GO_API_OUT)
	@protoc vHive.proto --proto_path=$(PROTOCOL_LOCATION) --go_out=plugins=grpc:$(GO_API_OUT)

build:
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@golangci-lint run ./...

test:
	@go test -cover -v ./...  -coverprofile=c.out
	@go tool cover -html=c.out
	@rm c.out

build_client: proto_gen_go build test

run_client:
	@mkdir -p "../latency-samples"
	@./client -o "../latency-samples" -g "../endpoints" -c "../experiments/tests/aws-test.json"

build_and_run_client: build_client run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

workflows_edit:
	@gedit ../../.github/workflows/pipeline.yml

replace_github_commit:
	@echo "WARNING: This will latest remote commit with a new one having the newest changes, risks are involved."
	@git add ../../.
	@git commit -S --no-edit --amend
	@git push --force

# Cloudlab deployment
remote="${CLOUDLAB_MACHINE}"

send_to_cloudlab:
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@scp ./client $(remote):~/
	@scp ./vHive-API-template-prod-oas30-apigateway.json $(remote):~/
	@scp ../initialize_machine.sh $(remote):~/
	@scp ../endpoints/* $(remote):~/ # optional
	@ssh $(remote) "mkdir -p ~/experiments/"
	@scp -r ../experiments/* $(remote):~/experiments/
	@ssh $(remote) "mkdir -p ~/setup/deployment/"
	@scp -rp ./setup/deployment/raw-code/ $(remote):~/setup/deployment/

connect_to_cloudlab:
	@ssh -X -p 22 $(remote)

download_cloudlab_results:
	@mkdir -p "../latency-samples/cloudlab"
	@scp -r $(remote):~/latency-samples/* "../latency-samples/cloudlab/"

.PHONY: proto_gen_go build test build_client run_client build_and_run_client workflows_edit replace_github_commit send_to_cloudlab connect_to_cloudlab download_cloudlab_results
