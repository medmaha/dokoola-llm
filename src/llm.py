import os

from cerebras.cloud.sdk import Cerebras

from src.constant import MESSAGES
from src.logger import Logger
from src.routers.models.user import AuthUser

logger = Logger(__name__)

LLM_API_KEY = os.getenv("LLM_API_KEY", None)

llm_client = Cerebras(api_key=LLM_API_KEY)


def engage_llm(prompt: str, user: AuthUser, _messges: list[dict] = []) -> str | None:
    try:
        if user:
            _messges.append({"role": "user", "content": f"Name: {user.name}"})

        completion = llm_client.chat.completions.create(
            messages=[
                *MESSAGES,
                *_messges,
                {"role": "user", "content": prompt},
            ],
            model="zai-glm-4.6",
            stream=False,
            max_completion_tokens=40960,
            temperature=0.6,
            top_p=0.95,
        )

        return completion.choices[0].message.content

    except Exception as e:
        logger.error(f"Error calling AI: {e}")
        return None
