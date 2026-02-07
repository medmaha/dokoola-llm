package handlers

import (
	"fmt"
	"net/http"

	"github.com/dokoola/llm-go/internal/clients"
	"github.com/dokoola/llm-go/internal/llm"
	"github.com/dokoola/llm-go/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TextCompletionHandler handles text completion requests
type TextCompletionHandler struct {
	llmClient     *llm.Client
	backendClient *clients.BackendClient
	logger        *zap.Logger
}

// NewTextCompletionHandler creates a new text completion handler
func NewTextCompletionHandler(llmClient *llm.Client, backendClient *clients.BackendClient, logger *zap.Logger) *TextCompletionHandler {
	return &TextCompletionHandler{
		llmClient:     llmClient,
		backendClient: backendClient,
		logger:        logger,
	}
}

// Complete handles POST /api/v1/chat/completion
func (h *TextCompletionHandler) Complete(c *gin.Context) {
	var req models.TextCompletionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.TextCompletionResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		})
		return
	}

	h.logger.Info("Received text completion request")

	// Get user ID from query parameter (optional for text completion)
	userID := c.Query("user_id")
	var user *models.AuthUser
	var err error

	if userID != "" {
		user, err = h.backendClient.GetUser(userID)
		if err != nil {
			h.logger.Warn("Failed to fetch user, continuing without user context",
				zap.String("user_id", userID),
				zap.Error(err),
			)
			errorMessage := fmt.Sprintf("User not found: %s", userID)
			c.JSON(http.StatusNotFound, models.TextCompletionResponse{
				Success: false,
				ErrorMessage: &errorMessage,
			})
			return 
		}
	}

	// Get LLM completion
	completion, err := h.llmClient.Complete(req.Text, user)
	if err != nil {
		h.logger.Error("LLM completion failed", zap.Error(err))
		errorMsg := fmt.Sprintf("Failed to generate completion: %s", err.Error())
		c.JSON(http.StatusInternalServerError, models.TextCompletionResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		})
		return
	}

	h.logger.Info("Text completion successful")

	c.JSON(http.StatusOK, models.TextCompletionResponse{
		Success:    true,
		Completion: &completion,
	})
}
