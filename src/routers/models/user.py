from pydantic import BaseModel, HttpUrl, Field
from typing import Optional


class UserProfile(BaseModel):
    avatar: Optional[HttpUrl] = Field(default=None)
    is_staff: Optional[bool] = Field(default=False)
    is_talent: Optional[bool] = Field(default=False)
    is_client: Optional[bool] = Field(default=False)

class AuthUser(UserProfile):
    name: str
    email: str
    public_id: str
    is_active: bool
    email_verified: bool
    username: Optional[str] = Field(default=None)
    complete_profile: Optional[bool] = Field(default=True)

    class Config:
        from_attributes = True
