#!/bin/bash

# Read the PID from the file
PID=$(cat website_monitor.pid)

# Send a SIGTERM signal to the process
kill -SIGTERM $PID

# Remove the PID file
rm website_monitor.pid