import os
import requests
from src.logger import Logger
from src.routers.models.jobs import (
    JobCategory,
)

logger = Logger(__name__)
BACKEND_SERVER_API = os.getenv("BACKEND_SERVER_API")

categories = None
def get_backend_categories():
    global categories

    if categories:
        return categories
        
    response = requests.get(f"{BACKEND_SERVER_API}/categories?scraper=true")
    if response.status_code != 200:
        logger.error(
            f"categories request failed with status code {response.status_code}"
        )
        return []
    data = response.json()
    categories = [JobCategory(**category) for category in data]
    return categories