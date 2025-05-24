# ===== Étape build =====
# FROM node:18.17.1-alpine3.18 AS builder
# WORKDIR /app
# COPY package*.json ./
# RUN npm ci
# COPY . .
# RUN npm run build

# # ===== Étape prod =====
# FROM node:18.17.1-alpine3.18
# RUN adduser -D -u 1001 appuser
# WORKDIR /app

# # On évite apk + chown en copiant directement avec les bons UID/GID
# COPY --from=builder --chown=appuser:appuser /app ./

# USER appuser
# EXPOSE 8080
# CMD ["npm", "run", "serve"]

# frontend/Dockerfile
FROM node:18.17.1-alpine3.18

WORKDIR /app

# Utilisateur non-root minimal pour Alpine
RUN adduser -D -u 1001 appuser \
 && chown -R appuser:appuser /app

USER appuser

CMD ["node", "-e", "console.log('Frontend image OK')"]