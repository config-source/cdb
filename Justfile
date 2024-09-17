alias dev-server := up

ci: lint test

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

lint:
    docker compose exec -it server golangci-lint run --fix

build:
    docker compose exec -it server go build ./cmd/cdbd

seed:
    docker compose exec -it server go run /app/scripts/seed.go

create-migration migration_name:
    docker compose exec -it server migrate create -ext sql -dir migrations -seq {{migration_name}}

migrate *FLAGS:
    docker compose exec -it server go run ./cmd/cdbd migrate {{FLAGS}}
