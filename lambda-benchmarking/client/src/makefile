LATENCY_SAMPLES_OUT := "../latency-samples"

build:
	@GOOS=linux go build -race -v -a

run_client:
	@mkdir -p $(LATENCY_SAMPLES_OUT)
	@./client \
		--outputPath $(LATENCY_SAMPLES_OUT) \
		--burstSize 1 \
		--frequencySeconds 600 \
		--gatewayEndpoint https://xw5co672m6.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://2szwlepb1j.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://q5boxdpf2d.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://hydwazqc6j.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://zcnmx4h9da.execute-api.eu-west-2.amazonaws.com/default

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

# Cloudlab helpers
send_code_to_cloudlab:
	@scp ./configure_cloudlab_machine.sh theodor@pc05.cloudlab.umass.edu:~/
	@scp ./client theodor@pc05.cloudlab.umass.edu:~/
	@scp ./gateways.csv theodor@pc05.cloudlab.umass.edu:~/
	@scp ./makefile theodor@pc05.cloudlab.umass.edu:~/

connect_to_cloudlab:
	@ssh -X -p 22 theodor@pc05.cloudlab.umass.edu

send_and_connect_to_cloudlab: send_code_to_cloudlab connect_to_cloudlab

download_cloudlab_results:
	@mkdir -p "$(LATENCY_SAMPLES_OUT)/cloudlab"
	@scp -r theodor@pc05.cloudlab.umass.edu:~/latency-samples/ "$(LATENCY_SAMPLES_OUT)/cloudlab/"

.PHONY: build run_client build_and_run_client edit_local_env_vars send_code_to_cloudlab connect_to_cloudlab send_and_connect_to_cloudlab get_cloudlab_results
