#!/bin/bash

# Define the directory where logs should be stored
LOG_DIR="/home/jed/yeah-buddy-lightweight/logs"

# Ensure the logs directory exists
mkdir -p "${LOG_DIR}"

# Get the current date in YYYY-MM-DD format
TODAY=$(date +"%Y-%m-%d")

# Path to the current log file
LOG_FILE="${LOG_DIR}/log-${TODAY}.txt"

# Use docker-compose to save logs in real-time
# Adjust the path to your docker-compose.yml if necessary
cd /home/jed/yeah-buddy-lightweight && sudo docker-compose logs --follow > "${LOG_FILE}" 2>&1 &
PID=$!

# Function to clean up when exiting
function finish {
  # Kill the background docker-compose logs process
  kill $PID
  # Delete log files older than 2 weeks
  find "${LOG_DIR}" -name 'log-*.txt' -type f -mtime +14 -exec rm -f {} \;
}
trap finish EXIT

# Wait for the background process to end (it won't, unless manually stopped)
wait $PID

