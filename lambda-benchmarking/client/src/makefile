build:
	@GOOS=linux go build -race -v -a

run_client:
	@mkdir -p "../latency-samples"
	@./client \
		--outputPath "../latency-samples" \
		--configPath "../config.csv"\
		--gatewaysPath "../gateways.csv"

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

# Cloudlab helpers
remote=hp060.utah.cloudlab.us
send_code_to_cloudlab:
	@scp ../configure_cloudlab_machine.sh theodor@$(remote):~/exp2/
	@scp ../gateways.csv theodor@$(remote):~/exp2/
	@scp ../config.csv theodor@$(remote):~/exp2/
	@scp ./client theodor@$(remote):~/exp2/

send_timer_to_cloudlab:
	@GOOS=linux go build -race -v -a -o ../timer ../timer.go
	@scp ../timer theodor@$(remote):~/

connect_to_cloudlab:
	@ssh -X -p 22 theodor@$(remote)

download_cloudlab_results:
	@mkdir -p "../cloudlab-latency-samples"
	@scp -r theodor@$(remote):~/latency-samples/ "../cloudlab-latency-samples/"

.PHONY: build run_client build_and_run_client edit_local_env_vars send_code_to_cloudlab send_timer_to_cloudlab connect_to_cloudlab get_cloudlab_results
