from typing import Optional
from enum import Enum
from pydantic import BaseModel, Field


class PromptTemplateEnum(str, Enum):
    NONE = "none"
    TALENT_BIO = "talent_bio"
    EMPLOYER_ABOUT_US = "client_about_us"
    #
    JOB_DESCRIPTION = "job_description"
    #
    PROPOSAL_COVER_LETTER = "proposal_cover_letter"


class PromptGenerationRequest(BaseModel):
    data: dict = Field(
        default={}, description="The data required to generate the prompt."
    )
    template_name: PromptTemplateEnum = Field(
        default=PromptTemplateEnum.NONE,
        description="The name of the prompt template to use.",
    )


class PromptGenerationResponse(BaseModel):
    completion: Optional[str] = Field(
        default=None, description="The generated prompt text."
    )
    error_message: Optional[str] = Field(default=None, description="The error message.")
    success: bool = Field(
        default=True, description="Indicates if the prompt generation was successful."
    )
