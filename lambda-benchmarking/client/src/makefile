LATENCY_SAMPLES_OUT := "../latency-samples"

build:
	@GOOS=linux go build -race -v -a

# 500 requests max to avoid `race: limit on 8128 simultaneously alive goroutines is exceeded, dying`
run_client:
	@mkdir -p $(LATENCY_SAMPLES_OUT)
	@./client \
		--payloadLengthBytes 3 \
		--outputPath $(LATENCY_SAMPLES_OUT) \
		--requests 15 \
		--bursts 2 \
		--frequencySeconds 1 \
		--randomized=true \
		--visualization "CDF" \
		--gatewayEndpoint https://xw5co672m6.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://2szwlepb1j.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://q5boxdpf2d.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://hydwazqc6j.execute-api.eu-west-2.amazonaws.com/default
#		--gatewayEndpoint https://zcnmx4h9da.execute-api.eu-west-2.amazonaws.com/default

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: build run_client build_and_run_client edit_local_env_vars