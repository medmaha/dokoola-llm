package models

// TextCompletionRequest is the request payload for text completion
type TextCompletionRequest struct {
	Text string `json:"text" binding:"required"`
}

// TextCompletionResponse is the response for text completion
type TextCompletionResponse struct {
	Completion   *string `json:"completion,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
	Success      bool    `json:"success"`
}
