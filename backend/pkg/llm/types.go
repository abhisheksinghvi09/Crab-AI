package llm

import (
	"context"
	"crab-ai/internal/models"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string                 `json:"role"`
	Content string                 `json:"content"`
	Type    string                 `json:"type,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Client defines the interface for LLM interactions
type Client interface {
	GenerateResponse(ctx context.Context, messages []*models.LLMMessage, dbType string, nonTechMode bool) (string, error)
	GenerateRecommendations(ctx context.Context, messages []*models.LLMMessage, dbType string) (string, error)
	GetModelInfo() ModelInfo
	HealthCheck(ctx context.Context) error
}

// ModelInfo contains information about the LLM model
type ModelInfo struct {
	Name                string
	Provider            string
	MaxCompletionTokens int
	ContextLimit        int
}

// Config holds configuration for LLM clients
type Config struct {
	Provider            string
	Model               string
	APIKey              string
	MaxCompletionTokens int
	Temperature         float64
	DBConfigs           []LLMDBConfig
	HealthCheckInterval time.Duration
	MaxRetries          int
	RequestTimeout      time.Duration
}

type LLMDBConfig struct {
	DBType       string
	Schema       interface{}
	SystemPrompt string
}

// HealthStatus represents the health status of an LLM provider
type HealthStatus struct {
	Provider     string
	IsHealthy    bool
	LastChecked  time.Time
	Error        error
	ResponseTime time.Duration
}

// RequestMetrics holds metrics for LLM requests
type RequestMetrics struct {
	Provider            string
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	AverageResponseTime time.Duration
	LastRequestTime     time.Time
}
