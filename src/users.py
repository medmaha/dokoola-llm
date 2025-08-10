import os

import requests

from src.logger import Logger
from src.routers.models.user import AuthUser

logger = Logger(__name__)
AUTH_SERVER_API = os.getenv("AUTH_SERVER_API")


def get_user_by_public_id(public_id: str):
    try:
        print(f"[AUTH_SERVER_API]: {AUTH_SERVER_API}")
        response = requests.get(f"{AUTH_SERVER_API}/users/auth/{public_id}/")
        response.raise_for_status()
        user = response.json()
        return AuthUser(**user)
    except Exception as e:
        logger.error(f"Failed to get user by public id {public_id} [Error:] {e}")
