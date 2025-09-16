package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	clients             map[string]Client
	healthStatus        map[string]*HealthStatus
	metrics             map[string]*RequestMetrics
	mu                  sync.RWMutex
	healthCheckInterval time.Duration
	stopHealthCheck     chan bool
}

func NewManager() *Manager {
	m := &Manager{
		clients:             make(map[string]Client),
		healthStatus:        make(map[string]*HealthStatus),
		metrics:             make(map[string]*RequestMetrics),
		healthCheckInterval: 30 * time.Second,
		stopHealthCheck:     make(chan bool),
	}

	// Start health check routine
	go m.healthCheckRoutine()

	return m
}

func (m *Manager) RegisterClient(name string, config Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var client Client
	var err error

	// Set default values for enhanced config
	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 30 * time.Second
	}

	switch config.Provider {
	case "openai":
		client, err = NewOpenAIClient(config)
	case "gemini":
		client, err = NewGeminiClient(config)
	default:
		return fmt.Errorf("unsupported LLM provider: %s", config.Provider)
	}

	if err != nil {
		return fmt.Errorf("failed to create LLM client: %v", err)
	}

	m.clients[name] = client
	m.healthStatus[name] = &HealthStatus{
		Provider:    config.Provider,
		IsHealthy:   true, // Assume healthy initially
		LastChecked: time.Now(),
	}
	m.metrics[name] = &RequestMetrics{
		Provider: config.Provider,
	}

	return nil
}

func (m *Manager) GetClient(name string) (Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("LLM client not found: %s", name)
	}

	// Check if client is healthy
	if status, ok := m.healthStatus[name]; ok && !status.IsHealthy {
		return nil, fmt.Errorf("LLM client %s is currently unhealthy: %v", name, status.Error)
	}

	return client, nil
}

// GetHealthyClient returns the first healthy client, implementing basic failover
func (m *Manager) GetHealthyClient() (Client, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, status := range m.healthStatus {
		if status.IsHealthy {
			if client, exists := m.clients[name]; exists {
				return client, name, nil
			}
		}
	}

	return nil, "", fmt.Errorf("no healthy LLM clients available")
}

// GetClientWithFallback tries to get the specified client, falls back to any healthy client
func (m *Manager) GetClientWithFallback(preferredName string) (Client, string, error) {
	// Try preferred client first
	client, err := m.GetClient(preferredName)
	if err == nil {
		return client, preferredName, nil
	}

	// Fall back to any healthy client
	return m.GetHealthyClient()
}

func (m *Manager) RemoveClient(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients, name)
	delete(m.healthStatus, name)
	delete(m.metrics, name)
}

// GetHealthStatus returns the health status of all clients
func (m *Manager) GetHealthStatus() map[string]*HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*HealthStatus)
	for name, status := range m.healthStatus {
		// Copy the status to avoid external mutation
		result[name] = &HealthStatus{
			Provider:     status.Provider,
			IsHealthy:    status.IsHealthy,
			LastChecked:  status.LastChecked,
			Error:        status.Error,
			ResponseTime: status.ResponseTime,
		}
	}
	return result
}

// GetMetrics returns metrics for all clients
func (m *Manager) GetMetrics() map[string]*RequestMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*RequestMetrics)
	for name, metrics := range m.metrics {
		// Copy the metrics to avoid external mutation
		result[name] = &RequestMetrics{
			Provider:            metrics.Provider,
			TotalRequests:       metrics.TotalRequests,
			SuccessfulRequests:  metrics.SuccessfulRequests,
			FailedRequests:      metrics.FailedRequests,
			AverageResponseTime: metrics.AverageResponseTime,
			LastRequestTime:     metrics.LastRequestTime,
		}
	}
	return result
}

// RecordRequest records metrics for a request
func (m *Manager) RecordRequest(clientName string, success bool, responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metrics, exists := m.metrics[clientName]; exists {
		metrics.TotalRequests++
		metrics.LastRequestTime = time.Now()

		if success {
			metrics.SuccessfulRequests++
		} else {
			metrics.FailedRequests++
		}

		// Update average response time (simple moving average)
		if metrics.TotalRequests == 1 {
			metrics.AverageResponseTime = responseTime
		} else {
			metrics.AverageResponseTime = (metrics.AverageResponseTime + responseTime) / 2
		}
	}
}

// Health check routine
func (m *Manager) healthCheckRoutine() {
	ticker := time.NewTicker(m.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.performHealthChecks()
		case <-m.stopHealthCheck:
			return
		}
	}
}

func (m *Manager) performHealthChecks() {
	m.mu.RLock()
	clients := make(map[string]Client)
	for name, client := range m.clients {
		clients[name] = client
	}
	m.mu.RUnlock()

	for name, client := range clients {
		go func(name string, client Client) {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := client.HealthCheck(ctx)
			responseTime := time.Since(start)

			m.mu.Lock()
			if status, exists := m.healthStatus[name]; exists {
				status.IsHealthy = (err == nil)
				status.LastChecked = time.Now()
				status.Error = err
				status.ResponseTime = responseTime
			}
			m.mu.Unlock()

			if err != nil {
				log.Printf("Health check failed for LLM client %s: %v", name, err)
			}
		}(name, client)
	}
}

// Stop the health check routine
func (m *Manager) Stop() {
	close(m.stopHealthCheck)
}

// Add helper function to properly format assistant response
func formatAssistantResponse(response map[string]interface{}) string {
	// Convert the response to JSON string
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error formatting assistant response: %v", err)
		return fmt.Sprintf("%v", response)
	}
	return string(jsonBytes)
}

// Helper functions
func mapRole(role string) string {
	switch strings.ToLower(role) {
	case "user":
		return "user"
	case "assistant":
		return "assistant"
	case "system":
		return "system"
	default:
		return "user"
	}
}
