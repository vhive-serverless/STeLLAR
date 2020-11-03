build:
	@GOOS=linux go build -v -a
	@golangci-lint run ./...

run_client:
	@mkdir -p "../latency-samples"
	@./client \
		--outputPath "../latency-samples" \
		--configPath "../config.csv"\
		--gatewaysPath "../gateways.csv"

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

# Cloudlab deployment
remote=theodor@hp104.utah.cloudlab.us

send_code_to_cloudlab:
	@scp ../initialize_machine.sh $(remote):~/
	@scp ../gateways.csv $(remote):~/
	@scp ../config.csv $(remote):~/
	@scp ./client $(remote):~/

send_timer_to_cloudlab:
	@GOOS=linux go build -v -a -o ../timer ../timer.go
	@scp ../timer $(remote):~/

connect_to_cloudlab:
	@ssh -X -p 22 $(remote)

download_cloudlab_results:
	@mkdir -p "../latency-samples/cloudlab"
	@scp -r $(remote):~/latency-samples/* "../latency-samples/cloudlab/"

.PHONY: build run_client build_and_run_client edit_local_env_vars send_code_to_cloudlab send_timer_to_cloudlab connect_to_cloudlab download_cloudlab_results
