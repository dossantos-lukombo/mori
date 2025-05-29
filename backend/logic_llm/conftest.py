# conftest.py
import os
import pytest
from dotenv import load_dotenv, find_dotenv

@pytest.fixture(scope="session", autouse=True)
def load_env():
    env_path = find_dotenv("../.env")
    if not load_dotenv(env_path):
        pytest.exit("❌ Impossible de charger le .env", returncode=1)
