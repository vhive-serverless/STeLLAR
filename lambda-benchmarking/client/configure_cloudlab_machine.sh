#!/bin/bash
sudo apt-get update
sudo apt-get install tmux
sudo apt-get install awscli
mkdir -p "latency-samples"
ulimit -n 16384
aws configure

# Export the two necessary environment variables
read -r -s -p 'Please enter your AWS_LAMBDA_ROLE: ' AWS_LAMBDA_ROLE
echo ""
# Below has to run as root
echo "export AWS_LAMBDA_ROLE=${AWS_LAMBDA_ROLE}" >>/etc/profile.d/benchmarking.sh
unset AWS_LAMBDA_ROLE

read -r -s -p 'Please enter your AWS_API_GATEWAY_KEY: ' AWS_API_GATEWAY_KEY
echo ""
# Below has to run as root
echo "export AWS_LAMBDA_ROLE=${AWS_API_GATEWAY_KEY}" >>/etc/profile.d/benchmarking.sh
unset AWS_API_GATEWAY_KEY

# Run this afterwards in current shell to get the new environment variables
# source /etc/profile.d/benchmarking.sh
