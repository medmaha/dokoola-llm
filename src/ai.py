import json
import os
import random

import requests

from src.constant import MESSAGES, MODELS
from src.logger import Logger
from src.routers.models.user import AuthUser

LLM_URL = os.getenv("LLM_URL", None)
LLM_API_KEY = os.getenv("LLM_API_KEY", None)

current_model_index = random.randrange(0, len(MODELS) - 1)

logger = Logger(__name__)


def get_ai_text_completion(
    prompt: str, user: AuthUser | None, max_retries=3
) -> tuple[str | None]:

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
    messages = MESSAGES[:]
    if user:
        messages.append({"role": "user", "content": f"Name: {user.name}"})

    payload = {
        "model": _model,
        "messages": [
            *messages,
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
            return (generated_text, status_code)
        else:
            # Log error for unsuccessful response
            logger.error(
                f"Failed to get a valid response from DOKOOLA_LLM model: [{_model}]. Status code: {response.status_code}"
            )
            return (None, status_code)

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

        return (None, status_code)
