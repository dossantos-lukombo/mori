from fastapi import FastAPI,Request, HTTPException, Depends
from pydantic import BaseModel
from llm_manager import treating_user_request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from typing import Literal
from datetime import datetime
from fastapi.responses import StreamingResponse
from datetime import datetime, timezone
import json,time,jwt,os
from dotenv import load_dotenv

app = FastAPI()

security = HTTPBearer()

load_dotenv("../.env")


"""
Data
"""
class Data(BaseModel):
    user_id: str
    conversation_id: str
    message: str
    
class SendData(BaseModel):
    status: Literal["success", "error"]
    user_id: str
    conversation_id: str
    response: str
    timestamp: str
    
global llm_response
llm_response = ""


def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security)):
    try:
        payload = jwt.decode(credentials.credentials, os.getenv("ACCESS_SECRET_KEY_LLM"), algorithms=["HS256"])
        return payload
    except jwt.ExpiredSignatureError:
        print("payload",jwt.decode(credentials.credentials, os.getenv("ACCESS_SECRET_KEY_LLM"), algorithms=["HS256"])
 )
        raise HTTPException(status_code=401, detail="Token expiré")
    except jwt.InvalidTokenError:
        print("payload",jwt.decode(credentials.credentials, os.getenv("ACCESS_SECRET_KEY_LLM"), algorithms=["HS256"])
 )

        raise HTTPException(status_code=401, detail="Token invalide")

async def generate_stream(entry_data):
    
    if entry_data["message"]!= "":
        for chunk in treating_user_request(entry_data):
            if chunk is not None:
                output = {
                    "status": "success",
                    "user_id":entry_data["user_id"],
                    "conversation_id": entry_data["conversation_id"],
                    "response": chunk,
                    "timestamp": chunk["created_at"]
                }
                yield f"data: {json.dumps(output)}\n\n"
    else:
        current_time = datetime.now(timezone.utc)
        output = {
            "status": "error",
            "user_id":entry_data["user_id"],
            "conversation_id": entry_data["conversation_id"],
            "response": "user message empty",
            "timestamp": current_time.isoformat(timespec='seconds').replace('+00:00', 'Z')
        }
        yield f"data: {json.dumps(output)}\n\n"

@app.post("/llm-protected")
async def receive_data(data: Data,credentials: HTTPAuthorizationCredentials = Depends(verify_token)):
    # Exemple de traitement des données
    print(f"Traitement des données pour Data {data.user_id}, {data.conversation_id}, {data.message}")
    entry_data = {
        "user_id": data.user_id,
        "conversation_id": data.conversation_id,
        "message": data.message
    }

    return StreamingResponse(generate_stream(entry_data),media_type="text/event-stream")

