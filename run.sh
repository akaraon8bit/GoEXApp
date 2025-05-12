#!/bin/bash
# run.sh

# Load environment variables
set -a
source .env
set +a

# Build the application
echo "Building the application..."
go build -o monolith ./cmd/mallbots

# Initialize the database
echo "Initializing database..."
./init_db.sh

# Run the application
echo "Starting MallBots application..."
./monolith
