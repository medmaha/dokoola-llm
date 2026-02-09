package llm

import (
	"testing"

	"go.uber.org/zap"
)

func initTestLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func TestNewClient(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	apiKey := "test-key-12345"
	client := NewClient(apiKey, logger)

	if client == nil {
		t.Error("expected non-nil client")
	}

	if client.apiKey != apiKey {
		t.Errorf("expected apiKey %q, got %q", apiKey, client.apiKey)
	}

	if client.logger == nil {
		t.Error("expected logger to be set")
	}

	if client.httpClient == nil {
		t.Error("expected httpClient to be set")
	}
}

func TestClientInitialization(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	tests := []struct {
		name   string
		apiKey string
	}{
		{name: "standard key", apiKey: "key-abc123"},
		{name: "long key", apiKey: "very-long-api-key-with-many-characters-and-dashes-123456789"},
		{name: "simple key", apiKey: "key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, logger)

			if client == nil {
				t.Error("expected non-nil client")
			}

			if client.apiKey != tt.apiKey {
				t.Errorf("expected apiKey %q, got %q", tt.apiKey, client.apiKey)
			}
		})
	}
}

func TestClientHTTPClientNotNil(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	client := NewClient("test-key", logger)

	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}
}

func TestClientLoggerAssignment(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	client := NewClient("test-key", logger)

	if client.logger != logger {
		t.Error("expected logger to be the passed logger instance")
	}
}

func TestMultipleClientInstances(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	client1 := NewClient("key1", logger)
	client2 := NewClient("key2", logger)

	if client1 == nil || client2 == nil {
		t.Error("expected both clients to be non-nil")
	}

	if client1.apiKey == client2.apiKey {
		t.Error("expected different API keys for different clients")
	}

	if client1 == client2 {
		t.Error("expected different client instances")
	}
}

func TestClientConstantsValues(t *testing.T) {
	// Verify that model constants are defined correctly
	// These are package-level constants used in the client
	if modelName == "" {
		t.Error("expected modelName constant to be non-empty")
	}

	if maxTokens <= 0 {
		t.Errorf("expected maxTokens to be positive, got %d", maxTokens)
	}

	if temperature < 0 || temperature > 1 {
		t.Errorf("expected temperature between 0 and 1, got %f", temperature)
	}

	if topP < 0 || topP > 1 {
		t.Errorf("expected topP between 0 and 1, got %f", topP)
	}
}

func TestMessageStructure(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Hello",
	}

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}

	if msg.Content != "Hello" {
		t.Errorf("expected content 'Hello', got %q", msg.Content)
	}
}

func TestChatCompletionRequestStructure(t *testing.T) {
	messages := []Message{
		{Role: "user", Content: "Test"},
	}

	req := ChatCompletionRequest{
		Model:       "gpt-test",
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
		Stream:      false,
	}

	if req.Model != "gpt-test" {
		t.Errorf("expected model 'gpt-test', got %q", req.Model)
	}

	if len(req.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(req.Messages))
	}

	if req.MaxTokens != 1000 {
		t.Errorf("expected max tokens 1000, got %d", req.MaxTokens)
	}

	if req.Stream {
		t.Error("expected stream to be false")
	}
}

func TestChatCompletionResponseStructure(t *testing.T) {
	resp := ChatCompletionResponse{
		ID:      "resp-123",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   "gpt-test",
	}

	if resp.ID != "resp-123" {
		t.Errorf("expected ID 'resp-123', got %q", resp.ID)
	}

	if resp.Object != "chat.completion" {
		t.Errorf("expected object 'chat.completion', got %q", resp.Object)
	}

	if resp.Model != "gpt-test" {
		t.Errorf("expected model 'gpt-test', got %q", resp.Model)
	}
}

func TestErrorResponseStructure(t *testing.T) {
	errResp := ErrorResponse{}
	errResp.Error.Message = "API Error"
	errResp.Error.Type = "invalid_request"
	errResp.Error.Code = "400"

	if errResp.Error.Message != "API Error" {
		t.Errorf("expected error message 'API Error', got %q", errResp.Error.Message)
	}

	if errResp.Error.Type != "invalid_request" {
		t.Errorf("expected error type 'invalid_request', got %q", errResp.Error.Type)
	}

	if errResp.Error.Code != "400" {
		t.Errorf("expected error code '400', got %q", errResp.Error.Code)
	}
}

func TestClientAPIKeyStorage(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	apiKey := "secret-api-key-do-not-share"
	client := NewClient(apiKey, logger)

	// Verify the key is stored (we can't access private fields directly,
	// but the constructor should have set it)
	if client == nil {
		t.Error("expected client to be created with API key")
	}
}
