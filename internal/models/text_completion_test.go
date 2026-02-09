package models

import (
	"testing"
)

func TestTextCompletionRequest(t *testing.T) {
	t.Run("creates request with text", func(t *testing.T) {
		req := TextCompletionRequest{
			Text: "Generate a bio",
		}

		if req.Text != "Generate a bio" {
			t.Errorf("expected text 'Generate a bio', got %q", req.Text)
		}
	})

	t.Run("handles empty text", func(t *testing.T) {
		req := TextCompletionRequest{
			Text: "",
		}

		if req.Text != "" {
			t.Errorf("expected empty text, got %q", req.Text)
		}
	})

	t.Run("handles long text", func(t *testing.T) {
		longText := "This is a very long text that contains multiple sentences and paragraphs. It should be stored correctly in the request."
		req := TextCompletionRequest{
			Text: longText,
		}

		if req.Text != longText {
			t.Errorf("expected long text to be preserved")
		}
	})
}

func TestTextCompletionResponse(t *testing.T) {
	t.Run("success response", func(t *testing.T) {
		completion := "Generated completion text"
		resp := TextCompletionResponse{
			Success:      true,
			Completion:   &completion,
			ErrorMessage: nil,
		}

		if !resp.Success {
			t.Error("expected success to be true")
		}

		if resp.Completion == nil || *resp.Completion != completion {
			t.Error("expected completion text")
		}

		if resp.ErrorMessage != nil {
			t.Error("expected error message to be nil")
		}
	})

	t.Run("error response", func(t *testing.T) {
		errorMsg := "Failed to complete"
		resp := TextCompletionResponse{
			Success:      false,
			Completion:   nil,
			ErrorMessage: &errorMsg,
		}

		if resp.Success {
			t.Error("expected success to be false")
		}

		if resp.Completion != nil {
			t.Error("expected completion to be nil")
		}

		if resp.ErrorMessage == nil || *resp.ErrorMessage != errorMsg {
			t.Error("expected error message")
		}
	})

	t.Run("response with both fields nil", func(t *testing.T) {
		resp := TextCompletionResponse{
			Success:      false,
			Completion:   nil,
			ErrorMessage: nil,
		}

		if resp.Completion != nil {
			t.Error("expected completion to be nil")
		}

		if resp.ErrorMessage != nil {
			t.Error("expected error message to be nil")
		}
	})
}
