FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o /app/cdbd ./cmd/cdbd

ENTRYPOINT /app/docker/entrypoint.sh
