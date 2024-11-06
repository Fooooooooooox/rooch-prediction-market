#!/bin/bash

cd /home/ubuntu/backend

# Pull the latest changes
git reset --hard origin/main
git pull origin main

ls -la

source .env

# Build your Go application (if necessary)
go build -o myapp .

# Stop the currently running server (if applicable)
pkill myapp || true

# Start the server
nohup ./myapp > app.log 2>&1 &
