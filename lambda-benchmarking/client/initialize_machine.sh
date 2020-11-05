#!/bin/bash
sudo apt-get update
sudo apt-get install tmux
sudo apt-get install awscli
mkdir -p "latency-samples"
aws configure
echo "Please now run 'ulimit -n 16384' in this shell before running the client."
echo "Please now run 'export AWS_ACCESS_KEY_ID=*******' in this shell before running the client."
echo "Please now run 'export AWS_SECRET_ACCESS_KEY=*******' in this shell before running the client."
