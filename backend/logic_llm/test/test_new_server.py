import pytest
import os
import json
from fastapi.testclient import TestClient
from mori.server import app
from mori.llm_manager import treating_user_request
from datetime import datetime, timedelta, timezone
from jose import jwt
from dotenv import load_dotenv

# Stub pour éviter l'appel réel à Ollama en CI\ n@pytest.fixture(autouse=True)
def stub_chat(monkeypatch):
    def fake_treating_user_request(entry_data):
        # Renvoie un chunk factice
        yield {
            "message": {"content": "Dummy response"},
            "created_at": datetime.now(timezone.utc).isoformat()
        }
    monkeypatch.setattr("mori.llm_manager.treating_user_request", fake_treating_user_request)

# Charge le .env pour les secrets
if not load_dotenv():
    pytest.skip("Could not load .env file", allow_module_level=True)

client = TestClient(app, base_url="http://test")

# Fonction utilitaire pour créer un JWT valide

def create_jwt_token():
    payload = {
        "sub": "test_user",
        "exp": datetime.now(tz=timezone.utc) + timedelta(hours=1),
        "iat": datetime.now(tz=timezone.utc),
        "scope": "user"
    }
    secret_key = os.getenv("ACCESS_SECRET_KEY_LLM")
    return jwt.encode(payload, secret_key, algorithm="HS256")


def test_health_check():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_receive_data():
    token = create_jwt_token()
    data = {"user_id": "test_user", "conversation_id": "test_convo", "message": "Hello, how are you?"}
    headers = {"Authorization": f"Bearer {token}"}
    response = client.post("/llm-protected", json=data, headers=headers)
    assert response.status_code == 200
    # Vérifie que le SSE contient notre chunk factice
    assert "data:" in response.text
    # Parse le JSON du premier chunk
    chunk = response.text.split("data: ")[1].strip()
    payload = json.loads(chunk)
    assert payload["status"] == "success"
    assert payload["response"] == "Dummy response"


def test_receive_data_empty_message():
    token = create_jwt_token()
    data = {"user_id": "test_user", "conversation_id": "test_convo", "message": ""}
    headers = {"Authorization": f"Bearer {token}"}
    response = client.post("/llm-protected", json=data, headers=headers)
    assert response.status_code == 200
    assert "user message empty" in response.text
