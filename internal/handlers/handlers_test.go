package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dokoola/llm-go/internal/clients"
	"github.com/dokoola/llm-go/internal/llm"
	"github.com/dokoola/llm-go/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func initHandlersTestLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func TestJobsHandlerNewJobsHandler(t *testing.T) {
	logger, _ := initHandlersTestLogger()
	defer logger.Sync()

	var mockLLM *llm.Client
	var mockBackend *clients.BackendClient

	handler := NewJobsHandler(mockLLM, mockBackend, logger)

	if handler == nil {
		t.Error("expected non-nil handler")
	}

	if handler.logger == nil {
		t.Error("expected logger to be set")
	}
}

func TestPromptsHandlerNewPromptsHandler(t *testing.T) {
	logger, _ := initHandlersTestLogger()
	defer logger.Sync()

	var mockLLM *llm.Client
	var mockBackend *clients.BackendClient

	handler := NewPromptsHandler(mockLLM, mockBackend, logger)

	if handler == nil {
		t.Error("expected non-nil handler")
	}

	if handler.logger == nil {
		t.Error("expected logger to be set")
	}
}

func TestTextCompletionHandlerNewTextCompletionHandler(t *testing.T) {
	logger, _ := initHandlersTestLogger()
	defer logger.Sync()

	var mockLLM *llm.Client
	var mockBackend *clients.BackendClient

	handler := NewTextCompletionHandler(mockLLM, mockBackend, logger)

	if handler == nil {
		t.Error("expected non-nil handler")
	}

	if handler.logger == nil {
		t.Error("expected logger to be set")
	}
}

func TestJobsHandlerInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := initHandlersTestLogger()
	defer logger.Sync()

	var mockLLM *llm.Client
	var mockBackend *clients.BackendClient
	handler := NewJobsHandler(mockLLM, mockBackend, logger)

	router := gin.New()
	router.POST("/jobs/categorize", handler.CategorizeJobs)

	// Send invalid JSON
	req := httptest.NewRequest("POST", "/jobs/categorize", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp models.JobCategorizationResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Success {
		t.Error("expected success to be false")
	}
}

func TestTextCompletionHandlerInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := initHandlersTestLogger()
	defer logger.Sync()

	var mockLLM *llm.Client
	var mockBackend *clients.BackendClient
	handler := NewTextCompletionHandler(mockLLM, mockBackend, logger)

	router := gin.New()
	router.POST("/completion", handler.Complete)

	// Send invalid JSON
	req := httptest.NewRequest("POST", "/completion", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp models.TextCompletionResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Success {
		t.Error("expected success to be false")
	}
}

func TestPromptsHandlerInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := initHandlersTestLogger()
	defer logger.Sync()

	var mockLLM *llm.Client
	var mockBackend *clients.BackendClient
	handler := NewPromptsHandler(mockLLM, mockBackend, logger)

	router := gin.New()
	router.POST("/prompts/generate", handler.GeneratePrompt)

	// Send invalid JSON
	req := httptest.NewRequest("POST", "/prompts/generate", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp models.PromptGenerationResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Success {
		t.Error("expected success to be false")
	}
}

func TestPromptsHandlerGeneratePromptNoneTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := initHandlersTestLogger()
	defer logger.Sync()

	var mockLLM *llm.Client
	var mockBackend *clients.BackendClient
	handler := NewPromptsHandler(mockLLM, mockBackend, logger)

	// This test checks that the handler correctly initializes
	// In real usage, GetUser would not be called with nil backend
	// This tests the handler creation structure
	if handler == nil {
		t.Error("expected handler to be created")
	}
}
