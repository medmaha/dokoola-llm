import os


def get_allowed_origins():
    return [service["host"] for service in ALLOWED_SERVICES.values()]


DOKOOLA_X_LLM_SERVICE_KEY_NAME = os.getenv("DOKOOLA_X_LLM_SERVICE_KEY_NAME")
DOKOOLA_X_LLM_SERVICE_CLIENT_NAME = os.getenv("DOKOOLA_X_LLM_SERVICE_CLIENT_NAME")
DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME = os.getenv("DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME")

ALLOWED_SERVICES = {
    "DKL7f2h8j4k9m5n3p6q1r8s4t7v2w9x5y": {
        "host": "https://dokoola.com",
        "client_name": "DOKOOLA_WEB",
        "secret_hash": "web_8f4a9c2e7b3d6k1m5n9p2r8s4t7v2w9x5"
    },
    "DKL3a9b5c7d2e8f4g6h1j7k4m9n5p2r": {
        "host": "http://129.151.181.32:8000",
        "client_name": "DOKOOLA_BACKEND", 
        "secret_hash": "backend_3k9m5n2p7r4s8t1v6w3x9y2z5a4b8c1d"
    },
    "DKL3a9b5c7d2e8f4g6h1j7k4m9n5p2s": {
        "host": "https://m.dokoola.com",
        "client_name": "DOKOOLA_MOBILE",
        "secret_hash": "mobile_5p2r7s4t1v8w3x6y9z2a5b8c4d7e1f3g"
    }
}

ALLOWED_ORIGNS = get_allowed_origins()
