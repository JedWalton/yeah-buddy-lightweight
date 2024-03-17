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

# Directory where backups are stored
BACKUP_DIR="./database_backups"

# Find the most recent backup file
LATEST_BACKUP=$(ls -t $BACKUP_DIR/db-backup-*.sql | head -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    echo "No backup file found in $BACKUP_DIR"
    exit 1
fi

echo "Restoring database from $LATEST_BACKUP"

# Prompt for password to avoid showing it in the command
export PGPASSWORD=$PASS

# Attempt to drop all tables (and other objects) in the database
echo "Attempting to clean the database: $DB"
psql -U "$USER" -h "$HOST" -p "$PORT" "$DB" <<EOSQL
BEGIN;
DO \$\$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END
\$\$;
COMMIT;
EOSQL

# Restore the database from the latest backup
psql -U "$USER" -h "$HOST" -p "$PORT" "$DB" < "$LATEST_BACKUP"

# Clear the password variable
unset PGPASSWORD

echo "Database restoration completed successfully."
