#!/bin/bash

## Define the token file path
#TOKEN_FILE="./influxdb_token.txt"
#
## Check if the token file exists
#if [ ! -f "$TOKEN_FILE" ]; then
#  echo "Error: Token file $TOKEN_FILE does not exist. Please ensure it contains the InfluxDB token."
#  exit 1
#fi
#
## Extract the token value
#TOKEN=$(cat "$TOKEN_FILE")
#if [ -z "$TOKEN" ]; then
#  echo "Error: Token file is empty or invalid. Please provide a valid InfluxDB token in $TOKEN_FILE."
#  exit 1
#fi

# Remove the existing stack if any
echo "Removing existing stack (if any)..."
docker stack rm scd3 2>/dev/null || echo "No stack to remove."
sleep 5  # Allow time for cleanup

# Deploy the full stack with the correct token
#echo "Deploying the full stack with token..."
#export INFLUXDB_TOKEN="$TOKEN"
echo "Deploying the full stack..."
docker stack deploy --compose-file stack.yml scd3

echo "Waiting for services to start..."
sleep 10

echo "Services have been successfully started!"
