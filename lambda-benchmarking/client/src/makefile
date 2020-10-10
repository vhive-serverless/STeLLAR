build:
	@GOOS=linux go build -race -v -a

run_client:
	@mkdir -p "../latency-samples"
	@./client \
		--outputPath "../latency-samples" \
		--configPath "../config.csv"

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

# Cloudlab helpers
send_code_to_cloudlab:
	@scp ../configure_cloudlab_machine.sh theodor@pc05.cloudlab.umass.edu:~/
	@scp ../config.csv theodor@pc05.cloudlab.umass.edu:~/
	@scp ./client theodor@pc05.cloudlab.umass.edu:~/
	@scp ./makefile theodor@pc05.cloudlab.umass.edu:~/

connect_to_cloudlab:
	@ssh -X -p 22 theodor@pc05.cloudlab.umass.edu

send_and_connect_to_cloudlab: send_code_to_cloudlab connect_to_cloudlab

download_cloudlab_results:
	@mkdir -p "../cloudlab-latency-samples"
	@scp -r theodor@pc05.cloudlab.umass.edu:~/latency-samples/ "../cloudlab-latency-samples/"

.PHONY: build run_client build_and_run_client edit_local_env_vars send_code_to_cloudlab connect_to_cloudlab send_and_connect_to_cloudlab get_cloudlab_results
