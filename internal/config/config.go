package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

// Settings holds the application configuration
type Settings struct {
	AppName    string
	AppVersion string
	Debug      bool
	APIPrefix  string
	Host       string
	Port       int
	ENV		string
}

// ServiceConfig holds configuration for an allowed service
type ServiceConfig struct {
	Host       string
	ClientName string
	SecretHash string
}

// Config holds all configuration
type Config struct {
	Settings         Settings
	AllowedServices  map[string]ServiceConfig
	AllowedOrigins   []string
	CerebrasAPIKey   string
	BackendServerAPI string
	ServiceKeyName   string
	ClientNameHeader string
	SecretHashHeader string
}

// LoadConfig loads configuration from environment variables and config.ini
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Settings: Settings{
			AppName:    getEnv("APP_NAME", "Dokoola LLM Service"),
			AppVersion: getEnv("APP_VERSION", "0.1.0"),
			APIPrefix:  getEnv("API_PREFIX", "/api/v1"),
			Host:       getEnv("HOST", "0.0.0.0"),
			Port:       getEnvInt("PORT", 8000),
			Debug:      getEnvBool("DEBUG", false) != false || os.Getenv("DEBUG") != "false",
			ENV:        getEnv("ENV", "production"),
		},
		CerebrasAPIKey:   strings.TrimSpace(os.Getenv("LLM_API_KEY")),
		BackendServerAPI: strings.TrimSpace(os.Getenv("BACKEND_SERVER_API")),
		ServiceKeyName:   strings.TrimSpace(os.Getenv("X_SERVICE_KEY_NAME")),
		ClientNameHeader: strings.TrimSpace(os.Getenv("X_SERVICE_CLIENT_NAME")),
		SecretHashHeader: strings.TrimSpace(os.Getenv("X_SERVICE_SECRET_NAME")),
	}

	// Validate required environment variables
	if cfg.CerebrasAPIKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY environment variable is required")
	}
	if cfg.BackendServerAPI == "" {
		return nil, fmt.Errorf("BACKEND_SERVER_API environment variable is required")
	}
	if cfg.ServiceKeyName == "" {
		return nil, fmt.Errorf("X_SERVICE_KEY_NAME environment variable is required")
	}
	if cfg.ClientNameHeader == "" {
		return nil, fmt.Errorf("X_SERVICE_CLIENT_NAME environment variable is required")
	}
	if cfg.SecretHashHeader == "" {
		return nil, fmt.Errorf("X_SERVICE_SECRET_NAME environment variable is required")
	}

	// Load allowed services from config.ini
	services, err := loadServicesFromConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load services from config.ini: %w", err)
	}

	cfg.AllowedServices = services
	cfg.AllowedOrigins = extractOrigins(services)

	return cfg, nil
}

// loadServicesFromConfig loads allowed services from config.ini file
func loadServicesFromConfig() (map[string]ServiceConfig, error) {
	configPath := "config.ini"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config.ini not found at %s", configPath)
	}

	iniFile, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config.ini: %w", err)
	}

	services := make(map[string]ServiceConfig)

	for _, section := range iniFile.Sections() {
		name := section.Name()
		if name == "DEFAULT" || name == "" {
			continue
		}

		// Only process sections that start with SERVICE_DKL
		if len(name) > 8 && name[:8] == "SERVICE_" {
			serviceKey := name[8:] // Remove "SERVICE_" prefix to get "DKL..."

			services[serviceKey] = ServiceConfig{
				Host:       section.Key("host").String(),
				ClientName: section.Key("client_name").String(),
				SecretHash: section.Key("secret_hash").String(),
			}
		}
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no services found in config.ini")
	}

	return services, nil
}

// extractOrigins extracts allowed origins from service configurations
func extractOrigins(services map[string]ServiceConfig) []string {
	origins := make([]string, 0, len(services))
	for _, service := range services {
		origins = append(origins, service.Host)
	}
	return origins
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		// Trim whitespace and control characters (like \r\n)
		return strings.TrimSpace(value)
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		value = strings.TrimSpace(value)
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		value = strings.TrimSpace(value)
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
