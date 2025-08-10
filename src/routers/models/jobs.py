from typing import List, Optional

from pydantic import BaseModel, Field


class JobCategory(BaseModel):
    slug: str
    description: str
    parent__slug: Optional[str]
    parent__description: Optional[str]


class JobData(BaseModel):
    public_id: str
    description: str


class JobResponseData(BaseModel):
    public_id: str
    category: str


class JobCategorizationRequest(BaseModel):
    data: List[JobData] = Field(default=[], description="The data to categorize.")


class JobCategorizationResponse(BaseModel):
    data: List[JobResponseData] = Field(default=[], description="The categorized data.")
    error_message: Optional[str] = Field(default=None, description="The error message.")
    success: bool = Field(
        default=True, description="Whether the request was successful."
    )
