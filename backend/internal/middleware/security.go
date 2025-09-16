package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Rate Limiter
type RateLimiter struct {
	requests map[string]*ClientRequests
	mu       sync.RWMutex
	limit    int
	window   time.Duration
	cleanup  time.Duration
}

type ClientRequests struct {
	count    int
	window   time.Time
	lastSeen time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*ClientRequests),
		limit:    limit,
		window:   window,
		cleanup:  time.Minute * 5,
	}

	// Start cleanup goroutine
	go rl.cleanup_routine()
	return rl
}

func (rl *RateLimiter) cleanup_routine() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanupExpired()
		}
	}
}

func (rl *RateLimiter) cleanupExpired() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, req := range rl.requests {
		if now.Sub(req.lastSeen) > rl.cleanup {
			delete(rl.requests, key)
		}
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	req, exists := rl.requests[key]

	if !exists {
		rl.requests[key] = &ClientRequests{
			count:    1,
			window:   now,
			lastSeen: now,
		}
		return true
	}

	req.lastSeen = now

	// Check if we're in a new window
	if now.Sub(req.window) >= rl.window {
		req.count = 1
		req.window = now
		return true
	}

	// Check if limit exceeded
	if req.count >= rl.limit {
		return false
	}

	req.count++
	return true
}

// Rate Limiting Middleware
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		// Get client identifier (IP address, could be enhanced with user ID)
		clientIP := c.ClientIP()

		if !limiter.Allow(clientIP) {
			ThrowRateLimitError(fmt.Sprintf("Rate limit exceeded. Maximum %d requests per %v", limit, window))
			return
		}

		c.Next()
	}
}

// Request Validation Middleware
func RequestValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate content length
		if c.Request.ContentLength > 10*1024*1024 { // 10MB limit
			ThrowValidationError("Request body too large")
			return
		}

		// Validate content type for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "" && !isValidContentType(contentType) {
				ThrowValidationError("Invalid content type")
				return
			}
		}

		c.Next()
	}
}

func isValidContentType(contentType string) bool {
	validTypes := []string{
		"application/json",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
		"text/plain",
	}

	for _, validType := range validTypes {
		if strings.HasPrefix(strings.ToLower(contentType), validType) {
			return true
		}
	}
	return false
}

// Security Headers Middleware
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// Request Sanitization Middleware
func RequestSanitizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize common dangerous characters in query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				sanitized := sanitizeString(value)
				c.Request.URL.Query()[key][i] = sanitized
			}
		}

		c.Next()
	}
}

func sanitizeString(input string) string {
	// Remove or escape dangerous characters
	replacer := strings.NewReplacer(
		"<script", "&lt;script",
		"</script>", "&lt;/script&gt;",
		"javascript:", "",
		"vbscript:", "",
		"onload", "",
		"onerror", "",
		"onclick", "",
	)

	return replacer.Replace(input)
}

// Request ID Middleware for tracing
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID
			requestID = fmt.Sprintf("%d-%s", time.Now().UnixNano(), c.ClientIP())
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		c.Next()
	}
}

// CORS Enhancement (builds on existing CORS middleware)
func EnhancedCORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is in allowed list
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, User-Agent, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", strconv.Itoa(12*3600)) // 12 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
