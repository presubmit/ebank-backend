#!/bin/sh

# Load env variables
export $(cat dev/local.env) 

# Create .redis folder to keep data persistant
mkdir -p .redis

# Stop already running redis container
docker stop redis >/dev/null 2>&1
docker rm redis >/dev/null 2>&1

# Create network if doesn't exist
docker network rm ebank >/dev/null 2>&1
docker network create -d bridge ebank >/dev/null 2>&1

# Start redis container
docker run \
  -p ${REDIS_PORT}:${REDIS_PORT} \
  -v "${PWD}/.redis:/data" \
  --network ebank \
  -d --name redis redis:alpine

echo
echo "Redis started successfully on port ${REDIS_PORT}"
echo