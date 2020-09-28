download_aws_lambda_library:
	@go get github.com/aws/aws-lambda-go/lambda

compile_executable:
	@GOOS=linux go build -o ./producer-handler src/producer/handler.go

create_deployment_package:
	@zip benchmarking.zip producer-handler

deploy_to_aws:
	@/usr/local/bin/aws lambda create-function \
		--function-name benchmarking \
		--runtime go1.x \
       	--role $$AWS_LAMBDA_ROLE \
       	--handler producer-handler \
       	--zip-file fileb://benchmarking.zip \
       	--tracing-config Mode=Active > deployments/$$( date '+%F_%H:%M:%S' ).txt
       	# Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray. PassThrough otherwise.

remove_from_aws:
	@/usr/local/bin/aws lambda delete-function \
		--function-name benchmarking \

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

all: compile_executable create_deployment_package deploy_to_aws

.PHONY: download_aws_lambda_library compile_executable create_deployment_package deploy_to_aws edit_local_env_vars all