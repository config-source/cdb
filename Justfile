alias dev-server := up

ci: lint test

up:
    docker compose up -d
    docker compose logs -f server frontend

test:
    # This monstrosity avoids go test ./... from trying to scan all of node_modules
    # and hitting ulimits.
    docker compose exec -it server bash -c 'go test -tags testing $(find . -path "./frontend/*" -prune -o -path "./.git/*" -prune -o -name "*.go" -printf "%h\n" | sort -u)'

lint:
    #!/usr/bin/env bash
    if [[ ! -x $(which golangci-lint 2>/dev/null) ]]; then
        echo "Linter not installed, installing..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v1.57.2
    fi

    docker compose exec -it server golangci-lint run --fix

build:
    docker compose exec -it server go build ./cmd/cdbd

seed:
    docker compose exec -it server go run /app/scripts/seed.go

create-migration migration_name:
    #!/usr/bin/env bash
    if [[ ! -x $(which migrate) ]]; then
        echo "Install the golang-migrate CLI."
        printf "\tInstructions available here: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate\n"
        exit 1
    fi

    migrate create -ext sql -dir migrations -seq {{migration_name}}
