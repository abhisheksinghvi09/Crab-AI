package middleware

import (
	"crab-ai/config"
	"crab-ai/internal/apis/dtos"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type ErrorType string

const (
	ErrorTypeValidation  ErrorType = "validation_error"
	ErrorTypeAuth        ErrorType = "authentication_error"
	ErrorTypePermission  ErrorType = "permission_error"
	ErrorTypeNotFound    ErrorType = "not_found_error"
	ErrorTypeInternal    ErrorType = "internal_error"
	ErrorTypeRateLimit   ErrorType = "rate_limit_error"
	ErrorTypeExternalAPI ErrorType = "external_api_error"
)

type APIError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"-"`
	Cause      error     `json:"-"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(errType ErrorType, message string, statusCode int) *APIError {
	return &APIError{
		Type:       errType,
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

func (e *APIError) WithCause(cause error) *APIError {
	e.Cause = cause
	return e
}

// Enhanced Recovery Middleware with structured error responses
func EnhancedRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				config.Error("Panic recovered:", err)
				config.Error("Stack trace:", string(debug.Stack()))

				// Check if it's our structured API error
				if apiErr, ok := err.(*APIError); ok {
					handleAPIError(c, apiErr)
					return
				}

				// Handle generic panic
				apiErr := &APIError{
					Type:       ErrorTypeInternal,
					Message:    "Internal server error",
					StatusCode: http.StatusInternalServerError,
				}

				if config.Env.Environment != "PRODUCTION" {
					if errStr, ok := err.(string); ok {
						apiErr.Details = errStr
					} else if errObj, ok := err.(error); ok {
						apiErr.Details = errObj.Error()
					}
				}

				handleAPIError(c, apiErr)
			}
		}()

		c.Next()
	}
}

// Error Handler Middleware
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			lastError := c.Errors.Last()

			// Check if it's our structured API error
			if apiErr, ok := lastError.Err.(*APIError); ok {
				handleAPIError(c, apiErr)
				return
			}

			// Handle generic error
			apiErr := &APIError{
				Type:       ErrorTypeInternal,
				Message:    "An error occurred while processing your request",
				StatusCode: http.StatusInternalServerError,
			}

			if config.Env.Environment != "PRODUCTION" {
				apiErr.Details = lastError.Error()
			}

			handleAPIError(c, apiErr)
		}
	}
}

func handleAPIError(c *gin.Context, apiErr *APIError) {
	// Log the error
	if apiErr.Cause != nil {
		config.Error("API Error:", apiErr.Message, "Cause:", apiErr.Cause.Error())
	} else {
		config.Error("API Error:", apiErr.Message)
	}

	// Prevent duplicate responses
	if c.Writer.Written() {
		return
	}

	c.JSON(apiErr.StatusCode, dtos.Response{
		Success: false,
		Message: apiErr.Message,
		Data: map[string]interface{}{
			"error_type": apiErr.Type,
			"details":    apiErr.Details,
		},
	})
	c.Abort()
}

// Helper functions for common errors
func ThrowValidationError(message string) {
	panic(NewAPIError(ErrorTypeValidation, message, http.StatusBadRequest))
}

func ThrowAuthError(message string) {
	panic(NewAPIError(ErrorTypeAuth, message, http.StatusUnauthorized))
}

func ThrowPermissionError(message string) {
	panic(NewAPIError(ErrorTypePermission, message, http.StatusForbidden))
}

func ThrowNotFoundError(message string) {
	panic(NewAPIError(ErrorTypeNotFound, message, http.StatusNotFound))
}

func ThrowInternalError(message string) {
	panic(NewAPIError(ErrorTypeInternal, message, http.StatusInternalServerError))
}

func ThrowRateLimitError(message string) {
	panic(NewAPIError(ErrorTypeRateLimit, message, http.StatusTooManyRequests))
}
