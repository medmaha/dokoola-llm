package middleware

import (
	"net/http"
	"time"

	"github.com/dokoola/llm-go/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware validates service authentication using custom headers
func AuthMiddleware(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for health check endpoint
		if c.Request.URL.Path == "/health/" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// Get authentication headers
		serviceKey := c.GetHeader(cfg.ServiceKeyName)
		clientName := c.GetHeader(cfg.ClientNameHeader)
		secretHash := c.GetHeader(cfg.SecretHashHeader)

		// Validate headers are present
		if serviceKey == "" || clientName == "" || secretHash == "" {
			logger.Warn("Missing authentication headers",
				zap.String("path", c.Request.URL.Path),
				zap.String("service_key", serviceKey),
				zap.String("client_name", clientName),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Missing required authentication headers",
			})
			c.Abort()
			return
		}

		// Validate service exists in allowed services
		service, exists := cfg.AllowedServices[serviceKey]
		if !exists {
			logger.Warn("Invalid service key",
				zap.String("service_key", serviceKey),
				zap.String("client_name", clientName),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Invalid service credentials",
			})
			c.Abort()
			return
		}

		// Validate client name and secret hash
		if service.ClientName != clientName || service.SecretHash != secretHash {
			logger.Warn("Invalid credentials",
				zap.String("service_key", serviceKey),
				zap.String("client_name", clientName),
				zap.Bool("client_match", service.ClientName == clientName),
				zap.Bool("secret_match", service.SecretHash == secretHash),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Invalid service credentials",
			})
			c.Abort()
			return
		}

		// Authentication successful
		logger.Debug("Authentication successful",
			zap.String("service_key", serviceKey),
			zap.String("client_name", clientName),
		)
		c.Next()
	}
}

// ProcessTimerMiddleware logs request processing time
func ProcessTimerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate processing time
		duration := time.Since(startTime)
		durationMs := float64(duration.Nanoseconds()) / 1e6

		// Add processing time to response headers
		c.Header("X-Process-Time", duration.String())

		// Log completion
		logger.Info("[REQUEST]", zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Float64("duration_ms", durationMs),
		)
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range cfg.AllowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			// Allow all origins (matching Python implementation)
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, "+cfg.ServiceKeyName+", "+cfg.ClientNameHeader+", "+cfg.SecretHashHeader)
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
