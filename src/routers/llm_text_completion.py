
import os
import json
import requests
from fastapi import APIRouter, Response
from src.logger import Logger
from src.routers.models.llm_text_completion import TextCompletionRequest, TextCompletionResponse
from src.routers.models.user import AuthUser

OLLAMA_HOST = os.getenv("OLLAMA_HOST","http://localhost:9999")
OLLAMA_MODEL = os.getenv("OLLAMA_MODEL", "gemma3:1b")
DOKOOLA_BACKEND_API = os.getenv("DOKOOLA_BACKEND_API", "http://localhost:8000/api")

logger = Logger(__name__)
router = APIRouter(
    tags=["LLM Text Completion"]
)

@router.post(
    "/text-completion/{userPublicId}",
    response_model=TextCompletionResponse,
    description="Completes a given text using a language model.",
    responses={
        200: {"detail": "Text completion finished successfully"},
        400: {"detail": "Invalid request body"},
        404: {"detail": "User not found"},
        500: {"detail": "Internal server error"},
    },
)

def llm_text_completion(request:TextCompletionRequest, userPublicId: str, respone: Response):
    try:
        logger.info(f"llm_text_completion request received for user {userPublicId}")

        user = get_user_by_public_id(userPublicId)
        if not user:
            respone.status_code = 404
            return TextCompletionResponse(success=False, error_message="User not found")
        
        completion = engage_ai_for_text_completion(request.text, user)
        return TextCompletionResponse(completion=completion, success=True) 
    except Exception as e:
        logger.error(f"llm_text_completion request failed: {e}")
        respone.status_code = 500
        return TextCompletionResponse(success=False, error_message="Internal server error")
    finally:
        logger.info(f"llm_text_completion request completed for user {userPublicId} with status code {respone.status_code}")
    

def get_user_by_public_id(public_id: str):
    try:
        response = requests.get(f"{DOKOOLA_BACKEND_API}/users/auth/{public_id}/")
        response.raise_for_status()
        user = response.json()
        return AuthUser(**user)
    except Exception as e:
        logger.error(f"Failed to get user by public id {public_id} [Error:] {e}")

def engage_ai_for_text_completion(prompt: str, user: AuthUser) -> str:

    payload = {
        "model": OLLAMA_MODEL,
        "messages": [
            {
                "role": "system",
                "content": "Welcome to Dokoola Platform! I am your dedicated AI assistant, ready to help you with any questions or tasks you have."
            },
            {
                "role": "user",
                "content": prompt
            }
        ]
    }

    try:
        generated_text = ""
        response = requests.post(f'{OLLAMA_HOST}/api/chat', json=payload, stream=True)
        if response.status_code == 200:
            logger.info("Streaming response from Ollama:")

            for line in response.iter_lines(decode_unicode=True):
                if line:
                    try:
                        json_line = json.loads(line)
                        if "message" in json_line and "content" in json_line["message"]:
                            logger.info(json_line["message"]["content"])
                            generated_text += json_line["message"]["content"]
                    except json.JSONDecodeError:
                        logger.error(f"Failed to decode JSON from line: {line}")
            return generated_text
        else:
            logger.error(f"Failed to get a valid response from Ollama. Status code: {response.status_code}")
            return f"Completed: {prompt}"

    except requests.exceptions.RequestException as e:
        logger.error(f"Error calling AI model: {e}")
        return f"Completed: {prompt}"
