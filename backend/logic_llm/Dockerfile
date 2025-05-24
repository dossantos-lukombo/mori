# backend/logic_llm/Dockerfile
# Image de base
FROM python:3.11-slim
# Installation de Poetry
ENV POETRY_VERSION=1.7.1 \
    POETRY_HOME=/opt/poetry \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1

RUN apt-get update \
 && apt-get install -y --no-install-recommends curl \
 && curl -sSL https://install.python-poetry.org | python3 - \
 && apt-get purge -y curl \
 && rm -rf /var/lib/apt/lists/*

ENV PATH="${POETRY_HOME}/bin:${PATH}"

# Répertoire de travail
WORKDIR /app

COPY logic_llm/pyproject.toml logic_llm/poetry.lock ./

# Installer uniquement les dépendances “main”, sans installer le projet
RUN poetry install --only main --no-root

COPY . ./logic_llm

# Exposition & lancement
EXPOSE 3000
CMD ["poetry", "run", "uvicorn", "logic_llm.server:app", "--host", "0.0.0.0", "--port", "3000"]