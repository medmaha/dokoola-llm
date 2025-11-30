import os

from src.logger import Logger
from src.routers.models.user import AuthUser
from cerebras.cloud.sdk import Cerebras

logger = Logger(__name__)

LLM_API_KEY = os.getenv("LLM_API_KEY", None)

llm_client = Cerebras(
    api_key=LLM_API_KEY
)

MESSAGES = [
    {
        "role": "system",
        "content": "You are a helpful assistant."
    },
    {"role": "system", "content": "Welcome to Dokoola Platform!"},
    {
        "role": "system",
        "content": "A digital marketplace connecting talent to opportunities across The Gambia and soon Africa",
    },
    {"role": "system", "content": "You are an expert career assistant."},
]


def engage_llm(prompt: str, user: AuthUser, _messges: list[dict]=[]) -> str | None:
    try:
        if user:
            _messges.append({"role": "user", "content": f"Name: {user.name}"})
            
        payload = {
            "top_p": 1,
            "stream": False,
            "temperature": 0.2,
            "model": "llama-3.3-70b",
            "max_completion_tokens": 1024,
            "messages": [
                *MESSAGES,
                *_messges,
                {"role": "user", "content": prompt},
            ],
        }
        completion = llm_client.chat.completions.create(
            top_p=payload["top_p"],
            model=payload["model"],
            stream=payload["stream"],
            messages=payload["messages"],
            temperature=payload["temperature"],
            max_completion_tokens=payload["max_completion_tokens"],
        )

        return completion.choices[0].message.content

    except Exception as e:
        logger.error(f"Error calling AI: {e}")
        return None
