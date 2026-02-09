package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dokoola/llm-go/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestAuthMiddleware_SkipsHealthCheck(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		ServiceKeyName:   "X-Service-Key",
		ClientNameHeader: "X-Client-Name",
		SecretHashHeader: "X-Secret-Hash",
		AllowedServices:  make(map[string]config.ServiceConfig),
	}

	router := gin.New()
	router.Use(AuthMiddleware(cfg, logger))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_RejectsMissingHeaders(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		ServiceKeyName:   "X-Service-Key",
		ClientNameHeader: "X-Client-Name",
		SecretHashHeader: "X-Secret-Hash",
		AllowedServices:  make(map[string]config.ServiceConfig),
	}

	router := gin.New()
	router.Use(AuthMiddleware(cfg, logger))
	router.POST("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestAuthMiddleware_RejectsInvalidServiceKey(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		ServiceKeyName:   "X-Service-Key",
		ClientNameHeader: "X-Client-Name",
		SecretHashHeader: "X-Secret-Hash",
		AllowedServices: map[string]config.ServiceConfig{
			"valid-key": {
				ClientName: "test-client",
				SecretHash: "test-hash",
			},
		},
	}

	router := gin.New()
	router.Use(AuthMiddleware(cfg, logger))
	router.POST("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/api/v1/test", nil)
	req.Header.Set("X-Service-Key", "invalid-key")
	req.Header.Set("X-Client-Name", "test-client")
	req.Header.Set("X-Secret-Hash", "test-hash")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestAuthMiddleware_RejectsInvalidCredentials(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		ServiceKeyName:   "X-Service-Key",
		ClientNameHeader: "X-Client-Name",
		SecretHashHeader: "X-Secret-Hash",
		AllowedServices: map[string]config.ServiceConfig{
			"valid-key": {
				ClientName: "test-client",
				SecretHash: "test-hash",
			},
		},
	}

	router := gin.New()
	router.Use(AuthMiddleware(cfg, logger))
	router.POST("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/api/v1/test", nil)
	req.Header.Set("X-Service-Key", "valid-key")
	req.Header.Set("X-Client-Name", "wrong-client")
	req.Header.Set("X-Secret-Hash", "test-hash")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestAuthMiddleware_AllowsValidCredentials(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		ServiceKeyName:   "X-Service-Key",
		ClientNameHeader: "X-Client-Name",
		SecretHashHeader: "X-Secret-Hash",
		AllowedServices: map[string]config.ServiceConfig{
			"valid-key": {
				ClientName: "test-client",
				SecretHash: "test-hash",
			},
		},
	}

	router := gin.New()
	router.Use(AuthMiddleware(cfg, logger))
	router.POST("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/api/v1/test", nil)
	req.Header.Set("X-Service-Key", "valid-key")
	req.Header.Set("X-Client-Name", "test-client")
	req.Header.Set("X-Secret-Hash", "test-hash")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestProcessTimerMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	router := gin.New()
	router.Use(ProcessTimerMiddleware(logger))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check if X-Process-Time header is set
	if w.Header().Get("X-Process-Time") == "" {
		t.Error("expected X-Process-Time header to be set")
	}
}
