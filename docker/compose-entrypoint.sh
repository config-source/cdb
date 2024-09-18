#!/usr/bin/env bash

function backend() {
  if [[ ! -f /usr/bin/psql ]]; then
    # Install the postgres client
    apt update && apt install -y postgresql-client
  fi

  if [[ ! -x $(which migrate) ]]; then
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  fi

  if [[ ! -x $(which golangci-migrate) ]]; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  fi

  if [[ ! -x $(which goimports) ]]; then
    go install golang.org/x/tools/cmd/goimports@latest
  fi

  until psql -c "select 1" >/dev/null 2>/dev/null; do
    echo "Waiting for postgres server..."
    sleep 1
  done

  go run /app/cmd/cdbd migrate

  go install github.com/air-verse/air@latest
  air
}

function live_build_cli() {
  go install github.com/air-verse/air@latest
  air \
    --build.cmd "go build -buildvcs=false -o ./tmp/cdb ./cmd/cdb" \
    --build.bin "./tmp/cdb" \
    --build.exclude_dir "assets,tmp,vendor,testdata,scripts,migrations,frontend" \
    --build.args_bin "\-\-help"
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

cli)
  live_build_cli
  ;;

*)
  backend
  ;;
esac
