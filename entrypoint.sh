#!/bin/sh

# Install dependencies if needed
go mod tidy

# Fix permissions for mounted volumes (useful if running on Linux)
chmod -R 777 /app/tmp || true

# Start air
exec air