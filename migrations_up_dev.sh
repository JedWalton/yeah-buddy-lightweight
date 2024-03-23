#!/bin/sh

## Install migrate with these commands.
# sudo apt-get install curl -y && \
# curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
# sudo mv migrate /usr/local/bin/migrate

# Load environment variables from .env file
. ./.env

# Run DB migrations
echo "[`date`] Running DB migrations..." && \
    migrate -database "${POSTGRESQL_URL_DEV}" -path ./app/migrations up
