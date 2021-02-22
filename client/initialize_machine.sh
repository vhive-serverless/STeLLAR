#!/bin/bash
sudo apt-get update

# For the client to run unattended
sudo apt-get install tmux

# Install Docker - used for automatic container deployment, but conflicts with vhive for now...
#curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
#sudo add-apt-repository \
#   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
#   $(lsb_release -cs) \
#   stable"
#sudo apt-get update
#sudo apt-get install docker-ce docker-ce-cli containerd.io

# Install & configure AWS CLI
echo "Y" | sudo apt-get install awscli
sudo aws configure

# Equivalent to "ulimit -n 1000000", see https://superuser.com/questions/1289345/why-doesnt-ulimit-n-work-when-called-inside-a-script
sudo prlimit --pid=$PPID --nofile=1000000

sudo mkdir latency-samples
echo "tmux new -s cloudlab"
echo "AWS example run: sudo ./main -o latency-samples -g endpoints -c experiments/tests/AWS/data-transfer.json"
