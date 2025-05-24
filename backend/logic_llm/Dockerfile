# backend/logic_llm/Dockerfile
FROM python:3.11-slim

# 1) Toutes les instructions RUN seront évaluées par bash avec pipefail
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# 2) Pour les installs non interactives
ARG DEBIAN_FRONTEND=noninteractive

# 3) Version de curl (vérifie avec: apt-cache policy curl | head -n 2)
ARG CURL_VER=7.88.1-10+deb12u12   # ← adapte si besoin

# 4) Installation de Poetry + ménage
RUN apt-get update \
 && apt-get install -y --no-install-recommends curl=${CURL_VER} \
 && curl -sSL https://install.python-poetry.org | python3 - \
 && apt-get purge -y --auto-remove curl \
 && rm -rf /var/lib/apt/lists/*

# Variables d’environnement Poetry
ENV POETRY_VERSION=1.7.1 \
    POETRY_HOME=/opt/poetry \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1 \
    PATH="${POETRY_HOME}/bin:${PATH}"

WORKDIR /app

# Copier les fichiers de dépendances
COPY logic_llm/pyproject.toml logic_llm/poetry.lock ./

# Installer seulement les dépendances “main”
RUN poetry install --only main --no-root

# Copier le reste du code
COPY . ./logic_llm

EXPOSE 3000
CMD ["poetry", "run", "uvicorn", "logic_llm.server:app", "--host", "0.0.0.0", "--port", "3000"]