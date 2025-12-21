from typing import Callable
from fastapi import APIRouter, Response

from src.llm import engage_llm
from src.llm_prompts.prompts import PromptTemplates
from src.logger import Logger
from src.routers.models.prompts import (
    PromptGenerationRequest,
    PromptGenerationResponse,
    PromptTemplateEnum,
)
from src.users import get_user_by_public_id

logger = Logger(__name__)
router = APIRouter(tags=["LLM Generation Endpoint"])


@router.post(
    "/actions/generation/{userPublicId}",
    response_model=PromptGenerationResponse,
    description="Generates content using a language model based on the provided template and data.",
    responses={
        200: {"detail": "Content generation completed successfully"},
        400: {"detail": "Invalid request body"},
        403: {"detail": "Forbidden: User not found"},
        429: {"detail": "Failed to generate completion"},
        500: {"detail": "Internal server error"},
    },
)
def process_prompt_completion(
    request: PromptGenerationRequest, userPublicId: str, respone: Response
):
    try:
        # Validate user exists
        user = get_user_by_public_id(userPublicId)
        if not user:
            respone.status_code = 403
            return PromptGenerationResponse(
                success=False, error_message="Forbidden: User not found"
            )

        # Validate template name
        if request.template_name not in PromptTemplateEnum.__members__.values():
            return PromptGenerationResponse(
                success=False, error_message="Invalid template name"
            )

        # Get the template function and generate content
        generate_prompt: Callable = PromptTemplates.get_template_function(
            PromptTemplateEnum(request.template_name)
        )

        prompt = generate_prompt(request.data, user)
        if not prompt:
            respone.status_code = 429
            return PromptGenerationResponse(
                success=False, error_message="Failed to generate completion"
            )

        completion = engage_llm(prompt, user)
        return PromptGenerationResponse(
            completion=completion, success=completion is not None
        )

    except Exception as e:
        logger.error(f"Content generation request failed: {e}")
        respone.status_code = 500
        return PromptGenerationResponse(
            success=False, error_message="Internal server error"
        )
    finally:
        logger.info(
            f"Content generation request completed for user {userPublicId} with status code {respone.status_code}"
        )
