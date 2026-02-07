package models

// HealthCheckResponse is the response for health check endpoint
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
