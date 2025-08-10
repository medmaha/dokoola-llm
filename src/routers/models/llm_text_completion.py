from typing import Optional

from pydantic import BaseModel, Field


class TextCompletionRequest(BaseModel):
    text: str = Field(default="", description="The text to complete.")


class TextCompletionResponse(BaseModel):
    completion: Optional[str] = Field(default=None, description="The completed text.")
    error_message: Optional[str] = Field(default=None, description="The error message.")
    success: bool = Field(
        default=True, description="Whether the request was successful."
    )
