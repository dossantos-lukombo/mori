import pytest
import os
import json
from fastapi.testclient import TestClient
import mori.server as server_mod
from datetime import datetime, timedelta, timezone
from jose import jwt
from dotenv import load_dotenv

# Stub Ollama via patch dans server_mod
@pytest.fixture(autouse=True)
def stub_chat(monkeypatch):
    def fake_treating_user_request(entry_data):
        # Renvoie un chunk factice pour tests
        yield {
            "message": {"content": "Dummy response"},
            "created_at": datetime.now(timezone.utc).isoformat()
        }
    # Patch dans le module server pour que generate_stream utilise fake
    monkeypatch.setattr(server_mod, 'treating_user_request', fake_treating_user_request)

# Charge les variables d'env depuis .env
if not load_dotenv():
    pytest.skip("Could not load .env file", allow_module_level=True)

# Client FastAPI
client = TestClient(server_mod.app, base_url="http://test")

# Utilitaire pour créer un JWT valide

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
    data = {
        "user_id": "test_user",
        "conversation_id": "test_convo",
        "message": "Hello, how are you?"
    }
    headers = {"Authorization": f"Bearer {token}"}
    response = client.post("/llm-protected", json=data, headers=headers)
    assert response.status_code == 200
    # Vérifie que la réponse SSE contient notre chunk factice
    assert "data:" in response.text
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
