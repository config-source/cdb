source /etc/bash_completion
source <"$(go run ./cmd/cdbd completion bash)"

function test() {
    go test "$@" -tags testing $(find . -path "./frontend/*" -prune -o -path "./.git/*" -prune -o -name "*.go" -printf "%h\n" | sort -u)
}

function cdbd() {
    go run ./cmd/cdbd "$@"
}
