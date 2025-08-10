import json
import os
import random

import requests
from fastapi import APIRouter, Response

from src.constant import MODELS
from src.logger import Logger
from src.routers.models.llm_text_completion import (
    TextCompletionRequest,
    TextCompletionResponse,
)
from src.routers.models.user import AuthUser

LLM_URL = os.getenv("LLM_URL", None)
LLM_API_KEY = os.getenv("LLM_API_KEY", None)
AUTH_SERVER_API = os.getenv("AUTH_SERVER_API")

# Define system messages that set the context and role for the AI assistant
# These messages help establish the AI's identity as a career assistant on the Dokoola platform
MESSAGES = [
    {"role": "system", "content": "Welcome to Dokoola Platform!"},
    {
        "role": "system",
        "content": "A digital marketplace connecting talent to opportunities across The Gambia and soon Africa",
    },
    {"role": "system", "content": "You are an expert career assistant."},
]

current_model_index = random.randrange(0, len(MODELS) - 1)

logger = Logger(__name__)
router = APIRouter(tags=["LLM Text Completion"])


@router.post(
    "/chat/completion/{userPublicId}",
    response_model=TextCompletionResponse,
    description="Completes a given text using a language model.",
    responses={
        200: {"detail": "Text completion finished successfully"},
        400: {"detail": "Invalid request body"},
        404: {"detail": "User not found"},
        500: {"detail": "Internal server error"},
    },
)
def llm_text_completion(
    request: TextCompletionRequest, userPublicId: str, respone: Response
):
    try:
        logger.info(f"llm_text_completion request received for user {userPublicId}")

        user = get_user_by_public_id(userPublicId)
        if not user:
            respone.status_code = 404
            return TextCompletionResponse(success=False, error_message="User not found")

        [completion, status_code] = get_ai_text_completion(request.text, user)
        if not completion:
            respone.status_code = status_code
            return TextCompletionResponse(
                success=False, error_message="Internal server error"
            )

        logger.info(
            f"llm_text_completion request completed for user {userPublicId} with status code {respone.status_code}"
        )
        return TextCompletionResponse(completion=completion, success=True)
    except Exception as e:
        logger.error(f"llm_text_completion request failed: {e}")
        respone.status_code = 500
        return TextCompletionResponse(
            success=False, error_message="Internal server error"
        )
    finally:
        logger.info(
            f"llm_text_completion request completed for user {userPublicId} with status code {respone.status_code}"
        )


def get_user_by_public_id(public_id: str):
    try:
        response = requests.get(f"{AUTH_SERVER_API}/users/auth/{public_id}/")
        response.raise_for_status()
        user = response.json()
        return AuthUser(**user)
    except Exception as e:
        logger.error(f"Failed to get user by public id {public_id} [Error:] {e}")


def get_ai_text_completion(prompt: str, user: AuthUser, max_retries=3) -> str | None:

    # Manage model rotation using global index
    global current_model_index
    _model = MODELS[current_model_index]
    current_model_index = (current_model_index + 1) % len(MODELS)

    # Reset index if it exceeds models length
    if current_model_index > len(MODELS):
        current_model_index = 0

    # Set up request headers with authentication
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {LLM_API_KEY}",
    }

    # Prepare the payload for the AI model
    payload = {
        "model": _model,
        "messages": [
            *MESSAGES,
            {"role": "user", "content": f"Name: {user.name}"},
            {"role": "user", "content": prompt},
        ],
    }

    # Initialize empty response text
    status_code = 200
    generated_text = ""
    try:
        # Make API request to the AI model
        response = requests.post(
            url=LLM_URL,
            data=json.dumps(payload),
            headers=headers,
        )
        status_code = response.status_code
        response.raise_for_status()

        # Process successful response
        if response.status_code == 200:
            generated_text = response.json()["choices"][0]["message"]["content"]
            return [generated_text, status_code]
        else:
            # Log error for unsuccessful response
            logger.error(
                f"Failed to get a valid response from DOKOOLA_LLM model: [{_model}]. Status code: {response.status_code}"
            )
            return [None, status_code]

    except requests.exceptions.RequestException as e:
        # Log any request-related errors
        logger.error(f"Error calling AI model [{_model}]: {e}")
        if max_retries > 0:
            max_retries = max_retries - 1
            logger.info(
                f"Retrying completion with AI model [{MODELS[current_model_index+1]}]"
            )
            return get_ai_text_completion(
                prompt=prompt, user=user, max_retries=max_retries
            )

        return [None, status_code]
