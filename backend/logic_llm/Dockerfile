FROM python:3.11-slim

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG DEBIAN_FRONTEND=noninteractive
ARG CURL_VER=7.88.1-10+deb12u12

RUN apt-get update \
 && apt-get install -y --no-install-recommends curl=${CURL_VER} \
 && curl -sSL https://install.python-poetry.org | python3 - \
 && apt-get purge -y --auto-remove curl \
 && rm -rf /var/lib/apt/lists/*

ENV POETRY_VERSION=1.7.1 \
    POETRY_HOME=/opt/poetry \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1
ENV PATH="${POETRY_HOME}/bin:${PATH}"

WORKDIR /app

COPY logic_llm/pyproject.toml logic_llm/poetry.lock ./
RUN poetry install --only main --no-root

COPY . ./logic_llm

EXPOSE 3000
CMD ["poetry", "run", "uvicorn", "logic_llm.server:app", "--host", "0.0.0.0", "--port", "3000"]