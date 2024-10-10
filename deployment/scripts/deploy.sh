#!/bin/bash

PROJECT_ID="ebank-299413"

for dockerfile in deployment/dockerfiles/*.Dockerfile; do
    name=$(basename "$dockerfile" .Dockerfile)
    echo "Building ${name}..."

    if docker build -f ${dockerfile} -t gcr.io/${PROJECT_ID}/${name} . ; then
      echo "Pushing $name..."
      docker push gcr.io/${PROJECT_ID}/${name}
    fi

done 