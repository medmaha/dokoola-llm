package prompts

import (
	"fmt"
	"strings"
	"time"

	"github.com/dokoola/llm-go/internal/models"
)

// ToneDescriptions maps tone enums to descriptions
var ToneDescriptions = map[models.ModelTuneEnum]string{
	models.TuneProfessional: "professional, polished, and confident",
	models.TuneConfident:    "bold, self-assured, and results-oriented",
	models.TuneFriendly:     "warm, approachable, and personable",
	models.TuneEnthusiastic: "energetic, passionate, and excited",
	models.TuneFormal:       "highly formal, corporate, and respectful",
	models.TuneWarm:         "friendly yet professional with a personal touch",
	models.TunePersuasive:   "convincing, compelling, and focused on benefits",

}

// LengthGuidelines maps length enums to guidelines
var LengthGuidelines = map[models.ModelResponseLengthEnum]string{
	models.LengthShort:    "Keep it concise: 120-180 words (3-4 short paragraphs)",
	models.LengthMedium:   "Balanced length: 200-300 words (4-5 paragraphs)",
	models.LengthDetailed: "In-depth and thorough: 320-450 words (5-7 paragraphs with specific examples)",
}

// BuildPrompt builds a prompt based on the template name and data
func BuildPrompt(templateName models.PromptTemplateEnum, data map[string]interface{}, user *models.AuthUser) (string, error) {
	switch templateName {
	case models.PromptTalentBio:
		return buildTalentBioPrompt(data, user)
	case models.PromptClientAboutUs:
		return buildClientAboutUsPrompt(data, user)
	case models.PromptJobDescription:
		return buildJobDescriptionPrompt(data, user)
	case models.PromptProposalCoverLetter:
		return buildProposalCoverLetterPrompt(data, user)
	case models.PromptNone:
		return "", nil
	default:
		return "", fmt.Errorf("unknown template: %s", templateName)
	}
}

// buildTalentBioPrompt creates a prompt for talent bio
func buildTalentBioPrompt(data map[string]interface{}, user *models.AuthUser) (string, error) {
	profile, ok := data["profile"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("profile field is required")
	}

	name := getString(profile, "name", "the talent")
	title := getString(profile, "title", "Professional")
	skills := getString(profile, "skills", "")
	rating := getFloat(profile, "rating", 0.0)
	jobsCompleted := getInt(profile, "jobs_completed", 0)
	pricing := getFloat(profile, "pricing", 0.0)
	badge := getString(profile, "badge", "")

	skillsList := ""
	if skills != "" {
		parts := strings.Split(skills, ",")
		if len(parts) > 5 {
			parts = parts[:5]
		}
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		skillsList = strings.Join(parts, ", ")
	}

	prompt := fmt.Sprintf(`You are an elite Dokoola profile copywriter.

Write a powerful, first-person bio for %s in MAX 500 characters (including spaces).

Current title: %s
Badge: %s
Rating: %.1f/5
Jobs completed: %d
Rate: $%.0f/hr
Top skills: %s

Rules:
- First person
- Zero filler words
- Instantly shows expertise + results
- Ends with a hook
- Under 500 chars total

Just output the bio simple rich markdown. Nothing else.`,
		name, title, badge, rating, jobsCompleted, pricing, skillsList)

	return prompt, nil
}

// buildClientAboutUsPrompt creates a prompt for client about us
func buildClientAboutUsPrompt(data map[string]interface{}, user *models.AuthUser) (string, error) {
	profile, ok := data["profile"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("profile field is required")
	}

	name := getString(profile, "name", "the company")
	industry := getString(profile, "industry", "technology")
	country := ""

	if company, ok := profile["company"].(map[string]interface{}); ok {
		if companyName := getString(company, "name", ""); companyName != "" {
			name = companyName
		}
		if companyIndustry := getString(company, "industry", ""); companyIndustry != "" {
			industry = companyIndustry
		}
		if companyCountry, ok := company["country"].(map[string]interface{}); ok {
			country = getString(companyCountry, "name", "")
		}
	}

	if country == "" {
		if profileCountry, ok := profile["country"].(map[string]interface{}); ok {
			country = getString(profileCountry, "name", "")
		}
	}

	companyText := fmt.Sprintf("Company: %s", name)
	industryText := ""
	if industry != "" {
		industryText = fmt.Sprintf("Industry: %s", industry)
	}
	countryText := ""
	if country != "" {
		countryText = fmt.Sprintf("Location: based in %s", country)
	}

	prompt := fmt.Sprintf(`Write a sharp, trust-building "About Us" blurb (third person) for this client.

%s
%s
%s

Requirements:
- Max 300 characters (including spaces)
- Professional & credible
- Shows they're serious about hiring top talent
- No generic fluff

Just output the final text in simple rich marckdown. Nothing else.`,
		companyText, industryText, countryText)

	return prompt, nil
}

// buildJobDescriptionPrompt creates a prompt for job description
func buildJobDescriptionPrompt(data map[string]interface{}, user *models.AuthUser) (string, error) {
	job, ok := data["job"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("job field is required")
	}

	client, ok := data["client"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("client field is required")
	}

	tone := "professional"
	length := "medium"

	title := getString(job, "title", "Untitled Role")
	categoryName := getString(job, "category_name", "Not specified")
	jobType := getString(job, "job_type", "")
	duration := getString(job, "estimated_duration", "")
	requiredSkills := getStringSlice(job, "required_skills")
	// description := getString(job, "description", "") // Not used in prompt
	address := getString(job, "address", "")
	isThirdParty := getBool(job, "is_third_party", false)

	clientName := getString(client, "name", "")
	clientAbout := getString(client, "about", "")

	companyName := clientName
	if isThirdParty {
		if thirdParty, ok := job["third_party_metadata"].(map[string]interface{}); ok {
			if tpName := getString(thirdParty, "company_name", ""); tpName != "" {
				companyName = tpName
			}
		}
	}

	// Budget info
	budgetInfo := ""
	if pricing, ok := job["pricing"].(map[string]interface{}); ok {
		budget := getFloat(pricing, "budget", 0)
		fixedPrice := getBool(pricing, "fixed_price", false)
		currencySymbol := "$"
		if currency, ok := pricing["currency"].(map[string]interface{}); ok {
			currencySymbol = getString(currency, "symbol", "$")
		}

		if budget > 0 {
			if fixedPrice {
				budgetInfo = fmt.Sprintf("%s%.0f fixed-price", currencySymbol, budget)
			} else {
				budgetInfo = fmt.Sprintf("%s%.0f/hr", currencySymbol, budget)
			}
		} else if fixedPrice {
			budgetInfo = "Fixed-price project"
		} else {
			budgetInfo = "Hourly rate (budget negotiable)"
		}
	}

	deadline := "Open until filled"
	if deadlineStr := getString(job, "application_deadline", ""); deadlineStr != "" {
		if t, err := time.Parse(time.RFC3339, strings.Replace(deadlineStr, "Z", "+00:00", 1)); err == nil {
			deadline = fmt.Sprintf("Apply by %s", t.Format("January 02, 2006"))
		}
	}

	skills := strings.Join(requiredSkills, " · ")
	if skills == "" {
		skills = "Relevant skills required"
	}

	location := ""
	if countryData, ok := job["country"].(map[string]interface{}); ok {
		countryName := getString(countryData, "name", "")
		if countryName != "" {
			if address != "" {
				location = fmt.Sprintf("%s - %s", countryName, address)
			} else {
				location = countryName
			}
		}
	}
	if location == "" && address != "" {
		location = address
	}

	thirdPartyText := ""
	if isThirdParty {
		thirdPartyText = " (posted via Dokoola)"
	}

	aboutUsText := ""
	if clientAbout != "" {
		aboutUsText = fmt.Sprintf("About Us: %s", clientAbout)
	}

	budgetText := ""
	if budgetInfo != "" {
		budgetText = fmt.Sprintf("Budget: %s", budgetInfo)
	}

	jobTypeText := ""
	if jobType != "" {
		jobTypeText = fmt.Sprintf("Job Type: %s", jobType)
	}

	durationText := ""
	if duration != "" {
		durationText = fmt.Sprintf("Estimated Duration: %s", duration)
	}

	locationText := ""
	if location != "" {
		locationText = fmt.Sprintf("Location Preference: %s", location)
	}

	prompt := fmt.Sprintf(`Rewrite this job posting into a clear, professional, and attractive job description that top talent actually wants to apply to.

=== BASIC INFO ===
Title: %s
Category: %s
Company: %s%s
%s

=== DETAILS FROM CLIENT ===

%s
%s
%s
Application Deadline: %s
Required Skills: %s
%s

=== INSTRUCTIONS ===
Write a polished, engaging job description that:
- Starts with a strong, specific openisng line (not "We are looking for...")
- Clearly explains the project/role and its impact
- Highlights what success looks like
- Mentions budget, timeline, and type upfront (no hiding)
- Lists required skills cleanly
- Ends with a confident, welcoming call-to-action
- Format your response strictly in rich-text using (basic-html-tags).

Tone: %s, direct, respectful of freelancers' time
Length: %s 250-450 words max
Style: Human, concise, zero fluff
    
No titles, no quotes, no extra text.
Just output the final job description in simple rich markdown.`,
		title, categoryName, companyName, thirdPartyText, aboutUsText,
		budgetText, jobTypeText, durationText, deadline, skills, locationText,
		tone, length)

	return prompt, nil
}

// buildProposalCoverLetterPrompt creates a prompt for proposal cover letter
func buildProposalCoverLetterPrompt(data map[string]interface{}, user *models.AuthUser) (string, error) {
	talent, ok := data["talent"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("talent field is required")
	}

	job, ok := data["job"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("job field is required")
	}

	resume := data["resume"]
	metadata, _ := data["metadata"].(map[string]interface{})

	tone := "professional"
	length := "medium"
	additionalNotes := ""

	if metadata != nil {
		tone = getString(metadata, "tone", "professional")
		length = getString(metadata, "length", "medium")
		additionalNotes = getString(metadata, "additional_notes", "")
	}

	// Talent info
	talentName := getString(talent, "name", "")
	talentTitle := getString(talent, "title", "")
	talentBio := getString(talent, "bio", "")
	talentSkills := getString(talent, "skills", "")
	talentPricing := getFloat(talent, "pricing", 0)
	talentRating := getFloat(talent, "rating", 0.5)
	talentBadge := getString(talent, "badge", "")

	// Job info
	jobTitle := getString(job, "title", "")
	jobCategory := ""
	if category, ok := job["category"].(map[string]interface{}); ok {
		jobCategory = getString(category, "name", "")
	}
	jobType := getString(job, "job_type", "")
	jobDescription := getString(job, "description", "No job description provided.")
	experienceLevel := getString(job, "experience_level", "")
	requiredSkills := getStringSlice(job, "required_skills")
	duration := getString(job, "estimated_duration", "")

	// Client info
	var clientName string
	if clientData, ok := job["client"].(map[string]interface{}); ok {
		clientName = getString(clientData, "name", "")
	}

	companyName := clientName
	if thirdParty, ok := job["third_party_metadata"].(map[string]interface{}); ok {
		if tpName := getString(thirdParty, "company_name", ""); tpName != "" {
			companyName = tpName
		}
	}

	// Badge text
	badgeText := "verified talent"
	if talentBadge == "pro" {
		badgeText = "Pro"
	} else if talentBadge == "star" {
		badgeText = "Top-Rated Star"
	}

	// Rating text
	ratingStars := "experienced talent"
	if talentRating > 0 {
		ratingStars = fmt.Sprintf("rated %.1f/5", talentRating)
	}

	// Resume text
	resumeText := ""
	if resume != nil {
		if resumeMap, ok := resume.(map[string]interface{}); ok {
			if desc := getString(resumeMap, "description", ""); desc != "" {
				resumeText = fmt.Sprintf("Relevant Experience (from resume): %s", desc)
			}
		}
	}

	skillsText := ""
	if talentSkills != "" {
		skillsText = fmt.Sprintf("Key Skills: %s", talentSkills)
	}

	bioText := ""
	if talentBio != "" {
		bioText = fmt.Sprintf("Bio Summary:: %s", talentBio)
	}

	pricingText := ""
	if talentPricing > 0 {
		pricingText = fmt.Sprintf("Pricing Rate: $%.0f", talentPricing)
	}

	additionalText := ""
	if additionalNotes != "" {
		additionalText = fmt.Sprintf("Additional instructions from me: %s", additionalNotes)
	}

	requiredSkillsStr := strings.Join(requiredSkills, ", ")
	if requiredSkillsStr == "" {
		requiredSkillsStr = "relevant skills"
	}

	jobTypeText := ""
	if jobType != "" {
		jobTypeText = fmt.Sprintf("Job Type: %s", jobType)
	}

	durationText := ""
	if duration != "" {
		durationText = fmt.Sprintf("Project Duration: %s", duration)
	}

	clientNameText := fmt.Sprintf("Client Name: %s", clientName)
	if companyName != "" {
		clientNameText = fmt.Sprintf("Company Name: %s", companyName)
	}

	prompt := fmt.Sprintf(`You are an expert freelance proposal writer for (Dokoola) who has helped hundreds of freelancers win high-value contracts on platforms like Upwork, Fiverr, and Toptal.

Write a compelling, personalized cover letter (proposal) for the following freelancer applying to this job.

=== FREELANCER PROFILE ===
Name: %s
Title: %s
Badge: %s
Rating: %s
%s
%s
%s
%s

=== JOB DETAILS ===
Client/Company: %s
%s
Job Title: %s
%s
Job Category: %s
Experience Level Expected: %s
Required Skills: %s
Project Duration: %s
%s
Full Job Description: """%s"""

=== INSTRUCTIONS ===
Write a winning proposal cover letter in first person as %s.

Structure:
1. Strong opening: Greet the client and express genuine interest in THEIR specific project (reference something unique from the job description).
2. Prove fit: Explain why I am the perfect match (highlight overlapping skills, past results, and relevant experience from my profile/resume).
3. Build trust: Mention my rating, badge, and success record.
4. Show understanding: Demonstrate that I fully understand the project goals and challenges.
5. Call to action: End with confidence and invite next steps (interview, questions, etc.).

Tone: %s
Length: %s
Style: Natural, human, engaging — never robotic or generic. Avoid this chars "—" or clichés like "I am passionate about" unless it feels authentic.

%s

Do NOT mention that this was AI-generated.
Do NOT say "As an AI language model".
Just write the cover letter — nothing else.
Format your response strictly in rich-text using (simple markdown).`,
		talentName, talentTitle, badgeText, ratingStars,
		skillsText, bioText, pricingText, resumeText,
		companyName, clientNameText, jobTitle, jobTypeText, jobCategory, experienceLevel,
		requiredSkillsStr, duration, durationText, jobDescription,
		talentName,
		tone, length,
		additionalText)

	return prompt, nil
}

// Helper functions
func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getFloat(m map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return defaultValue
}

func getInt(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return defaultValue
}

func getBool(m map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}

func getStringSlice(m map[string]interface{}, key string) []string {
	if val, ok := m[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
		if slice, ok := val.([]string); ok {
			return slice
		}
	}
	return []string{}
}
