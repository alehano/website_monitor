#!/bin/bash

# Load environment variables from .env file
export $(grep -v '^#' .env | xargs)

# Run website_monitor in the background and save its PID
./website_monitor &
echo $! > website_monitor.pid