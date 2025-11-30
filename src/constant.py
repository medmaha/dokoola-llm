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

MODELS = [
    
    "qwen/qwen3-30b-a3b:free",
	## "microsoft/mai-ds-r1:free",
	"openai/gpt-oss-20b:free",
	# "tencent/hunyuan-a13b-instruct:free",
	## "google/gemini-2.5-flash-lite",
	"mistralai/mistral-nemo:free",
	# "qwen/qwen3-235b-a22b:free",
	# "qwen/qwen-2.5-72b-instruct:free",
	# "google/gemini-2.5-flash-lite",
	# "qwen/qwen3-coder:free",
	"microsoft/mai-ds-r1:free",
	# "google/gemini-2.5-flash-lite",
	## "tngtech/deepseek-r1t2-chimera:free",
	"z-ai/glm-4.5-air:free",
	"nvidia/nemotron-nano-12b-v2-vl:free",
]

MESSAGES = [
    {"role": "system", "content": "Welcome to Dokoola Platform!"},
    {
        "role": "system",
        "content": "A digital marketplace connecting talent to opportunities across The Gambia and soon Africa",
    },
    {"role": "system", "content": "You are an expert career assistant."},
    {
        "role": "system",
        "content": "Avoid using none standard characters and tags. Do not wrap your responses in <answer> tags or any other XML/HTML-like tags (Unless specified).",
    },
]


ALLOWED_ORIGINS = get_allowed_origins()
