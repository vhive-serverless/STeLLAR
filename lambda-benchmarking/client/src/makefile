build:
	@GOOS=linux go build -race -v -a

run_client:
	@mkdir -p "../latency-samples"
	@./client \
		--outputPath "../latency-samples" \
		--configPath "../config.csv"\
		--gatewaysPath "../gateways.csv" \
		--runExperiment 2

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

# Cloudlab helpers
remote=hp151.utah.cloudlab.us
send_code_to_cloudlab:
	@scp ../configure_cloudlab_machine.sh theodor@$(remote):~/
	@scp ../gateways.csv theodor@$(remote):~/
	@scp ../config.csv theodor@$(remote):~/
	@scp ./client theodor@$(remote):~/
	@scp ./makefile theodor@$(remote):~/

connect_to_cloudlab:
	@ssh -X -p 22 theodor@$(remote)

send_and_connect_to_cloudlab: send_code_to_cloudlab connect_to_cloudlab

download_cloudlab_results:
	@mkdir -p "../cloudlab-latency-samples"
	@scp -r theodor@$(remote):~/latency-samples/ "../cloudlab-latency-samples/"

.PHONY: build run_client build_and_run_client edit_local_env_vars send_code_to_cloudlab connect_to_cloudlab send_and_connect_to_cloudlab get_cloudlab_results
