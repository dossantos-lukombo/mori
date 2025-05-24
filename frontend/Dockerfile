# ===== Étape build =====
FROM node:18.17.1-alpine3.18 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# ===== Étape prod =====
FROM node:18.17.1-alpine3.18
RUN adduser -D -u 1001 appuser
WORKDIR /app

# On évite apk + chown en copiant directement avec les bons UID/GID
COPY --from=builder --chown=appuser:appuser /app ./

USER appuser
EXPOSE 8080
CMD ["npm", "run", "serve"]