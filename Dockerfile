# syntax=docker/dockerfile:1.7-labs

# Build the server binary
FROM golang:1.22 as backend

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY --exclude=frontend,charts . .

RUN go build -o /app/cdbd ./cmd/cdbd

# Build the frontend
FROM node:lts as frontend

WORKDIR /app
COPY ./frontend/package.json .
COPY ./frontend/package-lock.json .
RUN npm install

COPY ./frontend .
RUN npm run build

# Build the final image
FROM debian:bookworm-slim as final

WORKDIR /app

RUN adduser --home /app --no-create-home --uid 600 --disabled-password cdb 

ENV ENV=deployed
COPY --chown=cdb --from=backend /app/cdbd /app/cdbd
COPY --chown=cdb --from=backend /app/docker/entrypoint.sh /app/entrypoint.sh
COPY --chown=cdb --from=frontend /app/build /app/frontend/build

ENTRYPOINT /app/entrypoint.sh
