import os


def get_allowed_origins():
    return [service["host"] for service in ALLOWED_SERVICES.values()]


DOKOOLA_X_LLM_SERVICE_KEY_NAME = os.getenv("DOKOOLA_X_LLM_SERVICE_KEY_NAME")
DOKOOLA_X_LLM_SERVICE_CLIENT_NAME = os.getenv("DOKOOLA_X_LLM_SERVICE_CLIENT_NAME")
DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME = os.getenv(
    "DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME"
)

ALLOWED_SERVICES = {
    "DKL7f2h8j4k9m5n3p6q1r8s4t7v2w9x5y": {
        "host": "https://dokoola.com",
        "client_name": "DOKOOLA_WEB",
        "secret_hash": "web_8f4a9c2e7b3d6k1m5n9p2r8s4t7v2w9x5",
    },
    "DKL3a9b5c7d2e8f4g6h1j7k4m9n5p2r": {
        "host": "http://129.151.181.32:8000",
        "client_name": "DOKOOLA_BACKEND",
        "secret_hash": "backend_3k9m5n2p7r4s8t1v6w3x9y2z5a4b8c1d",
    },
    "DKL2f9b5c7d2e8f4g6h1j7k4m9n5p2f": {
        "host": "http://129.151.181.32:8001",
        "client_name": "DOKOOLA_AGENT",
        "secret_hash": "agent_5p2r7s4t1v8w3x6y9z2a5b8c4d7e1f3g",
    },
}

MESSAGES = [
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

You are trusted to make every Dokoola user sound like the best version of themselves.
""",
    },
]


ALLOWED_ORIGINS = get_allowed_origins()
