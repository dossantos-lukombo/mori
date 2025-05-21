# backend/Dockerfile (Go Backend)
FROM golang:1.20-bullseye AS builder

WORKDIR /app

# 1. Copier uniquement go.mod/go.sum pour profiter du cache Docker
COPY go.mod go.sum ./

# 2. Télécharger les dépendances (git + SSL sont déjà présents)
RUN go mod download

# 3. Copier le code et compiler
COPY . .
RUN go build -o social-network .

# Étape finale : image minimale
FROM debian:bullseye-slim
# ca-certificates est déjà là, mais on le reconfirme
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/social-network ./

EXPOSE 8081
CMD ["./social-network"]
