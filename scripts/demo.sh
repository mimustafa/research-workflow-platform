#!/usr/bin/env bash

set -euo pipefail

echo "Starting the complete local platform..."
docker compose up --build -d

echo "Waiting for the Go API..."
until curl --fail --silent http://localhost:8080/health >/dev/null; do
  sleep 2
done

echo "Submitting a research brief..."
curl --fail --silent \
  -H "Content-Type: application/json" \
  -d '{"brief":"Compare approaches to reproducible scientific software."}' \
  http://localhost:8080/api/v1/workflows
printf "\n\n"

echo "Waiting for the Python worker to complete the job..."
sleep 2
curl --fail --silent http://localhost:8080/api/v1/workflows
printf "\n\n"

echo "Dashboard: http://localhost:3000"
echo "RabbitMQ management: http://localhost:15672"
