download_aws_lambda_library:
	@go get github.com/aws/aws-lambda-go/lambda

deploy_to_aws:
	@GOOS=linux go build -o ./producer-handler src/producer/handler.go
	@zip benchmarking.zip producer-handler
	@/usr/local/bin/aws lambda create-function \
		--function-name benchmarking \
		--runtime go1.x \
       	--role $$AWS_LAMBDA_ROLE \
       	--handler producer-handler \
       	--zip-file fileb://benchmarking.zip \
       	--tracing-config Mode=Active > deployments/create_$$( date '+%F_%H:%M:%S' ).txt
       	# Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray. PassThrough otherwise.

update_aws_deployment:
	@GOOS=linux go build -o ./producer-handler src/producer/handler.go
	@zip benchmarking.zip producer-handler
	@/usr/local/bin/aws lambda update-function-code \
		--function-name benchmarking \
       	--zip-file fileb://benchmarking.zip > deployments/update_$$( date '+%F_%H:%M:%S' ).txt

remove_from_aws:
	@/usr/local/bin/aws lambda delete-function \
		--function-name benchmarking \

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: download_aws_lambda_library deploy_to_aws update_aws_deployment remove_from_aws edit_local_env_vars