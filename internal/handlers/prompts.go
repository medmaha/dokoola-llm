package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dokoola/llm-go/internal/clients"
	"github.com/dokoola/llm-go/internal/llm"
	"github.com/dokoola/llm-go/internal/models"
	"github.com/dokoola/llm-go/internal/prompts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PromptsHandler handles prompt generation requests
type PromptsHandler struct {
	llmClient     *llm.Client
	backendClient *clients.BackendClient
	logger        *zap.Logger
}

// NewPromptsHandler creates a new prompts handler
func NewPromptsHandler(llmClient *llm.Client, backendClient *clients.BackendClient, logger *zap.Logger) *PromptsHandler {
	return &PromptsHandler{
		llmClient:     llmClient,
		backendClient: backendClient,
		logger:        logger,
	}
}

// GeneratePrompt handles POST /api/v1/llm/actions/generate-prompt
func (h *PromptsHandler) GeneratePrompt(c *gin.Context) {
	var req models.PromptGenerationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.PromptGenerationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		})
		return
	}

	h.logger.Info("Received prompt generation request",
		zap.String("template", string(req.TemplateName)),
	)

	// Get user ID from query parameter (optional)
	userID := c.Query("user_id")
	var user *models.AuthUser
	var err error

	user, err = h.backendClient.GetUser(userID)
	if err != nil {
		h.logger.Warn("Failed to fetch user, continuing without user context",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		errorMessage := fmt.Sprintf("User not found: %s", userID)
		c.JSON(http.StatusNotFound, models.TextCompletionResponse{
			Success:      false,
			ErrorMessage: &errorMessage,
		})
		return
	}

	// Build prompt from template
	prompt, err := prompts.BuildPrompt(req.TemplateName, req.Data, user)
	if err != nil {
		h.logger.Error("Failed to build prompt", zap.Error(err))
		errorMsg := fmt.Sprintf("Failed to build prompt: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.PromptGenerationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		})
		return
	}

	// Handle "none" template - just return empty completion
	if req.TemplateName == models.PromptNone {
		emptyStr := ""
		c.JSON(http.StatusOK, models.PromptGenerationResponse{
			Success:    true,
			Completion: &emptyStr,
		})
		return
	}

	h.logger.Debug("Prompt built successfully", zap.Int("prompt_length", len(prompt)))

	// Get LLM completion
	completion, err := h.llmClient.Complete(prompt, user)
	if err != nil {
		// If upstream LLM is rate-limited, return 503 to caller
		if errors.Is(err, llm.ErrRateLimited) {
			h.logger.Warn("Upstream LLM rate limited", zap.Error(err))
			msg := "Upstream LLM service overloaded; please try again later"
			c.JSON(http.StatusServiceUnavailable, models.PromptGenerationResponse{
				Success:      false,
				ErrorMessage: &msg,
			})
			return
		}

		h.logger.Error("LLM completion failed", zap.Error(err))
		errorMsg := fmt.Sprintf("Failed to generate completion: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.PromptGenerationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		})
		return
	}

	h.logger.Info("Prompt generation successful",
		zap.String("template", string(req.TemplateName)),
	)

	c.JSON(http.StatusOK, models.PromptGenerationResponse{
		Success:    true,
		Completion: &completion,
	})
}
