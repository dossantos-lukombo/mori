import pytest
from fastapi.testclient import TestClient
from logic_llm.server import app
import json
import jwt
from datetime import datetime, timedelta, timezone
import os

client = TestClient(app)

# Clé secrète utilisée pour signer le token JWT
SECRET_KEY = os.getenv("ACCESS_SECRET_KEY_LLM")

def create_jwt_token():
    # Définir les revendications du token
    payload = {
        "sub": "test_user",
        "exp": datetime.now(tz=timezone.utc) + timedelta(hours=1),  # Expiration dans 1 heure
        "iat": datetime.now(tz=timezone.utc),  # Heure d'émission
        "scope": "user"
    }
    # Générer le token JWT
    token = jwt.encode(payload, SECRET_KEY, algorithm="HS256")
    return token

@pytest.mark.asyncio
async def test_receive_data_success():
    # Générer un token JWT valide
    token = create_jwt_token()
    
    data = {
        "user_id": "test_user",
        "conversation_id": "test_conversation",
        "message": "Hello, this is a test message."
    }
    
    # Ajouter le token dans les en-têtes de la requête
    headers = {
        "Authorization": f"Bearer {token}"
    }
    
    response = client.post("/api/process-message", json=data, headers=headers)
    
    assert response.status_code == 200
    for line in response.iter_lines():
        print("line: ", line)
        if line:
            line_data = line.replace("data: ", "")
            response_json = json.loads(line_data)
            assert "status" in response_json
            assert response_json["status"] in ["success", "error"]

@pytest.mark.asyncio
async def test_receive_data_error():
    # Générer un token JWT valide
    token = create_jwt_token()
    
    # Cas où le message est vide
    data = {
        "user_id": "test_user",
        "conversation_id": "test_conversation",
        "message": ""
    }
    
    # Ajouter le token dans les en-têtes de la requête
    headers = {
        "Authorization": f"Bearer {token}"
    }
    
    response = client.post("/api/process-message", json=data, headers=headers)
    
    assert response.status_code == 200
    for line in response.iter_lines():
        print("line: ", line)
        if line:
            line_data = line.replace("data: ", "")
            response_json = json.loads(line_data)
            assert "status" in response_json
            assert response_json["status"] == "error"
