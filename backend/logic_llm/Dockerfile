FROM python:3.11-slim

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG DEBIAN_FRONTEND=noninteractive
ARG CURL_VER=7.88.1-10+deb12u12

ENV POETRY_VERSION=1.7.1 \
    POETRY_HOME=/opt/poetry \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1
ENV PATH="${POETRY_HOME}/bin:${PATH}"

RUN apt-get update \
 && apt-get install -y --no-install-recommends curl=${CURL_VER} \
 && POETRY_HOME=$POETRY_HOME curl -sSL https://install.python-poetry.org | python3 - \
 && apt-get purge -y --auto-remove curl \
 && rm -rf /var/lib/apt/lists/*

# backend/logic_llm/Dockerfile
WORKDIR /app

# On est déjà DANS backend/logic_llm : plus besoin du sous-dossier
COPY pyproject.toml poetry.lock ./
RUN poetry config virtualenvs.create false
RUN poetry install --no-root 

# On copie tout le code du dossier courant
COPY . .

# Le reste inchangé…
CMD ["uvicorn", "mori.server:app", "--host", "0.0.0.0", "--port", "3000"]