#!/usr/bin/env bash

set -Eeuo pipefail 
set -o xtrace

export $(cat .env | xargs)

db_url=postgres://$DB_USER:$DB_PASSWORD@$DB_HOST/$DB_NAME?sslmode=disable

export PGPASSWORD=$DB_PASSWORD
psql -h localhost -U $DB_USER -d $DB_NAME -a -f  db/migrations/create_tables.sql