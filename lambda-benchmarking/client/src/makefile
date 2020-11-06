build:
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@golangci-lint run ./...

run_client:
	@mkdir -p "../latency-samples"
	@./client \
		-o "../latency-samples" \
		-c "../experiments/test.json"\
		-g "../gateways.csv"

build_and_run_client: build run_client

workflows_edit:
	@gedit ../../../.github/workflows/go.yml

replace_github_commit:
	@echo "WARNING: This will latest remote commit with a new one having the newest changes, risks are involved."
	@git add ../../../.
	@git commit --no-edit --amend
	@git push --force

edit_local_env_vars:
	@sudo gedit /etc/profile.d/benchmarking.sh

# Cloudlab deployment
remote=theodor@hp084.utah.cloudlab.us

send_to_cloudlab:
	@scp ../initialize_machine.sh $(remote):~/
	@scp ../gateways.csv $(remote):~/
	@scp ../experiments/* $(remote):~/
	@CGO_ENABLED=0 GOOS=linux go build -v -a
	@scp ./client $(remote):~/
	@CGO_ENABLED=0 GOOS=linux go build -v -a -o ../timer ../timer.go
	@scp ../timer $(remote):~/

connect_to_cloudlab:
	@ssh -X -p 22 $(remote)

download_cloudlab_results:
	@mkdir -p "../latency-samples/cloudlab"
	@scp -r $(remote):~/latency-samples/* "../latency-samples/cloudlab/"

.PHONY: build run_client build_and_run_client edit_local_env_vars send_code_to_cloudlab connect_to_cloudlab download_cloudlab_results
