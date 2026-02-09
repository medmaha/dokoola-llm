package clients

import (
	"testing"

	"go.uber.org/zap"
)

func initTestLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func TestNewBackendClient(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	baseURL := "https://api.example.com"
	client := NewBackendClient(baseURL, logger)

	if client == nil {
		t.Error("expected non-nil client")
	}

	if client.baseURL != baseURL {
		t.Errorf("expected baseURL %q, got %q", baseURL, client.baseURL)
	}

	if client.logger == nil {
		t.Error("expected logger to be set")
	}

	if client.httpClient == nil {
		t.Error("expected httpClient to be set")
	}
}

func TestBackendClientFieldInitialization(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	baseURL := "https://test.com"
	client := NewBackendClient(baseURL, logger)

	// Verify categories cache is initialized empty
	if len(client.categories) != 0 {
		t.Errorf("expected empty categories cache, got %d items", len(client.categories))
	}

	// Verify mutex is initialized (we can't directly test mutex, but we can verify the client structure)
	if client == nil {
		t.Error("client should not be nil")
	}
}

func TestMultipleBackendClientInstances(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	client1 := NewBackendClient("https://api1.com", logger)
	client2 := NewBackendClient("https://api2.com", logger)

	if client1 == nil || client2 == nil {
		t.Error("expected both clients to be non-nil")
	}

	if client1.baseURL == client2.baseURL {
		t.Error("expected different baseURLs for different clients")
	}
}

func TestBackendClientDifferentBaseURLs(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{name: "local", baseURL: "http://localhost:8000"},
		{name: "production", baseURL: "https://api.example.com"},
		{name: "staging", baseURL: "https://staging-api.example.com"},
	}

	logger, _ := initTestLogger()
	defer logger.Sync()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewBackendClient(tt.baseURL, logger)

			if client.baseURL != tt.baseURL {
				t.Errorf("expected %q, got %q", tt.baseURL, client.baseURL)
			}
		})
	}
}

func TestBackendClientLoggerNotNil(t *testing.T) {
	logger, _ := initTestLogger()
	defer logger.Sync()

	client := NewBackendClient("https://api.com", logger)

	if client.logger == nil {
		t.Error("expected logger to be set in client")
	}

	if client.logger != logger {
		t.Error("expected client logger to be the passed logger instance")
	}
}

func TestBackendClientCachingBehavior(t *testing.T) {
	// This test verifies the caching structure is properly initialized
	logger, _ := initTestLogger()
	defer logger.Sync()

	client := NewBackendClient("https://api.com", logger)

	// Check that categories field exists and is empty
	if client.categories == nil {
		// categories might be nil or empty slice, both are acceptable initial states
		t.Log("categories field is nil (acceptable initial state)")
	} else if len(client.categories) != 0 {
		t.Errorf("expected empty categories cache initially, got %d items", len(client.categories))
	}
}
