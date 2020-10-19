download_aws_lambda_library:
	@go get github.com/aws/aws-lambda-go/lambda

build:
	@GOOS=linux go build -v -race -o producer-handler handler.go
	@rm -f benchmarking.zip
	@zip benchmarking.zip producer-handler

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: download_aws_lambda_library edit_local_env_vars build