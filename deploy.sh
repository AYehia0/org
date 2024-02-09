#!/bin/bash

echo "Fetching latest code from github"
git pull origin production

echo "Stopping Docker Compose"
docker-compose -f docker-compose.yaml down

echo "Running Docker Compose"
docker-compose -f docker-compose.yaml up -d --build
