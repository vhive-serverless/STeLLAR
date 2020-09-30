LATENCY_SAMPLES_OUT := "../latency-samples"

build:
	@GOOS=linux go build -race -v -a

run_client:
	@mkdir -p $(LATENCY_SAMPLES_OUT)
	@./client \
		--requests 10 \
		--execMilliseconds 50 \
		--payloadLengthBytes 5 \
		--outputPath $(LATENCY_SAMPLES_OUT)

build_and_run_client: build run_client

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: build run_client build_and_run_client edit_local_env_vars