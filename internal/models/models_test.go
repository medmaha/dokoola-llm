package models

import (
	"testing"
)

func TestPromptTemplateEnumValues(t *testing.T) {
	tests := []struct {
		name     string
		value    PromptTemplateEnum
		expected string
	}{
		{
			name:     "none template",
			value:    PromptNone,
			expected: "none",
		},
		{
			name:     "talent bio template",
			value:    PromptTalentBio,
			expected: "talent_bio",
		},
		{
			name:     "client about us template",
			value:    PromptClientAboutUs,
			expected: "client_about_us",
		},
		{
			name:     "job description template",
			value:    PromptJobDescription,
			expected: "job_description",
		},
		{
			name:     "proposal cover letter template",
			value:    PromptProposalCoverLetter,
			expected: "proposal_cover_letter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.value))
			}
		})
	}
}

func TestModelTuneEnumValues(t *testing.T) {
	tests := []struct {
		name     string
		value    ModelTuneEnum
		expected string
	}{
		{name: "professional", value: TuneProfessional, expected: "professional"},
		{name: "confident", value: TuneConfident, expected: "confident"},
		{name: "friendly", value: TuneFriendly, expected: "friendly"},
		{name: "enthusiastic", value: TuneEnthusiastic, expected: "enthusiastic"},
		{name: "formal", value: TuneFormal, expected: "formal"},
		{name: "warm", value: TuneWarm, expected: "warm"},
		{name: "persuasive", value: TunePersuasive, expected: "persuasive"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.value))
			}
		})
	}
}

func TestModelResponseLengthEnumValues(t *testing.T) {
	tests := []struct {
		name     string
		value    ModelResponseLengthEnum
		expected string
	}{
		{name: "short", value: LengthShort, expected: "short"},
		{name: "medium", value: LengthMedium, expected: "medium"},
		{name: "detailed", value: LengthDetailed, expected: "detailed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.value))
			}
		})
	}
}

func TestPromptGenerationRequest(t *testing.T) {
	t.Run("creates request with valid data", func(t *testing.T) {
		data := map[string]interface{}{
			"profile": map[string]interface{}{
				"name":  "John Doe",
				"title": "Developer",
			},
		}

		req := PromptGenerationRequest{
			Data:         data,
			TemplateName: PromptTalentBio,
		}

		if req.TemplateName != PromptTalentBio {
			t.Errorf("expected template %q, got %q", PromptTalentBio, req.TemplateName)
		}

		if req.Data["profile"] == nil {
			t.Error("expected profile data")
		}
	})
}

func TestPromptGenerationResponse(t *testing.T) {
	t.Run("success response with completion", func(t *testing.T) {
		completion := "Generated bio text"
		resp := PromptGenerationResponse{
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
		errorMsg := "Failed to generate"
		resp := PromptGenerationResponse{
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
}

func TestAuthUser(t *testing.T) {
	t.Run("creates auth user with valid fields", func(t *testing.T) {
		user := AuthUser{
			Name:            "Jane Doe",
			PublicID:        "user-123",
			IsActive:        true,
			CompleteProfile: true,
			IsStaff:         false,
			IsTalent:        true,
			IsClient:        false,
		}

		if user.Name != "Jane Doe" {
			t.Errorf("expected name Jane Doe, got %s", user.Name)
		}

		if user.PublicID != "user-123" {
			t.Errorf("expected public_id user-123, got %s", user.PublicID)
		}

		if !user.IsActive {
			t.Error("expected user to be active")
		}

		if !user.IsTalent {
			t.Error("expected user to be talent")
		}
	})
}

func TestJobData(t *testing.T) {
	t.Run("creates job data with required fields", func(t *testing.T) {
		job := JobData{
			PublicID:    "job-456",
			Description: "Seeking experienced developer",
		}

		if job.PublicID != "job-456" {
			t.Errorf("expected public_id job-456, got %s", job.PublicID)
		}

		if job.Description == "" {
			t.Error("expected non-empty description")
		}
	})
}

func TestJobResponseData(t *testing.T) {
	t.Run("creates response data with category", func(t *testing.T) {
		resp := JobResponseData{
			PublicID: "job-456",
			Category: "web-development",
		}

		if resp.PublicID != "job-456" {
			t.Errorf("expected public_id job-456, got %s", resp.PublicID)
		}

		if resp.Category != "web-development" {
			t.Errorf("expected category web-development, got %s", resp.Category)
		}
	})
}

func TestJobCategorizationRequest(t *testing.T) {
	t.Run("creates request with multiple jobs", func(t *testing.T) {
		jobs := []JobData{
			{PublicID: "job-1", Description: "Job 1"},
			{PublicID: "job-2", Description: "Job 2"},
		}

		req := JobCategorizationRequest{
			Data: jobs,
		}

		if len(req.Data) != 2 {
			t.Errorf("expected 2 jobs, got %d", len(req.Data))
		}
	})
}

func TestJobCategorizationResponse(t *testing.T) {
	t.Run("success response with results", func(t *testing.T) {
		results := []JobResponseData{
			{PublicID: "job-1", Category: "web-dev"},
			{PublicID: "job-2", Category: "mobile-dev"},
		}

		resp := JobCategorizationResponse{
			Success: true,
			Data:    results,
		}

		if !resp.Success {
			t.Error("expected success to be true")
		}

		if len(resp.Data) != 2 {
			t.Errorf("expected 2 results, got %d", len(resp.Data))
		}
	})

	t.Run("error response", func(t *testing.T) {
		errorMsg := "Failed to categorize"
		resp := JobCategorizationResponse{
			Success:      false,
			Data:         nil,
			ErrorMessage: &errorMsg,
		}

		if resp.Success {
			t.Error("expected success to be false")
		}

		if resp.ErrorMessage == nil {
			t.Error("expected error message")
		}
	})
}

func TestJobCategory(t *testing.T) {
	t.Run("creates category with parent", func(t *testing.T) {
		parentSlug := "development"
		parentDesc := "Development related"

		cat := JobCategory{
			Slug:              "web-development",
			Description:       "Web development category",
			ParentSlug:        &parentSlug,
			ParentDescription: &parentDesc,
		}

		if cat.Slug != "web-development" {
			t.Errorf("expected slug web-development, got %s", cat.Slug)
		}

		if cat.ParentSlug == nil || *cat.ParentSlug != parentSlug {
			t.Error("expected parent slug")
		}
	})

	t.Run("creates category without parent", func(t *testing.T) {
		cat := JobCategory{
			Slug:        "design",
			Description: "Design category",
		}

		if cat.ParentSlug != nil {
			t.Error("expected parent slug to be nil")
		}

		if cat.ParentDescription != nil {
			t.Error("expected parent description to be nil")
		}
	})
}
