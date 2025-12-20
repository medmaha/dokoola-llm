from pydantic import BaseModel, Field
from typing import Literal, Optional
from enum import Enum


class ModelTuneEnums(str, Enum):
    PROFESSIONAL = "professional"
    CONFIDENT = "confident"
    FRIENDLY = "friendly"
    ENTHUSIASTIC = "enthusiastic"
    FORMAL = "formal"
    WARM = "warm"


class ModelResponseLengthEnums(str, Enum):
    SHORT = "short"
    MEDIUM = "medium"
    DETAILED = "detailed"


class PromptMetadataModel(BaseModel):
    tone: Optional[ModelTuneEnums] = Field(default=ModelTuneEnums.PROFESSIONAL)
    length: Optional[ModelResponseLengthEnums] = Field(
        default=ModelResponseLengthEnums.MEDIUM
    )
    additional_notes: Optional[str] = Field(default=None)


class CountryModel(BaseModel):
    code: Optional[str] = Field(default=None)
    name: Optional[str] = Field(default=None)


class CategoryModel(BaseModel):
    slug: str
    name: str


class JobDescriptionPromptModel(BaseModel):
    class JobPricingModel(BaseModel):
        class CurrencyModel(BaseModel):
            code: Optional[str] = Field(default=None)
            symbol: Optional[str] = Field(default=None)
            name: Optional[str] = Field(default=None)

        country: Optional[CountryModel] = Field(default=None)
        currency: Optional[CurrencyModel] = Field(default=None)
        symbol: Optional[str] = Field(default=None)
        budget: Optional[float] = Field(default=None)
        fixed_price: Optional[bool] = Field(default=None)

    title: Optional[str] = Field(default=None)
    category: Optional[str] = Field(default=None)
    address: Optional[str] = Field(default=None)
    country: Optional[CountryModel] = Field(default=None)
    job_type: Optional[str] = Field(default="Flexible")
    job_type_other: Optional[str] = Field(default=None)
    is_third_party: Optional[bool] = Field(default=False)
    pricing: Optional[JobPricingModel] = Field(default=None)
    required_skills: Optional[list[str]] = Field(default=None)
    estimated_duration: Optional[str] = Field(default=None)
    application_deadline: Optional[str] = Field(default=None)
    third_party_metadata: Optional[dict] = Field(default=None)

    metadata: Optional[PromptMetadataModel] = Field(
        default_factory=lambda: PromptMetadataModel()
    )


class EmployerModel(BaseModel):
    class JobCompanyModel(BaseModel):
        name: str
        industry: Optional[str] = Field(default=None)
        date_established: Optional[str] = Field(default=None)

    name: str
    logo: Optional[str] = Field(default=None)
    company: Optional[JobCompanyModel] = Field(default=None)


class TalentModel(BaseModel):
    name: str
    badge: Optional[str] = Field(default=None)
    title: Optional[str] = Field(default=None)
    bio: Optional[str] = Field(default=None)
    skills: Optional[str] = Field(default=None)
    avatar: Optional[str] = Field(default=None)
    rating: Optional[float] = Field(default=None)
    verified: Optional[bool] = Field(default=False)
    country: Optional[CountryModel] = Field(default=None)
    pricing: Optional[float] = Field(default=None)


class JobDetailModel(BaseModel):
    title: Optional[str] = Field(default=None)
    category: CategoryModel
    address: Optional[str] = Field(default=None)
    country: Optional[CountryModel] = Field(default=None)
    job_type: str
    job_type_other: Optional[str] = Field(default=None)
    experience_level: str
    experience_level_other: Optional[str] = Field(default=None)
    is_third_party: bool
    pricing: Optional[JobDescriptionPromptModel.JobPricingModel] = Field(default=None)
    benefits: Optional[list[str]] = Field(default=None)
    required_skills: list[str]
    estimated_duration: Optional[str] = Field(default=None)
    application_deadline: Optional[str] = Field(default=None)
    third_party_metadata: Optional[dict] = Field(default=None)
    description: Optional[str] = Field(default=None)

    client: EmployerModel


class ProposalCoverLetterPromptModel(BaseModel):
    talent: TalentModel
    job: JobDetailModel
    resume: dict
    metadata: PromptMetadataModel = Field(default_factory=lambda: PromptMetadataModel())
