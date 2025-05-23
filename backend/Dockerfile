# backend/Dockerfile
# Étape 1 : builder sur une image Golang
FROM golang:1.23 AS builder
WORKDIR /app
# 1. Copier les fichiers de modules pour profiter du cache
COPY go.mod go.sum ./
RUN go mod download
# 2. Copier le code et compiler
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mori_backend .

# Étape 2 : image de production légère
FROM debian:bullseye-slim
# Installation des certificats et dépendances nécessaires
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates \
      tzdata \
      && rm -rf /var/lib/apt/lists/*
# Création d'un utilisateur non-root pour plus de sécurité
RUN useradd -u 1001 -r -s /bin/false appuser
# Configuration du fuseau horaire
ENV TZ=Europe/Paris
# Copie des fichiers de migration SQL si nécessaire
COPY --from=builder /app/pkg/db/migration/PostgreSql /app/pkg/db/migration/PostgreSql
# Copie de l'exécutable compilé
COPY --from=builder /app/mori_backend /app/
# Définition du répertoire de travail
WORKDIR /app
# Création du dossier pour les uploads de fichiers avec les bonnes permissions
RUN mkdir -p /app/fileUploads && chown appuser:appuser /app/fileUploads
# Exposition du port
EXPOSE 8081
# Point d'entrée
ENTRYPOINT ["/app/mori_backend"]

