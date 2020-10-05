LAMBDA_FUNCTION_NAME := "benchmarking"
DEPLOYMENTS_LOG_OUT := "deployments"

download_aws_lambda_library:
	@go get github.com/aws/aws-lambda-go/lambda

deploy_to_aws:
	@GOOS=linux go build -v -race -o ./producer-handler handler.go
	@zip $(LAMBDA_FUNCTION_NAME).zip producer-handler
	@mkdir -p $(DEPLOYMENTS_LOG_OUT)
	@/usr/local/bin/aws lambda create-function \
		--function-name $(LAMBDA_FUNCTION_NAME) \
		--runtime go1.x \
       	--role $$AWS_LAMBDA_ROLE \
       	--handler producer-handler \
       	--zip-file fileb://$(LAMBDA_FUNCTION_NAME).zip \
       	--tracing-config Mode=PassThrough > $(DEPLOYMENTS_LOG_OUT)/created_$$( date '+%F_%H:%M:%S' ).txt
       	# Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray. PassThrough otherwise.

update_aws_deployment:
	@GOOS=linux go build -o ./producer-handler handler.go
	@zip $(LAMBDA_FUNCTION_NAME).zip producer-handler
	@mkdir -p $(DEPLOYMENTS_LOG_OUT)
	@/usr/local/bin/aws lambda update-function-code \
		--function-name $(LAMBDA_FUNCTION_NAME) \
       	--zip-file fileb://$(LAMBDA_FUNCTION_NAME).zip > $(DEPLOYMENTS_LOG_OUT)/updated_$$( date '+%F_%H:%M:%S' ).txt

remove_from_aws:
	@/usr/local/bin/aws lambda delete-function \
		--function-name $(LAMBDA_FUNCTION_NAME) \

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: download_aws_lambda_library deploy_to_aws update_aws_deployment remove_from_aws edit_local_env_vars