from typing import TypedDict
from datetime import datetime

from src.routers.models.user import AuthUser

from .models import (
    JobDescriptionPromptModel,
    ProposalCoverLetterPromptModel,
    TalentInterface,
)
from .models import EmployerModel
from src.routers.models.prompts import (
    PromptTemplateEnum,
)

tone_descriptions = {
    "professional": "professional, polished, and confident",
    "confident": "bold, self-assured, and results-oriented",
    "friendly": "warm, approachable, and personable",
    "enthusiastic": "energetic, passionate, and excited",
    "formal": "highly formal, corporate, and respectful",
    "warm": "friendly yet professional with a personal touch",
}

length_guidelines = {
    "short": "Keep it concise: 120-180 words (3-4 short paragraphs)",
    "medium": "Balanced length: 200-300 words (4-5 paragraphs)",
    "detailed": "In-depth and thorough: 320-450 words (5-7 paragraphs with specific examples)",
}


class TalentBioParams(TypedDict):
    profile: dict


def get_client_short_about_prompt(profile: dict) -> str:
    """Generate a prompt for creating a client's 'About Us' section."""
    company = profile.get("company", {})
    name = company.get("name") or profile.get("name", "")
    industry = company.get("industry", "technology")

    country = ""
    if company.get("country"):
        country = company["country"].get("name", "")
    elif profile.get("user", {}).get("country"):
        country = profile["user"]["country"].get("name", "")

    stats = profile.get("stats", {})
    total_spend = stats.get("total_spend")
    spend = f"${total_spend:,}" if total_spend else None

    jobs_posted = stats.get("jobs_posted_count", 0)

    country_text = f" · {country}" if country else ""
    spend_text = f"Total spend: {spend}+" if spend else ""

    return f"""
Write a sharp, trust-building "About Us" blurb (third person) for this client.

Company: {name}
Industry: {industry}{country_text}
Jobs posted: {jobs_posted}
{spend_text}

Requirements:
- Max 700 characters (including spaces)
- Professional & credible
- Shows they're serious about hiring top talent
- No generic fluff

Just output the final text. Nothing else.
""".strip()


# =========================================================================
class TalentBiographyParams(TypedDict):
    profile: dict


def get_talent_short_bio_prompt(data: TalentBiographyParams, user: AuthUser) -> str:
    """Generate a prompt for creating a talent's bio."""

    profile = TalentInterface(**data["profile"])

    name = profile.name
    title = profile.title
    skills = profile.skills
    rating = profile.rating
    jobs_completed = profile.jobs_completed
    pricing = profile.pricing
    badge = profile.badge

    if isinstance(skills, str):
        skills_list = ", ".join([s.strip() for s in skills.split(",")][:5])
    else:
        skills_list = ""

    return f"""
You are an elite Dokoola profile copywriter.

Write a powerful, first-person bio for {name} in MAX 500 characters (including spaces).

Current title: {title}
Badge: {badge}
Rating: {rating:.1f}/5
Jobs completed: {jobs_completed}
Rate: ${pricing}/hr
Top skills: {skills_list}

Rules:
- First person
- Zero filler words
- Instantly shows expertise + results
- Ends with a hook
- Under 500 chars total

Just output the bio. Nothing else.
""".strip()


# =========================================================================
class CoverLetterParams(TypedDict):
    talent: dict
    job: dict
    resume: dict
    metadata: dict


def get_proposal_cover_letter_prompt(
    data: dict,
    user: AuthUser,
) -> str:
    """Generate a prompt for creating a proposal cover letter."""

    payload = ProposalCoverLetterPromptModel(**data)

    talent = payload.talent
    job = payload.job
    resume = payload.resume
    metadata = payload.metadata

    additional_notes = metadata.additional_notes
    length = metadata.length if metadata else "medium"
    tone = metadata.tone if metadata else "professional"

    client_name = (
        job.third_party_metadata.get("company_name")
        if job.third_party_metadata
        else job.client.name
    )

    job_type = job.job_type
    if job_type == "other" and job.job_type_other:
        job_type = job.job_type_other

    experience_level = job.experience_level
    if experience_level == "other" and job.experience_level_other:
        experience_level = job.experience_level_other

    required_skills = ", ".join(job.required_skills) or "relevant skills"
    duration = job.estimated_duration
    # benefits = job.benefits or []

    # Talent highlights
    rating = talent.rating or 0.5
    rating_stars = f"rated {rating:.1f}/5" if rating else "experienced talent"

    badge = talent.badge or ""
    badge_text = (
        "Pro"
        if badge == "pro"
        else "Top-Rated Star"
        if badge == "star"
        else "verified talent"
    )

    job_duration_text = (
        f"Project Duration: {job.estimated_duration}" if job.estimated_duration else ""
    )
    job_type_text = f"Job Type: {job_type}" if job_type else ""
    job_description_text = (
        f'Full Job Description: """{job.description}"""'
        if job.description
        else "No job description provided."
    )
    # job_benefits_text = f"Job Benefits Offered: {', '.join(benefits)}" if benefits else ""
    resume_text = (
        f"Relevant Experience (from resume): {resume.get('description', '')}"
        if resume
        else ""
    )
    bio_summary_text = f"Bio Summary:: {talent.bio}" if talent.bio else ""
    skills_text = f"Key Skills: {talent.skills}" if talent.skills else ""
    pricing_text = f"Pricing Rate: ${talent.pricing}" if talent.pricing else ""
    additional_text = (
        f"Additional instructions from me: {additional_notes}"
        if additional_notes
        else ""
    )
    if job.third_party_metadata and (
        company_name := job.third_party_metadata.get("company_name")
    ):
        client_name_text = f"Company Name: {company_name}"
    else:
        client_name_text = f"Client Name: {job.client.name}"

    return f"""
You are an expert freelance proposal writer for (Dokoola) who has helped hundreds of freelancers win high-value contracts on platforms like Upwork, Fiverr, and Toptal.

Write a compelling, personalized cover letter (proposal) for the following freelancer applying to this job.

=== FREELANCER PROFILE ===
Name: {talent.name}
Title: {talent.title or ""}
Badge: {badge_text}
Rating: {rating_stars}
{skills_text}
{bio_summary_text}
{pricing_text}
{resume_text}

=== JOB DETAILS ===
Client/Company: {client_name}
{client_name_text}
Job Title: {job.title}
{job_type_text}
Job Category: {job.category.name}
Experience Level Expected: {experience_level}
Required Skills: {required_skills}
Project Duration: {duration}
{job_duration_text}
{job_description_text}

=== INSTRUCTIONS ===
Write a winning proposal cover letter in first person as {talent.name}.

Structure:
1. Strong opening: Greet the client and express genuine interest in THEIR specific project (reference something unique from the job description).
2. Prove fit: Explain why I am the perfect match (highlight overlapping skills, past results, and relevant experience from my profile/resume).
3. Build trust: Mention my rating, badge, and success record.
4. Show understanding: Demonstrate that I fully understand the project goals and challenges.
5. Call to action: End with confidence and invite next steps (interview, questions, etc.).

Tone: {tone}
Length: {length}
Style: Natural, human, engaging — never robotic or generic. Avoid this chars "—" or clichés like "I am passionate about" unless it feels authentic.

{additional_text}

Do NOT mention that this was AI-generated.
Do NOT say "As an AI language model".
Just write the cover letter — nothing else.
Format your response strictly in rich-text using (basic-html-tags, e.g <p>, <strong> <ul,ol,li>, <br/>).
""".strip()


# =========================================================================


class DescriptionParams(TypedDict):
    job: dict
    client: dict
    metadata: dict


def get_job_description_prompt(data: DescriptionParams, user: AuthUser) -> str:
    """Generate a prompt for creating a job description."""

    job = JobDescriptionPromptModel(**data.get("job"))
    client = EmployerModel(**data.get("client"))

    tone = "professional"
    length = "medium"

    if not job or not client:
        raise ValueError("Invalid data provided.")

    is_third_party = job.is_third_party
    third_party = job.third_party_metadata or {}

    if is_third_party and third_party.get("company_name"):
        company_name = third_party["company_name"]
    else:
        company_name = client.name

    # Budget info
    pricing = job.pricing
    if pricing and pricing.budget and pricing.currency and pricing.currency.symbol:
        currency_symbol = pricing.currency.symbol
        project_budget = pricing.budget
        if pricing.fixed_price:
            budget_info = (
                f"{currency_symbol}{project_budget} fixed-price"
                if project_budget
                else "Fixed-price project"
            )
        else:
            budget_info = (
                f"{currency_symbol}{project_budget}/hr"
                if project_budget
                else "Hourly rate (budget negotiable)"
            )
    else:
        budget_info = None

    job_type = job.job_type
    if job_type == "other" and job.job_type_other:
        job_type = job.job_type_other

    duration = job.estimated_duration

    deadline = "Open until filled"

    if job.application_deadline:
        date_obj = datetime.fromisoformat(
            job.application_deadline.replace("Z", "+00:00")
        )
        deadline = f"Apply by {date_obj.strftime('%B %d, %Y')}"

    skills = " · ".join(job.required_skills or []) or "Relevant skills required"

    country = job.country
    address = job.address
    if country and country.name:
        if address:
            location = f"{country.name} - {address}"
        else:
            location = country.name
    else:
        location = address

    third_party_text = " (posted via Dokoola)" if is_third_party else ""

    return f"""
Rewrite this job posting into a clear, professional, and attractive job description that top talent actually wants to apply to.

=== BASIC INFO ===
Title: {job.title or "Untitled Role"}
Category: {job.category or "Not specified"}
Company: {company_name}{third_party_text}

=== DETAILS FROM CLIENT ===

{f"Budget: {budget_info}" if budget_info else ""}
{f"Job Type: {job_type}" if job_type else ""}
{f"Estimated Duration: {duration}" if duration else ""}
{f"Application Deadline: {deadline}" if deadline else ""}
Required Skills: {skills}
{f"Location Preference: {location}" if location else ""}

=== INSTRUCTIONS ===
Write a polished, engaging job description that:
- Starts with a strong, specific openisng line (not "We are looking for...")
- Clearly explains the project/role and its impact
- Highlights what success looks like
- Mentions budget, timeline, and type upfront (no hiding)
- Lists required skills cleanly
- Ends with a confident, welcoming call-to-action
- Format your response strictly in rich-text using (basic-html-tags).

Tone: {tone}, direct, respectful of freelancers' time
Length: {length} 250-450 words max
Style: Human, concise, zero fluff
    
Just output the final job description. No titles, no markdown, no quotes, no extra text.
""".strip()


def get_none(*args, **kwargs) -> str:
    return ""


class PromptTemplates:
    template_map = {
        PromptTemplateEnum.TALENT_BIO: get_talent_short_bio_prompt,
        PromptTemplateEnum.JOB_DESCRIPTION: get_job_description_prompt,
        PromptTemplateEnum.EMPLOYER_ABOUT_US: get_client_short_about_prompt,
        PromptTemplateEnum.PROPOSAL_COVER_LETTER: get_proposal_cover_letter_prompt,
        PromptTemplateEnum.NONE: get_none,
    }

    @staticmethod
    def get_template_function(template_name: PromptTemplateEnum):
        return PromptTemplates.template_map[template_name]
