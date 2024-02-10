#!/bin/sh
echo "Starting the server"

echo "Running Server with Gin for hot reload"

export ENV=development
export APORT=8080

CompileDaemon --build="go build -o main cmd/main.go" --command=./main
