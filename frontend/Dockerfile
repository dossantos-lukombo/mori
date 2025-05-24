# Étape 1: Build de l'application
FROM node:18.17.1-alpine3.18 AS builder
# Définition du répertoire de travail
WORKDIR /app
# Copie des fichiers de dépendances
COPY package*.json ./
# Installation des dépendances
RUN npm install
# Copie du reste du code source
COPY . .
# Build de l'application
RUN npm run build
# Étape 2: Image de production
# Image de développement Vue.js
FROM node:18.17.1-alpine3.18

# Création d'un utilisateur non-root
RUN adduser -D -u 1001 appuser

# Définition du répertoire de travail
WORKDIR /app

# Copie des fichiers de dépendances
COPY package*.json ./

# Installation des dépendances
RUN npm install

# Attribution des permissions
RUN chown -R appuser:appuser /app

# Changement vers l'utilisateur non-root
USER appuser

# Exposition du port
EXPOSE 8080

# Commande de démarrage
CMD ["npm", "run", "serve"]