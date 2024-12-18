alias dev-server := up

default:
    @just --list

ci: fmt lint test

fmt:
    docker compose exec -it server sh -c "find . -name '*.go' -exec goimports -w {} \\;"
    docker compose exec -it frontend npm run format

shell container-name:
    docker compose exec -it {{container-name}} bash

db:
    docker compose exec -it postgres psql -U postgres cdb

up:
    docker compose up -d
    docker compose logs -f server frontend cli

test:
    # This monstrosity avoids go test ./... from trying to scan all of node_modules
    # and hitting ulimits.
    docker compose exec -it server bash -c 'go test -tags testing $(find . -path "./frontend/*" -prune -o -path "./.git/*" -prune -o -name "*.go" -printf "%h\n" | sort -u)'
    # TODO: fix these tests
    # docker compose exec -it frontend npm test

lint:
    docker compose exec -it server golangci-lint run --fix
    docker compose exec -it frontend npm run lint

build:
    docker compose exec -it server go build ./cmd/cdbd

seed:
    docker compose exec -it server go run /app/scripts/seed.go

create-migration migration_name:
    docker compose exec -it server migrate create -ext sql -dir migrations -seq {{migration_name}}

migrate *FLAGS:
    docker compose exec -it server go run ./cmd/cdbd migrate {{FLAGS}}
