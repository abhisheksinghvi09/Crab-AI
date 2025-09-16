package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(2, time.Second)

	// First request should be allowed
	assert.True(t, limiter.Allow("test-key"))

	// Second request should be allowed
	assert.True(t, limiter.Allow("test-key"))

	// Third request should be rejected
	assert.False(t, limiter.Allow("test-key"))

	// Different key should be allowed
	assert.True(t, limiter.Allow("other-key"))
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	// Add error handler middleware to catch panics
	router.Use(EnhancedRecoveryMiddleware())
	router.Use(ErrorHandlerMiddleware())
	router.Use(RateLimitMiddleware(1, time.Second))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// First request should succeed
	req1, _ := http.NewRequest("GET", "/test", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request should be rate limited
	req2, _ := http.NewRequest("GET", "/test", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}

func TestRequestValidationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	// Add error handler middleware to catch panics
	router.Use(EnhancedRecoveryMiddleware())
	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestValidationMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test valid content type
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SecurityHeadersMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("request_id")
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal string", "normal string"},
		{"<script>alert('xss')</script>", "&lt;script>alert('xss')&lt;/script&gt;"},
		{"javascript:alert('xss')", "alert('xss')"},
		{"onclick='alert(1)'", "='alert(1)'"},
	}

	for _, test := range tests {
		result := sanitizeString(test.input)
		assert.Equal(t, test.expected, result, "Input: %s", test.input)
	}
}

func TestAPIError(t *testing.T) {
	err := NewAPIError(ErrorTypeValidation, "test message", http.StatusBadRequest)

	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Equal(t, "test message", err.Error())

	// Test chaining
	err.WithDetails("test details").WithCause(assert.AnError)
	assert.Equal(t, "test details", err.Details)
	assert.Equal(t, assert.AnError, err.Cause)
}

func TestIsValidContentType(t *testing.T) {
	validTypes := []string{
		"application/json",
		"application/json; charset=utf-8",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
		"text/plain",
	}

	for _, contentType := range validTypes {
		assert.True(t, isValidContentType(contentType), "Content type should be valid: %s", contentType)
	}

	invalidTypes := []string{
		"application/xml",
		"text/html",
		"image/jpeg",
		"",
	}

	for _, contentType := range invalidTypes {
		assert.False(t, isValidContentType(contentType), "Content type should be invalid: %s", contentType)
	}
}
