#!/bin/bash

# Define the directory where logs should be stored
LOG_DIR="/home/jed/yeah-buddy-lightweight/logs"

# Get the current date in YYYY-MM-DD format
TODAY=$(date +"%Y-%m-%d")

# Use docker-compose to save logs. Adjust the path to your docker-compose.yml if necessary
cd /home/jed/yeah-buddy-lightweight/logs && sudo docker-compose logs > "${LOG_DIR}/log-${TODAY}.txt" 2>&1

# Delete log files older than 2 weeks
find "${LOG_DIR}" -name 'log-*.txt' -type f -mtime +14 -exec rm -f {} \;

