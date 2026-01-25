import os

from dotenv import load_dotenv

load_dotenv()


class Settings:
    """Application settings and configuration."""

    # Application
    app_name: str = os.getenv("APP_NAME", "Dokoola LLM Service")
    app_version: str = os.getenv("APP_VERSION", "0.1.0")
    debug: bool = os.getenv("DEBUG", "true").lower() == "true"

    # API
    api_prefix: str = os.getenv("API_PREFIX", "/api/v1")
    host: str = os.getenv("HOST", "0.0.0.0")
    port: int = int(os.getenv("PORT", "8080"))


settings = Settings()
