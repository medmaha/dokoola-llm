package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dokoola/llm-go/internal/clients"
	"github.com/dokoola/llm-go/internal/llm"
	"github.com/dokoola/llm-go/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JobsHandler handles job categorization requests
type JobsHandler struct {
	llmClient     *llm.Client
	backendClient *clients.BackendClient
	logger        *zap.Logger
}

// NewJobsHandler creates a new jobs handler
func NewJobsHandler(llmClient *llm.Client, backendClient *clients.BackendClient, logger *zap.Logger) *JobsHandler {
	return &JobsHandler{
		llmClient:     llmClient,
		backendClient: backendClient,
		logger:        logger,
	}
}

// CategorizeJobs handles POST /api/v1/llm/jobs/categorize
func (h *JobsHandler) CategorizeJobs(c *gin.Context) {
	var req models.JobCategorizationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.JobCategorizationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		})
		return
	}

	h.logger.Info("Received job categorization request", zap.Int("job_count", len(req.Data)))

	// Fetch categories from backend
	categories, err := h.backendClient.GetCategories()
	if err != nil {
		h.logger.Error("Failed to fetch categories", zap.Error(err))
		errorMsg := fmt.Sprintf("Failed to fetch categories: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.JobCategorizationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		})
		return
	}

	// Build categories description for prompt
	categoriesDesc := h.buildCategoriesDescription(categories)

	// Process each job
	results := make([]models.JobResponseData, 0, len(req.Data))

	for _, job := range req.Data {
		h.logger.Debug("Processing job", zap.String("public_id", job.PublicID))

		// Build categorization prompt
		prompt := fmt.Sprintf(`You are a job categorization expert for Dokoola platform.

Analyze this job posting and select the SINGLE MOST RELEVANT category from the list below.

Job Description:
"""%s"""

Available Categories:
%s

Instructions:
- Return ONLY the category slug (e.g., "web-development")
- Choose the most specific matching category
- If no exact match, choose the closest parent category
- Return only the slug, nothing else

Category slug:`, job.Description, categoriesDesc)

		// Get LLM completion
		completion, err := h.llmClient.Complete(prompt, nil)
		if err != nil {
			h.logger.Error("LLM completion failed", zap.String("public_id", job.PublicID), zap.Error(err))
			// Use default category or empty string
			results = append(results, models.JobResponseData{
				PublicID: job.PublicID,
				Category: "",
			})
			continue
		}

		// Extract category slug from response
		categorySlug := strings.TrimSpace(completion)
		categorySlug = strings.Trim(categorySlug, "\"'`")

		h.logger.Debug("Job categorized",
			zap.String("public_id", job.PublicID),
			zap.String("category", categorySlug),
		)

		results = append(results, models.JobResponseData{
			PublicID: job.PublicID,
			Category: categorySlug,
		})
	}

	h.logger.Info("Job categorization completed", zap.Int("processed", len(results)))

	c.JSON(http.StatusOK, models.JobCategorizationResponse{
		Success: true,
		Data:    results,
	})
}

// buildCategoriesDescription creates a formatted string of categories
func (h *JobsHandler) buildCategoriesDescription(categories []models.JobCategory) string {
	var sb strings.Builder

	for i, cat := range categories {
		if cat.ParentSlug != nil {
			sb.WriteString(fmt.Sprintf("%d. %s (%s) - Child of: %s\n",
				i+1, cat.Slug, cat.Description, *cat.ParentSlug))
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s (%s)\n",
				i+1, cat.Slug, cat.Description))
		}
	}

	return sb.String()
}
