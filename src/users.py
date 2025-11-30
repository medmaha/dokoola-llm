import os

import requests

from src.logger import Logger
from src.routers.models.user import AuthUser

logger = Logger(__name__)
BACKEND_SERVER_API = os.getenv("BACKEND_SERVER_API")


def get_user_by_public_id(public_id: str):
    try:
        response = requests.get(f"{BACKEND_SERVER_API}/users/{public_id}/llm")
        response.raise_for_status()
        user = response.json()
        return AuthUser(**user)
    except Exception as e:
        logger.error(f"Failed to get user by public id {public_id} [Error:] {e}")
