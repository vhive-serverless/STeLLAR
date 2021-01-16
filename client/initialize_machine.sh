#!/bin/bash
sudo apt-get update

# For the client to run unattended
sudo apt-get install tmux

# For the client to build binaries and deploy on-the-fly
sudo apt-get install golang-go
sudo go get github.com/aws/aws-lambda-go/events
sudo go get github.com/aws/aws-lambda-go/lambda
sudo go get github.com/aws/aws-lambda-go/lambdacontext

# For the client to deploy container images
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io

sudo apt-get install awscli
aws configure

mkdir -p "latency-samples"
echo "Please now run 'ulimit -n 16384' in this shell before running the client."
echo "Suggested operation: tmux new -s cloudlab"
echo "Example client run: ./client -c experiments/tests/test.json"
