#!/bin/bash

# Load environment variables from .env file in project root
set -a
source ./.env.production.local
set +a

# Extract the components from POSTGRESQL_URL
PROTO="$(echo $POSTGRESQL_URL | grep :// | sed -e's,^\(.*://\).*,\1,g')"
URL="$(echo ${POSTGRESQL_URL/$PROTO/})"
USER="$(echo $URL | grep @ | cut -d: -f1)"
PASS="$(echo $URL | grep @ | cut -d@ -f1 | cut -d: -f2)"
HOST="$(echo ${URL/$USER:$PASS@/} | cut -d: -f1)"
PORT="$(echo ${URL/$USER:$PASS@/} | cut -d: -f2 | cut -d/ -f1)"
DB="$(echo $URL | grep / | cut -d/ -f2-)"

# Prompt for password to avoid showing it in the command
export PGPASSWORD=$PASS

# Set the date format and backup directory
BACKUP_DATE=$(date +%Y%m%d%H%M%S)
BACKUP_DIR="./database_backups"  # Ensure this directory exists

# Execute pg_dump
pg_dump -U "$USER" -h "$HOST" -p "$PORT" -d "$DB" > "$BACKUP_DIR/db-backup-$BACKUP_DATE.sql"

# Clear the password variable
unset PGPASSWORD
