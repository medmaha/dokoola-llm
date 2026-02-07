package handlers

import (
	"net/http"

	"github.com/dokoola/llm-go/internal/models"
	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthCheckResponse{
		Status:  "ok",
		Message: "Dokoola LLM Service is running",
	})
}
