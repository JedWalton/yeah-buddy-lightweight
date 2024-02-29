#!/bin/bash

# Load environment variables from .env file in project root
set -a
source ../.env
set +a

# Set the date format and backup directory
BACKUP_DATE=$(date +%Y%m%d%H%M%S)
BACKUP_DIR="../db.bak"  # Ensure this directory exists

# Create a backup
pg_dump -h localhost -p 5432 -U $POSTGRES_USER -d $POSTGRES_DB > "$BACKUP_DIR/db-backup-$BACKUP_DATE.sql"


echo "Backup created at $BACKUP_DIR/db-backup-$BACKUP_DATE.sql"

