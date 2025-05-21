# backend/Dockerfile

# Étape 1 : builder sur une image Debian (inclut git + ca-certificates)
FROM golang:1.20 AS builder
WORKDIR /app

# 1. Copier les fichiers de modules pour profiter du cache
COPY go.mod go.sum ./
RUN go mod download

# 2. Copier le code et compiler
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o social-network .

# Étape 2 : image de production légère, on réutilise Debian slim
FROM debian:bullseye-slim

# → Pin dynamique de ca-certificates pour DL3008
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates=$(apt-cache policy ca-certificates | awk '/Candidate/ {print $2}') && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/social-network ./
EXPOSE 8081
ENTRYPOINT ["./social-network"]

