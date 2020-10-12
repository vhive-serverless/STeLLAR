#!/bin/bash
sudo apt-get update
sudo apt-get install tmux
sudo apt-get install awscli
mkdir -p "latency-samples"
ulimit -n 4096
aws configure
