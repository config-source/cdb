#!/usr/bin/env bash

function backend() {
  if [[ ! -f /usr/bin/psql ]]; then
    # Install the postgres client
    apt update && apt install -y postgresql-client
  fi

  until psql -c "select 1" >/dev/null 2>/dev/null; do
    echo "Waiting for postgres server..."
    sleep 1
  done

  go run /app/cmd/cdbd migrate

  go install github.com/air-verse/air@latest
  air
}

function frontend() {
  echo "Checking for node_modules"
  if [[ ! -f node_modules/.package-lock.json ]]; then
    echo "node_modules not found, installing deps."
    npm ci
  fi

  npm run dev -- --host
}

case $1 in
frontend)
  frontend
  ;;

backend)
  backend
  ;;

*)
  backend
  ;;
esac
