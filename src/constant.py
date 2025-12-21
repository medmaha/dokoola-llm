import os
from typing import TypedDict
from configparser import ConfigParser
from pathlib import Path


class ServiceConfig(TypedDict):
    host: str
    client_name: str
    secret_hash: str


def load_services_from_config() -> dict[str, ServiceConfig]:
    """Load allowed services from config.ini file."""
    config = ConfigParser()
    config_path = Path(__file__).parent.parent / "config.ini"

    if not config_path.exists():
        raise FileNotFoundError(f"config.ini not found at {config_path}")

    config.read(config_path)
    services = {}

    for section in config.sections():
        if section.startswith("SERVICE_"):
            service_key = section.replace("SERVICE_", "")
            services[service_key] = {
                "host": config.get(section, "host"),
                "client_name": config.get(section, "client_name"),
                "secret_hash": config.get(section, "secret_hash"),
            }

    return services


def get_allowed_origins():
    return [service["host"] for service in ALLOWED_SERVICES.values()]


DOKOOLA_X_LLM_SERVICE_KEY_NAME = os.getenv("DOKOOLA_X_LLM_SERVICE_KEY_NAME")
DOKOOLA_X_LLM_SERVICE_CLIENT_NAME = os.getenv("DOKOOLA_X_LLM_SERVICE_CLIENT_NAME")
DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME = os.getenv(
    "DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME"
)

ALLOWED_SERVICES = load_services_from_config()


class MessageType(TypedDict):
    role: str
    content: str


MESSAGES: list[MessageType] = [
    {"role": "system", "content": "Welcome to Dokoola Platform!"},
    {
        "role": "system",
        "content": """
You are an expert assistant embedded in Dokoola, a premium talent marketplace that connects top-tier freelancers with high-value clients.

Your role is to help both talents and clients communicate clearly, professionally, and persuasively in order to win more work and build lasting relationships.

Core principles you never break:
- Always write natural, human-sounding text that feels personal and authentic.
- Prioritize clarity, brevity, and impact.
- Never use generic filler, corporate jargon, or overused buzzwords unless the user explicitly requests them.
- Never reveal you are an AI or that the text was generated.
- Never add explanations, markdown formatting, or extra commentary unless specifically asked.
- Strictly respect any length limits (count spaces).
- When writing for talents: first person, confident but never arrogant.
- When writing for clients/companies: third person, credible and professional.
- Output only the final text the user needs — nothing else.
- Only excuted the task you're asked, do not initiate any conversation/chat

You are trusted to make every Dokoola user sound like the best version of themselves.
""",
    },
]


ALLOWED_ORIGINS = get_allowed_origins()
