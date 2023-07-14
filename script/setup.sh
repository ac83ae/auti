#!/usr/bin/env bash

sudo apt-get update -y

# Install Go
wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
rm go1.20.5.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >>~/.bashrc
source ~/.bashrc

# Install jq
sudo apt-get install -y jq

# Install Python and pip
sudo apt-get install -y python3 python3-pip
sudo apt-get install -y python-is-python3

# Install pyyaml
pip3 install pyyaml

# Install Docker
# Uninstall all conflicting packages
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt-get remove $pkg; done
# Update the apt package index and install packages to allow apt to use a repository over HTTPS
sudo apt-get update -y
sudo apt-get install -y ca-certificates curl gnupg
# Add Docker's official GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
# Add Docker apt repository
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" |
  sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
# Update apt package index, install docker engine, and docker-compose
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-compose
# Verify that the Docker Engine installation is successful by running the hello-world image.
sudo docker run hello-world

# Add current user to docker group
# Create the docker group
sudo groupadd docker
# Add your user to the docker group
sudo usermod -aG docker "$USER"
echo "Please log out and log back in so that your group membership is re-evaluated."
