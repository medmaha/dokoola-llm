package config

import (
	"os"
	"testing"
)

func TestGetEnvString(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultVal   string
		envValue     string
		expected     string
		shouldSetEnv bool
	}{
		{
			name:         "returns env value when set",
			key:          "TEST_ENV_VAR",
			defaultVal:   "default",
			envValue:     "from_env",
			expected:     "from_env",
			shouldSetEnv: true,
		},
		{
			name:         "returns default when not set",
			key:          "NONEXISTENT_VAR",
			defaultVal:   "default",
			envValue:     "",
			expected:     "default",
			shouldSetEnv: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up first
			os.Unsetenv(tt.key)

			if tt.shouldSetEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultVal   int
		envValue     string
		expected     int
		shouldSetEnv bool
	}{
		{
			name:         "returns parsed int when valid",
			key:          "TEST_PORT",
			defaultVal:   8000,
			envValue:     "9000",
			expected:     9000,
			shouldSetEnv: true,
		},
		{
			name:         "returns default when not set",
			key:          "NONEXISTENT_PORT",
			defaultVal:   8000,
			envValue:     "",
			expected:     8000,
			shouldSetEnv: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)

			if tt.shouldSetEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvInt(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultVal   bool
		envValue     string
		expected     bool
		shouldSetEnv bool
	}{
		{
			name:         "returns true for true value",
			key:          "TEST_DEBUG",
			defaultVal:   false,
			envValue:     "true",
			expected:     true,
			shouldSetEnv: true,
		},
		{
			name:         "returns false for false value",
			key:          "TEST_DEBUG2",
			defaultVal:   true,
			envValue:     "false",
			expected:     false,
			shouldSetEnv: true,
		},
		{
			name:         "returns default when not set",
			key:          "NONEXISTENT_DEBUG",
			defaultVal:   true,
			envValue:     "",
			expected:     true,
			shouldSetEnv: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)

			if tt.shouldSetEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvBool(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSettingsStructure(t *testing.T) {
	settings := Settings{
		AppName:    "Test App",
		AppVersion: "1.0.0",
		Debug:      true,
		APIPrefix:  "/api/v1",
		Host:       "0.0.0.0",
		Port:       8000,
		ENV:        "development",
	}

	if settings.AppName != "Test App" {
		t.Errorf("expected AppName 'Test App', got %q", settings.AppName)
	}

	if settings.Port != 8000 {
		t.Errorf("expected Port 8000, got %d", settings.Port)
	}

	if !settings.Debug {
		t.Error("expected Debug to be true")
	}
}

func TestServiceConfigStructure(t *testing.T) {
	service := ServiceConfig{
		Host:       "api.example.com",
		ClientName: "test-client",
		SecretHash: "hash123",
	}

	if service.Host != "api.example.com" {
		t.Errorf("expected Host 'api.example.com', got %q", service.Host)
	}

	if service.ClientName != "test-client" {
		t.Errorf("expected ClientName 'test-client', got %q", service.ClientName)
	}

	if service.SecretHash != "hash123" {
		t.Errorf("expected SecretHash 'hash123', got %q", service.SecretHash)
	}
}

func TestConfigStructure(t *testing.T) {
	settings := Settings{
		AppName:    "App",
		AppVersion: "1.0",
		Debug:      false,
		APIPrefix:  "/api",
		Host:       "localhost",
		Port:       8000,
		ENV:        "prod",
	}

	services := map[string]ServiceConfig{
		"key1": {
			Host:       "host1",
			ClientName: "client1",
			SecretHash: "hash1",
		},
	}

	origins := []string{"https://example.com"}

	cfg := Config{
		Settings:         settings,
		AllowedServices:  services,
		AllowedOrigins:   origins,
		CerebrasAPIKey:   "key123",
		BackendServerAPI: "https://backend.example.com",
		ServiceKeyName:   "X-Service-Key",
		ClientNameHeader: "X-Client-Name",
		SecretHashHeader: "X-Secret-Hash",
	}

	if cfg.Settings.AppName != "App" {
		t.Errorf("expected Settings.AppName 'App', got %q", cfg.Settings.AppName)
	}

	if len(cfg.AllowedServices) != 1 {
		t.Errorf("expected 1 service, got %d", len(cfg.AllowedServices))
	}

	if len(cfg.AllowedOrigins) != 1 {
		t.Errorf("expected 1 origin, got %d", len(cfg.AllowedOrigins))
	}

	if cfg.CerebrasAPIKey != "key123" {
		t.Errorf("expected CerebrasAPIKey 'key123', got %q", cfg.CerebrasAPIKey)
	}
}
