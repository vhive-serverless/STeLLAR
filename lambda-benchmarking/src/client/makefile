LATENCY_SAMPLES_OUT := "latency-samples"

build:
	@GOOS=linux go build -race -v -a client.go benchmarkWriter.go

run_client:
	@mkdir -p $(LATENCY_SAMPLES_OUT)
	@./client --requests 50 -execMilliseconds 40 -outputPath $(LATENCY_SAMPLES_OUT)

build_and_run_client: build run_client

.PHONY: build run_client build_and_run_client