# syntax=docker/dockerfile:1.7-labs

# Build the server binary
FROM golang:1.22 as backend

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY --exclude=frontend --exclude=charts . .

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
COPY --chown=cdb --from=backend /app/migrations /app/migrations
COPY --chown=cdb --from=backend /app/docker/entrypoint.sh /app/entrypoint.sh
COPY --chown=cdb --from=frontend /app/build /app/frontend/build

RUN chown -R cdb:cdb /app
RUN apt update && apt install -y bash-completion

USER cdb
RUN echo "source /etc/bash_completion" > /app/.bashrc
RUN /app/cdbd completion bash >> /app/.bashrc

ENTRYPOINT /app/entrypoint.sh
