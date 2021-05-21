.PHONY: build_container_and_run_server
build_container_and_run_server:
	@docker build --no-cache --tag hellopy-aws:latest .
	@docker run -it -p 50051:50051 hellopy-aws

.PHONY: build_and_push_to_ecr
build_and_push_to_ecr:
	@aws ecr get-login --region us-west-1
	@docker build --no-cache -t vhive-bench:hellopy .
	@docker tag vhive-bench:hellopy 356764711652.dkr.ecr.us-west-1.amazonaws.com/vhive-bench:hellopy
	@docker push 356764711652.dkr.ecr.us-west-1.amazonaws.com/vhive-bench:hellopy
