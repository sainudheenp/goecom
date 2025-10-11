#!/bin/bash

# Migration script for running SQL migrations

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR/.."

# Load environment variables from .env file if it exists
if [ -f "$PROJECT_ROOT/.env" ]; then
    echo "Loading environment from .env file..."
    export $(grep -v '^#' "$PROJECT_ROOT/.env" | xargs)
fi

# Get the direction (up or down)
DIRECTION=${1:-up}

# Database connection from environment
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is not set"
    echo "Make sure you have a .env file with DATABASE_URL in the project root"
    exit 1
fi

echo "Running migrations: $DIRECTION"

MIGRATIONS_DIR="$PROJECT_ROOT/migrations"

# Function to run migrations
run_migrations() {
    local direction=$1
    
    # Find all migration files for the given direction
    for file in "$MIGRATIONS_DIR"/*."$direction".sql; do
        if [ -f "$file" ]; then
            echo "Running migration: $(basename "$file")"
            psql "$DATABASE_URL" -f "$file"
        fi
    done
}

# Run migrations
if [ "$DIRECTION" = "up" ]; then
    run_migrations "up"
    echo "Migrations completed successfully"
elif [ "$DIRECTION" = "down" ]; then
    # Reverse order for down migrations
    for file in $(ls -r "$MIGRATIONS_DIR"/*.down.sql 2>/dev/null); do
        echo "Running migration: $(basename "$file")"
        psql "$DATABASE_URL" -f "$file"
    done
    echo "Rollback completed successfully"
else
    echo "Error: Invalid direction. Use 'up' or 'down'"
    exit 1
fi
