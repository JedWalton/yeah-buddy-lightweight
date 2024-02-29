#!/bin/bash

# Load environment variables from .env file in project root
set -a
source ../.env
set +a

BACKUP_DIR="../db.bak"  # Adjust this directory if needed

# Find the most recent backup file
LATEST_BACKUP=$(ls -t $BACKUP_DIR/db-backup-*.sql | head -1)

if [ -z "$LATEST_BACKUP" ]; then
    echo "No backup found in $BACKUP_DIR"
    exit 1
fi

# Restore from the latest backup
PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -p 5432 -U $POSTGRES_USER -d $POSTGRES_DB -f "$LATEST_BACKUP"

echo "Database restored from $LATEST_BACKUP"

