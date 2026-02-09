package prompts

import (
	"testing"

	"github.com/dokoola/llm-go/internal/models"
)

func TestToneDescriptionsMapExists(t *testing.T) {
	expectedTones := []models.ModelTuneEnum{
		models.TuneProfessional,
		models.TuneConfident,
		models.TuneFriendly,
		models.TuneEnthusiastic,
		models.TuneFormal,
		models.TuneWarm,
		models.TunePersuasive,
	}

	for _, tone := range expectedTones {
		if _, exists := ToneDescriptions[tone]; !exists {
			t.Errorf("tone %q not found in ToneDescriptions map", tone)
		}
	}
}

func TestLengthGuidelinesMapExists(t *testing.T) {
	expectedLengths := []models.ModelResponseLengthEnum{
		models.LengthShort,
		models.LengthMedium,
		models.LengthDetailed,
	}

	for _, length := range expectedLengths {
		if _, exists := LengthGuidelines[length]; !exists {
			t.Errorf("length %q not found in LengthGuidelines map", length)
		}
	}
}

func TestBuildPromptWithNoneTemplate(t *testing.T) {
	data := map[string]interface{}{}
	prompt, err := BuildPrompt(models.PromptNone, data, nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if prompt != "" {
		t.Errorf("expected empty prompt for 'none' template, got %q", prompt)
	}
}

func TestBuildPromptWithUnknownTemplate(t *testing.T) {
	data := map[string]interface{}{}
	prompt, err := BuildPrompt("unknown_template", data, nil)

	if err == nil {
		t.Error("expected error for unknown template")
	}

	if prompt != "" {
		t.Error("expected empty prompt on error")
	}
}

func TestBuildTalentBioPromptWithValidData(t *testing.T) {
	data := map[string]interface{}{
		"profile": map[string]interface{}{
			"name":           "John Doe",
			"title":          "Senior Developer",
			"skills":         "Go, Python, JavaScript",
			"rating":         4.9,
			"jobs_completed": 15,
			"pricing":        75.0,
			"badge":          "Top Rated",
		},
	}

	user := &models.AuthUser{
		Name:     "John Doe",
		PublicID: "user-123",
	}

	prompt, err := BuildPrompt(models.PromptTalentBio, data, user)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}

	// Verify prompt contains key information
	if !contains(prompt, "John Doe") {
		t.Error("expected prompt to contain talent name")
	}

	if !contains(prompt, "Senior Developer") {
		t.Error("expected prompt to contain title")
	}
}

func TestBuildTalentBioPromptMissingProfile(t *testing.T) {
	data := map[string]interface{}{
		"other_field": "value",
	}

	prompt, err := BuildPrompt(models.PromptTalentBio, data, nil)

	if err == nil {
		t.Error("expected error when profile is missing")
	}

	if prompt != "" {
		t.Error("expected empty prompt on error")
	}
}

func TestBuildClientAboutUsPromptWithValidData(t *testing.T) {
	data := map[string]interface{}{
		"profile": map[string]interface{}{
			"name":     "Tech Corp",
			"industry": "Software",
		},
	}

	prompt, err := BuildPrompt(models.PromptClientAboutUs, data, nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}

	if !contains(prompt, "Tech Corp") {
		t.Error("expected prompt to contain company name")
	}
}

func TestBuildClientAboutUsPromptMissingCompany(t *testing.T) {
	data := map[string]interface{}{
		"other_field": "value",
	}

	prompt, err := BuildPrompt(models.PromptClientAboutUs, data, nil)

	if err == nil {
		t.Error("expected error when profile is missing")
	}

	if prompt != "" {
		t.Error("expected empty prompt on error")
	}
}

func TestBuildJobDescriptionPromptWithValidData(t *testing.T) {
	data := map[string]interface{}{
		"job": map[string]interface{}{
			"title":         "Frontend Developer",
			"description":   "Build React apps",
			"category_name": "Web Development",
		},
		"client": map[string]interface{}{
			"name": "Tech Company",
		},
	}

	prompt, err := BuildPrompt(models.PromptJobDescription, data, nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}

	if !contains(prompt, "Frontend Developer") {
		t.Error("expected prompt to contain job title")
	}
}

func TestBuildProposalCoverLetterPromptWithValidData(t *testing.T) {
	data := map[string]interface{}{
		"talent": map[string]interface{}{
			"name":  "Jane Developer",
			"title": "Full Stack Engineer",
		},
		"job": map[string]interface{}{
			"title": "Senior Developer",
		},
	}

	prompt, err := BuildPrompt(models.PromptProposalCoverLetter, data, nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}
}

func TestToneDescriptionsHaveValidContent(t *testing.T) {
	for tone, description := range ToneDescriptions {
		if description == "" {
			t.Errorf("tone %q has empty description", tone)
		}
	}
}

func TestLengthGuidelinesHaveValidContent(t *testing.T) {
	for length, guideline := range LengthGuidelines {
		if guideline == "" {
			t.Errorf("length %q has empty guideline", length)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
