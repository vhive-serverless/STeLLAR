#!/bin/bash
sudo apt-get update
sudo apt-get install tmux
sudo apt-get install awscli
mkdir -p "latency-samples"
aws configure
echo "Please now run 'ulimit -n 16384' in this shell before running the client."
echo "Suggested operation: tmux new -s cloudlab"
echo "Example client run: ./client -o latency-samples/ -c experiments/test.json -g gateways.csv"
