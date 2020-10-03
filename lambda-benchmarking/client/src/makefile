LATENCY_SAMPLES_OUT := "../latency-samples"

build:
	@GOOS=linux go build -race -v -a

# 500 requests max to avoid `race: limit on 8128 simultaneously alive goroutines is exceeded, dying`
run_client:
	@mkdir -p $(LATENCY_SAMPLES_OUT)
	@./client \
		--requests 50 \
		--execMs 50 \
		--payloadLengthBytes 3 \
		--outputPath $(LATENCY_SAMPLES_OUT) \
		--bursts 13 \
#		--frequencySeconds 20 \

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: build run_client build_and_run_client edit_local_env_vars