.PHONY: build_container_and_run_server
build_container_and_run_server:
	@docker build --no-cache --tag chameleon-vhive:latest .
	@docker run -it -p 50051:50051 chameleon-vhive

.PHONY: build_and_push
build_and_push:
	@docker build --no-cache --tag vhive-bench:chameleon .
	@docker tag vhive-bench:chameleon vhiveease/vhive-bench:chameleon
	@docker push vhiveease/vhive-bench:chameleon
