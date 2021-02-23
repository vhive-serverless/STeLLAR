# MIT License
#
# Copyright (c) 2021 Theodor Amariucai
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

.PHONY: setup_machine
setup_machine:
	@sudo bash initialize_machine.sh

.PHONY: build_and_run
build_and_run: build run_aws

.PHONY: build
build:
	@docker build -t vhive-bench .
	@docker create --name temp vhive-bench /bin/bash
	@docker cp temp:/app/client ./main
	@docker cp temp:/app/setup/deployment/raw-code/producer-consumer/go1.x/aws/handler setup/deployment/raw-code/producer-consumer/go1.x/aws/handler
	@docker cp temp:/app/setup/deployment/raw-code/producer-consumer/go1.x/vhive/handler setup/deployment/raw-code/producer-consumer/go1.x/vhive/handler
	@docker container rm temp

.PHONY: run_vhive
run_vhive:
	@./main -o ../latency-samples -g ../endpoints -c ../experiments/tests/vhive/burstiness-image-package.json

.PHONY: run_aws_S3
run_aws_S3:
	@./main -o ../latency-samples -g ../endpoints -c ../experiments/tests/AWS/data-transfer-S3.json

.PHONY: empty_S3_bucket
empty_S3_bucket:
	@aws s3 rm --recursive s3://vhive-bench/

.PHONY: run_aws
run_aws:
	@./main -o ../latency-samples -g ../endpoints -c ../experiments/tests/AWS/data-transfer.json

.PHONY: test
test:
	@go test -short -v ./...

.PHONY: start_docker_daemon
start_docker_daemon:
	@sudo systemctl start docker

.PHONY: edit_local_env_vars
edit_local_env_vars:
	@kate /etc/profile.d/benchmarking.sh

.PHONY: workflows_edit
workflows_edit:
	@kate ../../.github/workflows/pipeline.yml

.PHONY: replace_github_commit
replace_github_commit: ## This will replace latest commit with a new one having the newest changes, risks are involved.
	@echo "WARNING: This will replace latest commit with a new one having the newest changes, risks are involved."
	@git add ../../.
	@git commit -S --no-edit --amend --no-gpg-sign
	@git push --force

.PHONY: vhive_local_function_deployment
vhive_local_docker_deployment:
	@docker system prune
	@docker build -t vhive-bench:latest setup/deployment/raw-code/producer-consumer/go1.x/vhive/
	@docker run -dp 50051:50051 vhive-function:latest

# Cloudlab deployment
.PHONY: send_to_cloudlab
send_to_cloudlab:
	@scp -rp ./setup/deployment/raw-code/ ${CLOUDLAB_REMOTE}:~/vhive-bench-deployment/setup/deployment/

.PHONY: connect_to_cloudlab
connect_to_cloudlab:
	@ssh -X -p 22 ${CLOUDLAB_REMOTE}

.PHONY: download_cloudlab_results
download_cloudlab_results:
	@mkdir -p "../latency-samples/cloudlab"
	@scp -r ${CLOUDLAB_REMOTE}:~/vhive-bench/latency-samples/* "../latency-samples/cloudlab/"

.PHONY: help
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
