.PHONY: build_container_and_run_server
build_container_and_run_server:
	@docker build --no-cache --tag chameleon-aws:latest .
	@docker run -it -p 50051:50051 chameleon-aws

.PHONY: build_and_push_to_ecr
build_and_push_to_ecr:
	@aws ecr get-login --region us-west-1
	@docker build --no-cache -t vhive-bench:chameleon .
	@docker tag vhive-bench:chameleon 356764711652.dkr.ecr.us-west-1.amazonaws.com/vhive-bench:chameleon
	@docker push 356764711652.dkr.ecr.us-west-1.amazonaws.com/vhive-bench:chameleon
