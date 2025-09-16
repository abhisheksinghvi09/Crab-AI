package handlers

import (
	"context"
	"crab-ai/config"
	"crab-ai/internal/apis/dtos"
	"crab-ai/pkg/llm"
	"crab-ai/pkg/mongodb"
	"crab-ai/pkg/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	mongoClient *mongodb.MongoDBClient
	redisRepo   redis.IRedisRepositories
	llmManager  *llm.Manager
}

func NewHealthHandler(mongoClient *mongodb.MongoDBClient, redisRepo redis.IRedisRepositories, llmManager *llm.Manager) *HealthHandler {
	return &HealthHandler{
		mongoClient: mongoClient,
		redisRepo:   redisRepo,
		llmManager:  llmManager,
	}
}

type HealthResponse struct {
	Status       string                          `json:"status"`
	Timestamp    time.Time                       `json:"timestamp"`
	Environment  string                          `json:"environment"`
	Version      string                          `json:"version"`
	Dependencies map[string]DependencyHealth     `json:"dependencies"`
	LLMProviders map[string]*llm.HealthStatus    `json:"llm_providers"`
	Metrics      map[string]*llm.RequestMetrics  `json:"metrics,omitempty"`
}

type DependencyHealth struct {
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time"`
	Error        string        `json:"error,omitempty"`
}

// BasicHealthCheck provides a simple health check endpoint
func (h *HealthHandler) BasicHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Message: "Server is healthy",
		Data: map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now(),
		},
	})
}

// DetailedHealthCheck provides comprehensive health information
func (h *HealthHandler) DetailedHealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	healthResponse := HealthResponse{
		Status:       "healthy",
		Timestamp:    time.Now(),
		Environment:  config.Env.Environment,
		Version:      "1.0.0", // You can make this configurable
		Dependencies: make(map[string]DependencyHealth),
	}

	// Check MongoDB
	mongoHealth := h.checkMongoDB(ctx)
	healthResponse.Dependencies["mongodb"] = mongoHealth
	if mongoHealth.Status != "healthy" {
		healthResponse.Status = "degraded"
	}

	// Check Redis
	redisHealth := h.checkRedis(ctx)
	healthResponse.Dependencies["redis"] = redisHealth
	if redisHealth.Status != "healthy" {
		healthResponse.Status = "degraded"
	}

	// Get LLM provider health status
	healthResponse.LLMProviders = h.llmManager.GetHealthStatus()
	
	// Check if any LLM provider is unhealthy
	hasHealthyLLM := false
	for _, status := range healthResponse.LLMProviders {
		if status.IsHealthy {
			hasHealthyLLM = true
			break
		}
	}
	
	if !hasHealthyLLM {
		healthResponse.Status = "degraded"
	}

	// Include metrics if requested
	if c.Query("include_metrics") == "true" {
		healthResponse.Metrics = h.llmManager.GetMetrics()
	}

	// Set response status code based on health
	statusCode := http.StatusOK
	if healthResponse.Status == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, dtos.Response{
		Success: healthResponse.Status == "healthy",
		Message: "Health check completed",
		Data:    healthResponse,
	})
}

func (h *HealthHandler) checkMongoDB(ctx context.Context) DependencyHealth {
	start := time.Now()
	
	// Try to ping MongoDB
	err := h.mongoClient.Client.Ping(ctx, nil)
	responseTime := time.Since(start)
	
	if err != nil {
		config.Error("MongoDB health check failed:", err)
		return DependencyHealth{
			Status:       "unhealthy",
			ResponseTime: responseTime,
			Error:        err.Error(),
		}
	}
	
	return DependencyHealth{
		Status:       "healthy",
		ResponseTime: responseTime,
	}
}

func (h *HealthHandler) checkRedis(ctx context.Context) DependencyHealth {
	start := time.Now()
	
	// Try to ping Redis
	err := h.redisRepo.Ping(ctx)
	responseTime := time.Since(start)
	
	if err != nil {
		config.Error("Redis health check failed:", err)
		return DependencyHealth{
			Status:       "unhealthy",
			ResponseTime: responseTime,
			Error:        err.Error(),
		}
	}
	
	return DependencyHealth{
		Status:       "healthy",
		ResponseTime: responseTime,
	}
}

// ReadinessCheck for Kubernetes readiness probes
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check critical dependencies
	ready := true
	issues := []string{}

	// MongoDB is critical
	if err := h.mongoClient.Client.Ping(ctx, nil); err != nil {
		ready = false
		issues = append(issues, "mongodb")
	}

	// Redis is critical
	if err := h.redisRepo.Ping(ctx); err != nil {
		ready = false
		issues = append(issues, "redis")
	}

	if ready {
		c.JSON(http.StatusOK, dtos.Response{
			Success: true,
			Message: "Service is ready",
			Data:    map[string]interface{}{"ready": true},
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, dtos.Response{
			Success: false,
			Message: "Service is not ready",
			Data: map[string]interface{}{
				"ready":  false,
				"issues": issues,
			},
		})
	}
}

// LivenessCheck for Kubernetes liveness probes
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Message: "Service is alive",
		Data:    map[string]interface{}{"alive": true},
	})
}