#!/usr/bin/env sh

until psql -c "select 1" >/dev/null 2>/dev/null; do
  echo "Waiting for postgres server..."
  sleep 1
done

/app/cdbd migrate

go install github.com/cosmtrek/air@latest
air
