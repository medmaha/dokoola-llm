## [2026.02.09] - 2026-02-09

- feat(prompts): enhance talent bio prompt with user details and improved guidelines (46c0295 - medmaha)


- feat(llm): implement rate limiting handling in LLM client and update logging in handlers (3a62e91 - medmaha)


- refactor(ci): simplify formatting and testing steps in CI workflow (ff5f874 - medmaha)
- Add unit tests for constants, handlers, middleware, models, and prompts (fc24d31 - medmaha)


- feat(ci): update Docker image handling and add deployment secrets (c9fcbbc - medmaha)


- feat(config): add ENV field to Settings and load from environment variables (9312cb8 - medmaha)


- feat(prompts): add description for persuasive tone in ToneDescriptions (39c1e39 - medmaha)


- feat(models): add persuasive tune option to ModelTuneEnum (acfd239 - medmaha)


- chore(env): remove example environment configuration file (d745108 - medmaha)


- refactor(ci): update CI workflow to use local actions repository for changelog, build, and deploy steps (5e7722a - medmaha)
- refactor(ci): simplify CI workflow by removing redundant steps and using reusable actions (f8b1c8e - medmaha)


- chore(config): remove example configuration file (01257cf - medmaha)


- fix(ci): handle line endings for SSH key and set permissions for known hosts (8641eb8 - medmaha)


- fix(ci): update Docker image tag to remove 'latest' for consistency (ee59297 - medmaha)


- fix(ci): remove SSH connection test step from deployment workflow (b250aca - medmaha)
- fix(ci): update Go version and Docker image in CI configuration (8aa6ea7 - medmaha)


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

