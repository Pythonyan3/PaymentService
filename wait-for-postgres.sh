#!/bin/sh

cmd="$1"

until PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -c '\q'; do
  echo "Waiting for postgres..."
  sleep 1
done

echo "Running migrations..."

migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -path migrations up

exec $cmd