import pytest
from fastapi.testclient import TestClient
from logic_llm.server import app
import json

client = TestClient(app)

@pytest.mark.asyncio
async def test_receive_data_success():
    data = {
        "user_id": "test_user",
        "conversation_id": "test_conversation",
        "message": "Hello, this is a test message."
    }
    
    response = client.post("/api/process-message", json=data)
    
    assert response.status_code == 200
    # assert response != ""
    for line in response.iter_lines():
        print("line: ",line)
        if line != "":
            
            line_data = line.replace("data: ", "")
            response_json = json.loads(line_data)
            assert "status" in response_json
            assert response_json["status"] == "success" or response_json["status"] == "error"

@pytest.mark.asyncio
async def test_receive_data_error():
    # Simulez un cas où `treating_user_request` retourne une réponse vide
    data = {
        "user_id": "test_user",
        "conversation_id": "test_conversation",
        "message": ""
    }
    
    # Envoyer une requête POST à l'endpoint
    response = client.post("/api/process-message", json=data)
    
    # Assurez-vous que la réponse est correcte
    assert response.status_code == 200
    
    # Lisez le contenu de la réponse en tant que flux
    for line in response.iter_lines():
        print("line: ", line)
        if line != "":
            
            line_data = line.replace("data: ", "")
            response_json = json.loads(line_data)
            assert "status" in response_json
            assert response_json["status"] == "error"

