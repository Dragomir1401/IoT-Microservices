#!/bin/bash

# Remove the existing stack if any
echo "Removing existing stack (if any)..."
docker stack rm scd3 2>/dev/null || echo "No stack to remove."
sleep 5

# Deploy the full stack
echo "Deploying the full stack..."
docker stack deploy --compose-file stack.yml scd3

echo "Waiting for services to start..."
sleep 10

echo "Services have been successfully started!"
