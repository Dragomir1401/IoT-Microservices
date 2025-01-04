#!/bin/bash

# Define the token file path
TOKEN_FILE="./influxdb_token.txt"

# Set default value for SCD_DVP if not provided
export SCD_DVP=${SCD_DVP:-/var/lib/docker/volumes}
echo "Using Docker volume path: $SCD_DVP"

# Ensure directories for volumes exist
sudo mkdir -p "$SCD_DVP/influxdb-data"
sudo mkdir -p "$SCD_DVP/grafana-data"

# Check if InfluxDB is already initialized
if [ ! -f "$SCD_DVP/influxdb-data/influxd.bolt" ]; then
  echo "InfluxDB not initialized. Cleaning up and setting up a fresh configuration..."
  rm -rf "$SCD_DVP/influxdb-data/*" # Clear any leftover data
  sudo mkdir -p "$SCD_DVP/influxdb-data/engine"
  sudo chown -R 1000:1000 "$SCD_DVP/influxdb-data"
fi

# Check if Grafana is already initialized
if [ ! -d "$SCD_DVP/grafana-data" ]; then
  echo "Grafana not initialized. Setting up a fresh configuration..."
  rm -rf "$SCD_DVP/grafana-data" # Clear any leftover data
  sudo mkdir -p "$SCD_DVP/grafana-data"
  sudo chown -R 472:472 "$SCD_DVP/grafana-data"
fi

# Check if the token file exists
if [ ! -f "$TOKEN_FILE" ]; then
  echo "Error: Token file $TOKEN_FILE does not exist. Please ensure it contains the InfluxDB token."
  exit 1
fi

# Extract the token value
TOKEN=$(cat "$TOKEN_FILE")
if [ -z "$TOKEN" ]; then
  echo "Error: Token file is empty or invalid. Please provide a valid InfluxDB token in $TOKEN_FILE."
  exit 1
fi

# Remove the existing stack if any
echo "Removing existing stack (if any)..."
docker stack rm scd3 2>/dev/null || echo "No stack to remove."
sleep 5  # Allow time for cleanup

# Deploy the full stack with the correct token
echo "Deploying the full stack with token..."
export INFLUXDB_TOKEN="$TOKEN"
docker stack deploy --compose-file stack.yml scd3

echo "Waiting for services to start..."
sleep 5

echo "Services have been successfully started!"
