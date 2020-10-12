download_aws_lambda_library:
	@go get github.com/aws/aws-lambda-go/lambda

deploy:
	./deployment_scripts/deploy-producers.sh 25

update:
	./deployment_scripts/update-producers.sh 25

remove:
	./deployment_scripts/remove-producers.sh 25

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

.PHONY: download_aws_lambda_library deploy_to_aws update_aws_deployment remove_from_aws edit_local_env_vars