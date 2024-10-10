#!/bin/sh

go mod tidy
go mod vendor

if ! [ -x "$(command -v fswatch)" ]; then
  echo 'Error: fresh is not installed.\nSee https://github.com/emcrisostomo/fswatch.' >&2
  exit 1
fi

if ! [ -x "$(command -v protoc)" ]; then
  echo 'Error: protoc is not installed.\nSee http://google.github.io/proto-lens/installing-protoc.html.' >&2
  exit 1
fi

# Create network if doens't exists
docker network create -d bridge ebank &> /dev/null 

# Automatically generate protos on .proto file change
fswatch ./services/*/*.proto | xargs -n1 -I "{}" yarn proto &
P1=$!

docker-compose -f dev/docker-compose.yml down
docker-compose -f dev/docker-compose.yml build
docker-compose -f dev/docker-compose.yml up --remove-orphans
P2=$!

wait $P2

while [ -e /proc/$P1 ]
do
    echo "Process $P1 is still running"
    sleep .6
done
echo "Process $P1 has finished"