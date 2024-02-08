#!/bin/sh
echo "Starting the server"

echo "Running Server with Gin for hot reload"
/usr/local/bin/gin --port 8080 --appPort 8080 run cmd/main.go

# exec "$@"
