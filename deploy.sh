#!/bin/bash

echo "Fetching latest code from github"
git pull origin production

echo "Stopping Docker Compose"
docker-compose -f docker-compose.yml down

echo "Running Docker Compose"
docker-compose -f docker-compose.yml up -d --build
