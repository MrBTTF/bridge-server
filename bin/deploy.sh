#!/usr/bin/env bash

set -Eeuo pipefail

export $(cat prod/.env | xargs)

echo Building
./bin/build.sh
echo Building Docker image
docker build --tag fuji:5000/bridge-server:latest .
echo Pushing Docker image
version=$(IFS=. read -r a b c<<<"$(cat version.txt)";echo "$a.$b.$((c+1))")
echo $version > version.txt
docker tag fuji:5000/bridge-server:latest fuji:5000/bridge-server:$version
docker push fuji:5000/bridge-server:latest
docker push fuji:5000/bridge-server:$version
echo Deploying Helm chart
helm upgrade --install bridge-server \
    --set image.tag=$version \
    --set postgresql.auth.postgresPassword=$DB_ROOT_PASSWORD \
    --set postgresql.auth.username=$DB_USER \
    --set postgresql.auth.password=$DB_PASSWORD \
    --set-file db_init=db/migrations/create_tables.sql \
    ./chart/