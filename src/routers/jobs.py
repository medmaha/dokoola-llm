import os
import re

import requests
from fastapi import APIRouter, Response

from src.ai import get_ai_text_completion
from src.logger import Logger
from src.routers.models.jobs import (
    JobCategorizationRequest,
    JobCategorizationResponse,
    JobCategory,
    JobResponseData,
)
from src.users import get_user_by_public_id

LLM_URL = os.getenv("LLM_URL", None)
LLM_API_KEY = os.getenv("LLM_API_KEY", None)

logger = Logger(__name__)
router = APIRouter(tags=["Jobs Management Endpoint"])


@router.post(
    "/jobs/categorization/",
    response_model=JobCategorizationResponse,
    description="Categorizes a list of jobs.",
    responses={
        200: {"detail": "Job categorization finished successfully"},
        400: {"detail": "Invalid request body"},
        404: {"detail": "User not found"},
        500: {"detail": "Internal server error"},
    },
)
def job_categorization(request: JobCategorizationRequest, response: Response):
    try:
        data = []
        categories = get_categories()

        for job in request.data:
            prompt = build_prompt(job.description, categories)

            [completion, status_code] = get_ai_text_completion(prompt, user=None)

            if not completion:
                logger.error(
                    f"job_categorization request failed for user {job.public_id} with status code {status_code}"
                )
            else:
                completion = re.sub(r"<.*?>", "", completion).strip().lower()
                for category in categories:
                    if category.slug == completion:
                        data.append(
                            JobResponseData(
                                public_id=job.public_id, category=category.slug
                            )
                        )
                        break

        logger.info(
            f"job_categorization request completed with status code {response.status_code}"
        )
        return JobCategorizationResponse(data=data, success=True)
    except Exception as e:
        logger.error(f"job_categorization request failed: {e}")
        response.status_code = 500
        return JobCategorizationResponse(
            success=False, error_message="Internal server error"
        )
    finally:
        logger.info(
            f"job_categorization request completed with status code {response.status_code}"
        )


categories = None
CATEGORIES_API_URL = os.getenv("CATEGORIES_API_URL")


def get_categories():
    global categories

    if categories:
        return categories

    response = requests.get(CATEGORIES_API_URL)
    if response.status_code != 200:
        logger.error(
            f"categories request failed with status code {response.status_code}"
        )
        return None
    data = response.json()
    categories = [JobCategory(**category) for category in data]
    return categories


def build_prompt(job_description: str, categories: list[JobCategory]):
    prompt = f"""
    You are a job categorization assistant. Your task is to match jobs to predefined categories.
    
    IMPORTANT: You must ONLY respond with the exact category slug from the list. Do not add any explanation or additional text.
    
    Job Description: {job_description}
    
    Available category slugs: {[category.slug for category in categories]}
    
    Rules:
    - Return ONLY the category slug that best matches the job description
    - If no category matches, return exactly "other"
    - Do not add any other text or explanation
    - The response must be lowercase
    """
    return prompt
