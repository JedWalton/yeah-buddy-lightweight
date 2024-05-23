#!/bin/sh

# Install migrate with these commands.
# sudo apt-get install curl -y
# curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
# sudo mv migrate /usr/local/bin/migrate

# Load environment variables from .env file
. ../.env

# Function to clear all tables in the database
clear_tables() {
  # Connect to the PostgreSQL database using psql and execute the SQL command to drop all tables
  psql "${POSTGRESQL_URL_DEV}" <<-EOSQL
    DO
    \$\$
    DECLARE
        _tbl text;
    BEGIN
        -- Loop through all tables in the 'public' schema
        FOR _tbl IN
            SELECT tablename
            FROM pg_tables
            WHERE schemaname = 'public'
        LOOP
            -- Drop each table with cascade to automatically remove dependent entities
            EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(_tbl) || ' CASCADE';
        END LOOP;
    END
    \$\$
    ;
EOSQL
}

# Output the date and start message, then run the function to clear tables
echo "[`date`] Clearing all tables in the database..."
clear_tables
echo "[`date`] All tables cleared."

# Run DB migrations
echo "[`date`] Running DB migrations..."
migrate -database "${POSTGRESQL_URL_DEV}" -path ../migrations up

# Output success message
echo "[`date`] Database setup complete."

