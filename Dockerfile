# syntax=docker/dockerfile:1.7-labs

# Build the server binary
FROM golang:1.22 as backend

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY --exclude=frontend . .

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
FROM fedora:latest as final

WORKDIR /app

COPY --from=backend /app/cdbd /app/cdbd
COPY --from=backend /app/docker/entrypoint.sh /app/entrypoint.sh
COPY --from=frontend /app/build /app/frontend/build

ENTRYPOINT /app/entrypoint.sh
