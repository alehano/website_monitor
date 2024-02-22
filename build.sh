#!/bin/bash

# Set the target OS and architecture
# GOOS=linux GOARCH=arm go build -o website_monitor
GOOS=linux GOARCH=amd go build -o website_monitor

echo "Build complete."