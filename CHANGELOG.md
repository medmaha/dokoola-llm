## [2026.02.08] - 2026-02-08

- fix(ci): update CI configuration to simplify paths, enhance coverage reporting, and adjust Docker build context (e825bbb - medmaha)
- fix(ci): remove coverage report uploads and lint job, update job dependencies (b44b004 - medmaha)
- fix(ci): comment out golangci-lint installation and execution steps, add placeholder message (f7bc054 - medmaha)
- fix(ci): remove unnecessary defaults section from lint job (168c472 - medmaha)
- feat(ci): add CI/CD workflow for testing, linting, changelog generation, and deployment (c9f1858 - medmaha)
- fix(Dockerfile): remove config example copy to streamline image (9d99208 - medmaha)
- Add handlers for health check, job categorization, prompt generation, and text completion (cbc67ce - medmaha)
- feat(golang): refactored codebase to use golang instead of python. for improved performance and maintainability (fd9325f - medmaha)
- fix: update application port to 8000 and refactor health check endpoint (78cc62c - medmaha)
- fix: Update health check interval in Dockerfile for improved monitoring (301f4de - medmaha)
- fix(Dockerfile): restore health check command and remove docker-compose.yaml (3fcfdf2 - medmaha)
- fix(.gitignore): add .vscode directory to ignore list (d00370e - medmaha)
- feat(config): add configuration settings for application and API (e902100 - medmaha)
- fix(models, prompts): add about field to EmployerModel and include it in job description prompt (33ccf90 - medmaha)
- fix(models): update JobDescriptionPromptModel to use Any for category and add _category_name method (7fb18a6 - medmaha)
- fix(Makefile): remove health check step from rebuild process and streamline cleanup (6a5323c - medmaha)
- fix(models): add ProposalTalentResumeModel and update ProposalCoverLetterPromptModel to use it (5540d47 - inoblend)
- fix(models): change avatar field type from HttpUrl to str in UserProfile model (c024c82 - Touray Mahammed)
- fix(models): update EMPLOYER_ABOUT_US value in PromptTemplateEnum to match client terminology (5e501a6 - Touray Mahammed)
- fix(models): add country field to EmployerModel and update TalentInterface to TalentModel (6f29afb - Touray Mahammed)


All notable changes to the LLM service will be documented here.

