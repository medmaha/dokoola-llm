package models

// PromptTemplateEnum defines available prompt templates
type PromptTemplateEnum string

const (
	PromptNone                PromptTemplateEnum = "none"
	PromptTalentBio           PromptTemplateEnum = "talent_bio"
	PromptClientAboutUs       PromptTemplateEnum = "client_about_us"
	PromptJobDescription      PromptTemplateEnum = "job_description"
	PromptProposalCoverLetter PromptTemplateEnum = "proposal_cover_letter"
)

// ModelTuneEnum defines tone options for prompts
type ModelTuneEnum string

const (
	TuneProfessional ModelTuneEnum = "professional"
	TuneConfident    ModelTuneEnum = "confident"
	TuneFriendly     ModelTuneEnum = "friendly"
	TuneEnthusiastic ModelTuneEnum = "enthusiastic"
	TuneFormal       ModelTuneEnum = "formal"
	TuneWarm         ModelTuneEnum = "warm"
)

// ModelResponseLengthEnum defines length options for responses
type ModelResponseLengthEnum string

const (
	LengthShort    ModelResponseLengthEnum = "short"    // 120-180 words
	LengthMedium   ModelResponseLengthEnum = "medium"   // 200-300 words
	LengthDetailed ModelResponseLengthEnum = "detailed" // 320-450 words
)

// PromptGenerationRequest is the request payload for prompt-based generation
type PromptGenerationRequest struct {
	Data         map[string]interface{} `json:"data" binding:"required"`
	TemplateName PromptTemplateEnum     `json:"template_name" binding:"required"`
}

// PromptGenerationResponse is the response for prompt generation
type PromptGenerationResponse struct {
	Completion   *string `json:"completion,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
	Success      bool    `json:"success"`
}
