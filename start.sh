#!/usr/bin/env bash

docker-compose up -d database
RETRIES=6
until docker-compose exec database psql -U postgres -c "select 1" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
  sleep 10
done
docker-compose exec database psql -U postgres -c "CREATE DATABASE orderdb WITH ENCODING 'UTF8' TEMPLATE template0;"
docker-compose exec database psql -U postgres -c "CREATE USER shaumux WITH ENCRYPTED PASSWORD 'shaumux';"
docker-compose exec database psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE orderdb TO shaumux;"
docker-compose exec database psql -U postgres orderdb -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
docker-compose build goOrderAPI
docker-compose up -d goOrderAPI