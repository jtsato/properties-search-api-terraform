#!/bin/bash

set -e

# Update system packages
sudo apt update -y

# Install necessary packages to configure Docker repository
sudo apt install -y apt-transport-https lsb-release ca-certificates curl gnupg

# Download and decrypt Docker's GPG key to authenticate packages
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg -f --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Add Docker repository entry to sources list file if not already added
if ! grep -q "docker.com" /etc/apt/sources.list.d/docker.list; then
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
fi

# Update system packages again after adding Docker repository
sudo apt update -y

# Install Docker Engine, Docker CLI, and Containerd
sudo apt install -y docker-ce docker-ce-cli containerd.io

# Check installed Docker version
sudo docker version

# Add current user to docker group to manage Docker containers without sudo
sudo usermod -aG docker ${USER}

# Install Docker Compose
sudo apt install -y docker-compose

# Restart Docker service to apply configuration changes
sudo systemctl restart docker

# Check installed Docker Compose version
sudo docker-compose version
