package models

// JobCategory represents a job category from the backend
type JobCategory struct {
	Slug              string  `json:"slug"`
	Description       string  `json:"description"`
	ParentSlug        *string `json:"parent_slug,omitempty"`
	ParentDescription *string `json:"parent_description,omitempty"`
}

// JobData represents input job data for categorization
type JobData struct {
	PublicID    string `json:"public_id" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// JobResponseData represents categorized job output
type JobResponseData struct {
	PublicID string `json:"public_id"`
	Category string `json:"category"`
}

// JobCategorizationRequest is the request payload for job categorization
type JobCategorizationRequest struct {
	Data []JobData `json:"data" binding:"required,dive"`
}

// JobCategorizationResponse is the response for job categorization
type JobCategorizationResponse struct {
	Data         []JobResponseData `json:"data"`
	ErrorMessage *string           `json:"error_message,omitempty"`
	Success      bool              `json:"success"`
}
