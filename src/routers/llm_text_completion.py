from fastapi import APIRouter, Response

from src.llm import engage_llm
from src.logger import Logger
from src.routers.models.llm_text_completion import (
    TextCompletionRequest,
    TextCompletionResponse,
)
from src.users import get_user_by_public_id

logger = Logger(__name__)
router = APIRouter(tags=["LLM Text Completion"])


@router.post(
    "/chat/completion/{userPublicId}",
    response_model=TextCompletionResponse,
    description="Completes a given text using a language model.",
    responses={
        200: {"detail": "Text completion finished successfully"},
        400: {"detail": "Invalid request body"},
        404: {"detail": "User not found"},
        500: {"detail": "Internal server error"},
    },
)
def llm_text_completion(
    request: TextCompletionRequest, userPublicId: str, respone: Response
):
    try:
        user = get_user_by_public_id(userPublicId)
        if not user:
            respone.status_code = 404
            return TextCompletionResponse(success=False, error_message="User not found")

        completion = engage_llm(request.text, user)
        if not completion:
            respone.status_code = 429
            return TextCompletionResponse(
                success=False, error_message="Failed to generate completion"
            )

        return TextCompletionResponse(completion=completion, success=True)
    except Exception as e:
        logger.error(f"llm_text_completion request failed: {e}")
        respone.status_code = 500
        return TextCompletionResponse(
            success=False, error_message="Internal server error"
        )
    finally:
        logger.info(
            f"llm_text_completion request completed for user {userPublicId} with status code {respone.status_code}"
        )
