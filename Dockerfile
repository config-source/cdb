# syntax=docker/dockerfile:1.7-labs

# Build the server binary
FROM golang:1.22 as backend

WORKDIR /app
RUN adduser --home /app --no-create-home --uid 600 --disabled-password cdb 
ENV ENV=deployed

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /app/cdbd ./cmd/cdbd
ENTRYPOINT /app/cdbd
