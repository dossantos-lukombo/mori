# Étape 1 : build
FROM node:18.17.1-alpine3.18 AS builder
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Étape 2 : image de production
FROM node:18.17.1-alpine3.18
WORKDIR /app

#  Une seule couche RUN pour tout faire :
RUN adduser -D -u 1001 appuser \
 && apk add --no-cache --virtual .tmp-chown shadow \
 && chown -R appuser:appuser /app \
 && apk del .tmp-chown

COPY package*.json ./
USER root            # on reste root le temps d’installer
RUN npm install --omit=dev \
 && chown -R appuser:appuser /app

USER appuser
EXPOSE 8080
CMD ["npm", "run", "serve"]