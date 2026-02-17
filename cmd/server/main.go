package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dokoola/llm-go/internal/clients"
	"github.com/dokoola/llm-go/internal/config"
	"github.com/dokoola/llm-go/internal/handlers"
	"github.com/dokoola/llm-go/internal/llm"
	"github.com/dokoola/llm-go/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Initialize logger
	logger, err := initLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Dokoola LLM Service (Go)")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded\n",
		zap.Int("port", cfg.Settings.Port),
		zap.String("host", cfg.Settings.Host),
		zap.String("app_name", cfg.Settings.AppName),
		zap.String("version", cfg.Settings.AppVersion),
		zap.Bool("debug", cfg.Settings.Debug),
		zap.String("api_prefix", cfg.Settings.APIPrefix),
		zap.Int("allowed_services", len(cfg.AllowedServices)),
	)

	// Set Gin mode
	if !cfg.Settings.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize clients
	llmClient := llm.NewClient(cfg.CerebrasAPIKey, logger)
	backendClient := clients.NewBackendClient(cfg.BackendServerAPI, logger)

	// Initialize handlers
	jobsHandler := handlers.NewJobsHandler(llmClient, backendClient, logger)
	textCompletionHandler := handlers.NewTextCompletionHandler(llmClient, backendClient, logger)
	promptsHandler := handlers.NewPromptsHandler(llmClient, backendClient, logger)

	// Create router
	router := gin.New()
	router.RedirectTrailingSlash = true

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.ProcessTimerMiddleware(logger))
	router.Use(middleware.CORSMiddleware(cfg))

	// API routes with authentication
	apiPrefix := cfg.Settings.APIPrefix
	
	// Health check endpoint (no auth required)
	router.GET(apiPrefix+"/health", handlers.HealthCheck)

	api := router.Group(apiPrefix)
	api.Use(middleware.AuthMiddleware(cfg, logger))
	{
		// Jobs description
		api.POST("/jobs/describe", jobsHandler.GenerateJobDesc)
		
		// Jobs categorization
		api.POST("/jobs/categorize", jobsHandler.CategorizeJobs)

		// Text completion
		api.POST("/chat/completion", textCompletionHandler.Complete)

		// Prompt generation
		api.POST("/actions/generate-prompt", promptsHandler.GeneratePrompt)
	}

	// Create server
	addr := fmt.Sprintf("%s:%d", cfg.Settings.Host, cfg.Settings.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully",
		zap.String("address", addr),
		zap.String("health_check", "http://"+addr+"/health"),
	)

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// initLogger initializes the zap logger
func initLogger() (*zap.Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")

	var logger *zap.Logger
	var err error

	if logLevel == "production" || logLevel == "prod" {
		logger, err = zap.NewProduction()
	} else {
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

		// Adjust log level based on environment
		switch logLevel {
		case "debug":
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "warn", "warning":
			config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		}

		logger, err = config.Build()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}
