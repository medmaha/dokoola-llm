import os
from typing import Any

from cerebras.cloud.sdk import Cerebras 

from src.logger import Logger
from src.routers.models.user import AuthUser

logger = Logger(__name__)

LLM_API_KEY = os.getenv("LLM_API_KEY", None)
llm_client = Cerebras(api_key=LLM_API_KEY)


def engage_llm(prompt: str, user: AuthUser, _messges: list[dict] = []) -> str | None:
    try:
        # if user:
        # _messges.append({"role": "user", "content": f"Name: {user.name}"})

        completion: Any = llm_client.chat.completions.create(
            messages=[
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
                {"role": "user", "content": prompt},
            ],
            model="zai-glm-4.6",
            stream=False,
            max_completion_tokens=40960,
            temperature=0.6,
            top_p=0.95,
        )
        if not completion:
            raise ValueError("No completion returned from LLM.")

        if not completion.choices:
            return ""
        
        return completion.choices[0].message.content
    except Exception as e:
        logger.error(f"Error calling AI: {e}")
        return None

