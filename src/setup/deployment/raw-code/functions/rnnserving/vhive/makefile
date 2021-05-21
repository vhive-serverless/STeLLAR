.PHONY: build_container_and_run_server
build_container_and_run_server:
	@docker build --no-cache --tag rnnserving-vhive:latest .
	@docker run -it -p 50051:50051 rnnserving-vhive

.PHONY: build_and_push
build_and_push:
	@docker build --no-cache --tag vhive-bench:rnnserving .
	@docker tag vhive-bench:rnnserving vhiveease/vhive-bench:rnnserving
	@docker push vhiveease/vhive-bench:rnnserving
