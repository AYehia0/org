#!/bin/sh

# is script is meant to run on production server only
echo "Running Gin Server on production"
exec "$@"
