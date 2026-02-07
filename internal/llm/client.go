package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dokoola/llm-go/internal/constants"
	"github.com/dokoola/llm-go/internal/models"
	"go.uber.org/zap"
)

const (
	cerebrasAPIURL = "https://api.cerebras.ai/v1/chat/completions"
	modelName      = "gpt-oss-120b"
	maxTokens      = 40960
	temperature    = 0.6
	topP           = 0.95
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is the request to Cerebras API
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	Stream      bool      `json:"stream"`
}

// ChatCompletionResponse is the response from Cerebras API
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Client handles LLM API requests
type Client struct {
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewClient creates a new LLM client
func NewClient(apiKey string, logger *zap.Logger) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

// Complete sends a completion request to the LLM API
func (c *Client) Complete(userPrompt string, user *models.AuthUser) (string, error) {
	messages := c.buildMessages(userPrompt, user)

	reqBody := ChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	c.logger.Debug("Sending completion request",
		zap.String("model", modelName),
		zap.Int("message_count", len(messages)),
	)

	req, err := http.NewRequest("POST", cerebrasAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
			c.logger.Error("LLM API error",
				zap.Int("status_code", resp.StatusCode),
				zap.String("error", errResp.Error.Message),
			)
			return "", fmt.Errorf("LLM API error: %s", errResp.Error.Message)
		}
		return "", fmt.Errorf("LLM API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var completionResp ChatCompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	completion := completionResp.Choices[0].Message.Content

	c.logger.Info("LLM completion successful",
		zap.Int("prompt_tokens", completionResp.Usage.PromptTokens),
		zap.Int("completion_tokens", completionResp.Usage.CompletionTokens),
		zap.Int("total_tokens", completionResp.Usage.TotalTokens),
	)

	return completion, nil
}

// buildMessages constructs the message array for the LLM request
func (c *Client) buildMessages(userPrompt string, user *models.AuthUser) []Message {
	messages := make([]Message, 0, len(constants.SystemMessages)+2)

	// Add system messages
	for _, sysMsg := range constants.SystemMessages {
		messages = append(messages, Message{
			Role:    sysMsg["role"],
			Content: sysMsg["content"],
		})
	}

	// Add user context if available
	if user != nil {
		userContext := fmt.Sprintf(`Current User Context:
- Name: %s
- User ID: %s
- Profile Status: %s
- Account Type: %s`,
			user.Name,
			user.PublicID,
			func() string {
				if user.CompleteProfile {
					return "Complete"
				}
				return "Incomplete"
			}(),
			func() string {
				if user.IsTalent {
					return "Talent/Freelancer"
				}
				if user.IsClient {
					return "Client/Employer"
				}
				if user.IsStaff {
					return "Staff"
				}
				return "User"
			}(),
		)

		messages = append(messages, Message{
			Role:    "system",
			Content: userContext,
		})
	}

	// Add user prompt
	messages = append(messages, Message{
		Role:    "user",
		Content: userPrompt,
	})

	return messages
}
