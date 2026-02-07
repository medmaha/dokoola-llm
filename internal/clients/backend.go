package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/dokoola/llm-go/internal/models"
	"go.uber.org/zap"
)

// BackendClient handles requests to the backend API
type BackendClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
	categories []models.JobCategory
	mu         sync.RWMutex
}

// NewBackendClient creates a new backend API client
func NewBackendClient(baseURL string, logger *zap.Logger) *BackendClient {
	return &BackendClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

// GetUser fetches user data from the backend
func (c *BackendClient) GetUser(userID string) (*models.AuthUser, error) {
	url := fmt.Sprintf("%s/api/users/%s/llm/", c.baseURL, userID)

	c.logger.Debug("Fetching user from backend", zap.String("user_id", userID))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backend API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var user models.AuthUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	c.logger.Debug("User fetched successfully",
		zap.String("user_id", userID),
		zap.String("name", user.Name),
	)

	return &user, nil
}

// GetCategories fetches job categories from the backend
func (c *BackendClient) GetCategories() ([]models.JobCategory, error) {
	// Check cache first
	c.mu.RLock()
	if len(c.categories) > 0 {
		cached := c.categories
		c.mu.RUnlock()
		c.logger.Debug("Returning cached categories", zap.Int("count", len(cached)))
		return cached, nil
	}
	c.mu.RUnlock()

	// Fetch from backend
	url := fmt.Sprintf("%s/api/categories?scraper=true", c.baseURL)

	c.logger.Debug("Fetching categories from backend")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backend API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var categories []models.JobCategory
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal categories: %w", err)
	}

	// Cache the categories
	c.mu.Lock()
	c.categories = categories
	c.mu.Unlock()

	c.logger.Info("Categories fetched and cached", zap.Int("count", len(categories)))

	return categories, nil
}
