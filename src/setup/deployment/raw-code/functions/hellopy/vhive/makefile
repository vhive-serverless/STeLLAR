.PHONY: build_container_and_run_server
build_container_and_run_server:
	@docker build --no-cache --tag hellopy-vhive:latest .
	@docker run -it -p 50051:50051 hellopy-vhive

.PHONY: build_and_push
build_and_push:
	@docker build --no-cache --tag vhive-bench:hellopy .
	@docker tag vhive-bench:hellopy vhiveease/vhive-bench:hellopy
	@docker push vhiveease/vhive-bench:hellopy
