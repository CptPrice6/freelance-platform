#!/bin/bash

# Prerequisites: 
# sudo apt update
# sudo apt install git, clone repository ( git clone https://github.com/CptPrice6/freelance-platform.git )
# Add backend external VM port that forwards to 8080 internally to .env in frontend

# Update package list and install necessary tools
sudo apt update

# Install Docker prerequisites
sudo apt install apt-transport-https ca-certificates curl software-properties-common lsb-release -y

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up the stable Docker repository for Debian
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Update apt package index and install Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Build and run Docker containers
sudo docker compose -p go-website up --build -d

# Output helpful information on how to access the project
echo "Docker containers are now up and running!"
echo "To access the project:"
echo "1. If you're running this on your local machine, open your browser and go to http://localhost:3000."
echo "2. If running on a remote server, use the server's public IP address, e.g., http://<your-server-ip>:3000."
