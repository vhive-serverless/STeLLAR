LATENCY_SAMPLES_OUT := "../latency-samples"

build:
	@GOOS=linux go build -race -v -a

# 500 requests max to avoid `race: limit on 8128 simultaneously alive goroutines is exceeded, dying`
run_client:
	@mkdir -p $(LATENCY_SAMPLES_OUT)
	@./client \
		--payloadLengthBytes 3 \
		--outputPath $(LATENCY_SAMPLES_OUT) \
		--requests 10 \
		--bursts 10 \
		--frequencySeconds 1 \
		--randomized=true \
		--visualization "CDF" \
		--gatewayEndpoint https://xw5co672m6.execute-api.eu-west-2.amazonaws.com/default

# Gateway 1: https://xw5co672m6.execute-api.eu-west-2.amazonaws.com/default
# Gateway 2:
# Gateway 3:
# Gateway 4:
# Gateway 5:

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: build run_client build_and_run_client edit_local_env_vars