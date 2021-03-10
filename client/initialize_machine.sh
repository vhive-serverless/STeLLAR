#!/bin/bash
sudo apt-get update

# For the client to run unattended
sudo apt-get install tmux

# Equivalent to "ulimit -n 1000000", see https://superuser.com/questions/1289345/why-doesnt-ulimit-n-work-when-called-inside-a-script
sudo sh -c "echo \"* soft nofile 1000000\" >> /etc/security/limits.conf"
sudo sh -c "echo \"* hard nofile 1000000\" >> /etc/security/limits.conf"
sudo sh -c "echo \"root soft nofile 1000000\" >> /etc/security/limits.conf"
sudo sh -c "echo \"root hard nofile 1000000\" >> /etc/security/limits.conf"

sudo mkdir latency-samples
echo "tmux new -s vhive-bench"
echo "AWS example run: sudo ./main -o latency-samples -g endpoints -c experiments/tests/aws/data-transfer.json"
