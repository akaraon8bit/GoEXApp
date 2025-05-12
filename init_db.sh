#!/bin/bash
# init_db.sh

# Install psql if not already installed
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL client (psql) not found. Please install it first."
    exit 1
fi

# Connection string from .env
PG_CONN=$(grep PG_CONN .env | cut -d '=' -f2)

# Execute all schema creation scripts in order
echo "Initializing database schemas..."
psql "$PG_CONN" -f customers/internal/postgres/schema.sql
psql "$PG_CONN" -f baskets/internal/postgres/schema.sql
psql "$PG_CONN" -f depot/internal/postgres/schema.sql
psql "$PG_CONN" -f ordering/internal/postgres/schema.sql
psql "$PG_CONN" -f payments/internal/postgres/schema.sql
psql "$PG_CONN" -f stores/internal/postgres/schema.sql

echo "Loading initial data..."
psql "$PG_CONN" -f customers/internal/postgres/data.sql

echo "Database initialization complete!"
