#!/bin/sh

if ! [ -x "$(command -v psql)" ]; then
  echo 'Error: psql is not installed.\nSee https://www.postgresql.org.' >&2
  exit 1
fi

# Load env variables
export $(cat ./dev/local.env) 

# Create .db folders to keep data persistant
mkdir -p .db

# Stop already running postgres container
docker stop postgres >/dev/null 2>&1
docker rm postgres >/dev/null 2>&1

# Create network if doens't exists
docker network rm ebank >/dev/null 2>&1
docker network create -d bridge ebank >/dev/null 2>&1

# Start postgres container
docker run -e POSTGRES_USER=${DB_USER} \
  -e POSTGRES_PASSWORD=${DB_PASS} \
  -e POSTGRES_DB=${DB_NAME} \
  -p ${DB_PORT}:${DB_PORT} \
  -v "${PWD}/.db:/var/lib/postgresql/data/" \
  --network ebank \
  -d --name postgres postgres:alpine

# Wait until postgres is ready
CMD="host=localhost port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASS}"
# Sleep up to 1 minute until postgres has started
RETRIES=60
until psql "$CMD" -c "select 1" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "Waiting for postgres..."
  sleep 1
done

echo
echo "Postgres started successfully on port ${DB_PORT}"
echo