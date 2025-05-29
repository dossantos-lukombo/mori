# backend/Dockerfile
# Étape 1 : builder sur une image Golang
FROM golang:1.24 AS builder
WORKDIR /app

# 1. Copier les fichiers de modules pour profiter du cache
COPY go.mod go.sum ./
RUN go mod download

# 2. Copier le code et compiler
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mori_backend .

# Étape 2 : image de production légère
FROM debian:bullseye-slim

ENV DEBIAN_FRONTEND=noninteractive \
    TZ=Europe/Paris

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates \
      tzdata \
      curl &&\
    rm -rf /var/lib/apt/lists/*

# Utiliser Bash pour une meilleure robustesse des scripts
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Création d'un utilisateur non-root pour plus de sécurité
RUN useradd -u 1001 -r -s /bin/false appuser

# Configuration du fuseau horaire
ENV TZ=Europe/Paris

# Variables d'environnement pour la connexion Supabase
# À surcharger au runtime avec -e ou via Docker Compose / Secrets
ENV DATABASE_URL="" \
    SUPABASE_SERVICE_ROLE_KEY=""

# Copie des fichiers de migration SQL
COPY --from=builder /app/pkg/db/migration/PostgreSql /app/pkg/db/migration/PostgreSql

# Copie de l'exécutable compilé
COPY --from=builder /app/mori_backend /app/

# Définition du répertoire de travail
WORKDIR /app

# Création du dossier pour les uploads de fichiers avec les bonnes permissions
RUN mkdir -p /app/fileUploads \
    && chown appuser:appuser /app/fileUploads

# Passage en utilisateur non-root
USER appuser

# Exposition du port
EXPOSE 8081

# Commande de démarrage
CMD ["/app/mori_backend"]
