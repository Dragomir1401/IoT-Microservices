#!/bin/bash

echo "Building Docker images..."
docker build -t mqtt-adaptor .

echo "Deploying stack..."
docker stack deploy -c stack.yml my-stack
