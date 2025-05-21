# backend/logic_llm/Dockerfile (FastAPI avec Poetry)
FROM python:3.11-slim

# Poetry installation
ENV POETRY_VERSION=1.5.1 \
    POETRY_HOME="/opt/poetry" \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1

RUN apt-get update && apt-get install -y --no-install-recommends curl \
    && curl -sSL https://install.python-poetry.org | python3 - \
    && apt-get purge -y curl \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="$POETRY_HOME/bin:$PATH"
WORKDIR /app

# Copier la config Poetry depuis la racine
COPY ../../pyproject.toml ../../poetry.lock* ./
RUN poetry install --only main

# Copier le code
COPY . .

EXPOSE 3000
CMD ["poetry", "run", "uvicorn", "server:app", "--host", "0.0.0.0", "--port", "3000"]