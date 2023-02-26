#!/usr/bin/env bash

set -Eeuo pipefail 
set -o xtrace

export $(cat .env | xargs)

docker rm postgresql || true
docker run -d  -e POSTGRES_USER=$DB_USER -e POSTGRES_PASSWORD=$DB_PASSWORD -e POSTGRES_DB=$DB_NAME \
    -p 5432:5432 \
    -v  /etc/bridge/pgdata:/var/lib/postgresql/data \
    --name postgresql postgres:latest 
