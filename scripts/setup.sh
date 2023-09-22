#!/bin/bash
sudo add-apt-repository ppa:longsleep/golang-backports
sudo add-apt-repository ppa:linuxuprising/java
sudo add-apt-repository ppa:cwchien/gradle
sudo apt-get install apt-transport-https ca-certificates gnupg curl sudo
echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
sudo apt-get update

# Set up Docker Registry
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
"deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
"$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Install Golang
sudo apt install golang

# Install OpenJDK 11
sudo apt-get install openjdk-11-jdk

# Install Gradle
sudo apt install gradle

# Install Node
NODE_MAJOR=16
sudo apt-get install nodejs -y

# Install npm and serverless
sudo apt install npm
sudo npm install -g serverless serverless-azure-functions --save-dev

# For the client to run unattended
sudo apt-get install --no-install-recommends --assume-yes tmux

# Equivalent to "ulimit -n 1000000", see https://superuser.com/questions/1289345/why-doesnt-ulimit-n-work-when-called-inside-a-script
sudo sh -c "echo \"* soft nofile 1000000\" >> /etc/security/limits.conf"
sudo sh -c "echo \"* hard nofile 1000000\" >> /etc/security/limits.conf"
sudo sh -c "echo \"root soft nofile 1000000\" >> /etc/security/limits.conf"
sudo sh -c "echo \"root hard nofile 1000000\" >> /etc/security/limits.conf"

sudo mkdir ../latency-samples
echo "Please run: ulimit -n 8192"
echo "Recommended: tmux new -s stellar"
echo "AWS example run: sudo ./main -o latency-samples -g endpoints -c experiments/tests/aws/data-transfer.json"
