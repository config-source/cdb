FROM golang:1.22

WORKDIR /app

RUN apt update && apt install -y postgresql-client

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /app/cdbd ./cmd/cdbd

ENTRYPOINT /app/docker/entrypoint.sh
