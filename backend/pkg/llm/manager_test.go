package llm

import (
	"context"
	"crab-ai/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the LLM Client interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) GenerateResponse(ctx context.Context, messages []*models.LLMMessage, dbType string, nonTechMode bool) (string, error) {
	args := m.Called(ctx, messages, dbType, nonTechMode)
	return args.String(0), args.Error(1)
}

func (m *MockClient) GenerateRecommendations(ctx context.Context, messages []*models.LLMMessage, dbType string) (string, error) {
	args := m.Called(ctx, messages, dbType)
	return args.String(0), args.Error(1)
}

func (m *MockClient) GetModelInfo() ModelInfo {
	args := m.Called()
	return args.Get(0).(ModelInfo)
}

func (m *MockClient) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestManager_RegisterClient(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	config := Config{
		Provider:            "test",
		Model:               "test-model",
		APIKey:              "test-key",
		MaxCompletionTokens: 1000,
		Temperature:         0.7,
	}

	// Test unsupported provider
	err := manager.RegisterClient("test", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported LLM provider")
}

func TestManager_GetClient(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	// Test getting non-existent client
	client, err := manager.GetClient("non-existent")
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "LLM client not found")
}

func TestManager_GetHealthyClient(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	// Test when no clients are registered
	client, name, err := manager.GetHealthyClient()
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Empty(t, name)
	assert.Contains(t, err.Error(), "no healthy LLM clients available")
}

func TestManager_GetClientWithFallback(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	// Test fallback when preferred client doesn't exist
	client, name, err := manager.GetClientWithFallback("preferred")
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Empty(t, name)
}

func TestManager_HealthStatus(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	// Manually add a mock health status for testing
	manager.mu.Lock()
	manager.healthStatus["test"] = &HealthStatus{
		Provider:    "test",
		IsHealthy:   true,
		LastChecked: time.Now(),
	}
	manager.mu.Unlock()

	statuses := manager.GetHealthStatus()
	assert.Len(t, statuses, 1)
	assert.Contains(t, statuses, "test")
	assert.True(t, statuses["test"].IsHealthy)
	assert.Equal(t, "test", statuses["test"].Provider)
}

func TestManager_Metrics(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	// Test recording metrics
	manager.mu.Lock()
	manager.metrics["test"] = &RequestMetrics{
		Provider: "test",
	}
	manager.mu.Unlock()

	manager.RecordRequest("test", true, 100*time.Millisecond)

	metrics := manager.GetMetrics()
	assert.Len(t, metrics, 1)
	assert.Contains(t, metrics, "test")
	assert.Equal(t, int64(1), metrics["test"].TotalRequests)
	assert.Equal(t, int64(1), metrics["test"].SuccessfulRequests)
	assert.Equal(t, int64(0), metrics["test"].FailedRequests)
}

func TestManager_RecordRequest(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	// Setup metrics
	manager.mu.Lock()
	manager.metrics["test"] = &RequestMetrics{
		Provider: "test",
	}
	manager.mu.Unlock()

	// Test successful request
	manager.RecordRequest("test", true, 100*time.Millisecond)
	metrics := manager.GetMetrics()
	assert.Equal(t, int64(1), metrics["test"].TotalRequests)
	assert.Equal(t, int64(1), metrics["test"].SuccessfulRequests)

	// Test failed request
	manager.RecordRequest("test", false, 200*time.Millisecond)
	metrics = manager.GetMetrics()
	assert.Equal(t, int64(2), metrics["test"].TotalRequests)
	assert.Equal(t, int64(1), metrics["test"].FailedRequests)
}

func TestManager_RemoveClient(t *testing.T) {
	manager := NewManager()
	defer manager.Stop()

	// Manually add data to test removal
	manager.mu.Lock()
	manager.healthStatus["test"] = &HealthStatus{Provider: "test"}
	manager.metrics["test"] = &RequestMetrics{Provider: "test"}
	manager.mu.Unlock()

	manager.RemoveClient("test")

	statuses := manager.GetHealthStatus()
	metrics := manager.GetMetrics()
	assert.Empty(t, statuses)
	assert.Empty(t, metrics)
}
