# Étape 1 : builder
FROM golang:1.20-alpine AS builder

# Installer git et les certificats pour que go mod download puisse fonctionner
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 1. Copier uniquement go.mod/go.sum pour tirer parti du cache Docker
COPY go.mod go.sum ./

# 2. Télécharger les dépendances Go
RUN go mod download

# 3. Copier le reste et compiler
COPY . .
RUN go build -o social-network .

# Étape 2 : image finale minimale
FROM alpine:latest
# Installer seulement ca-certificates si vous communiquez en HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /app

# 4. Copier le binaire compilé
COPY --from=builder /app/social-network ./

EXPOSE 8081

# 5. Lancer le binaire
CMD ["./social-network"]
