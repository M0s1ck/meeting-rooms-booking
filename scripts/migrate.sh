#!/bin/sh
# One-off admin process: applies DB migrations using the same env config as the app.
# Usage: migrate.sh [migrate args...]   (default: up)
set -eu

: "${POSTGRES_HOST:?POSTGRES_HOST is required}"
: "${POSTGRES_DB:?POSTGRES_DB is required}"
: "${POSTGRES_USER:?POSTGRES_USER is required}"
: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"

DB_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT:-5432}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE:-disable}"

if [ "$#" -eq 0 ]; then
  set -- up
fi

exec migrate -path "${MIGRATIONS_PATH:-/app/migrations}" -database "$DB_URL" "$@"
