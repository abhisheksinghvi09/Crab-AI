package routes

import (
	"crab-ai/internal/apis/dtos"
	"crab-ai/internal/di"
	"crab-ai/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRoutes(router *gin.Engine) {
	// Add recovery middleware
	router.Use(middleware.CustomRecoveryMiddleware())

	// Get health handler
	healthHandler, err := di.GetHealthHandler()
	if err != nil {
		log.Fatalf("Failed to get health handler: %v", err)
	}

	// Health check routes
	router.GET("/health", healthHandler.BasicHealthCheck)
	router.GET("/health/detailed", healthHandler.DetailedHealthCheck)
	router.GET("/health/ready", healthHandler.ReadinessCheck)
	router.GET("/health/live", healthHandler.LivenessCheck)

	// Legacy health check route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, dtos.Response{
			Success: true,
			Data:    "Server is healthy!",
		})
	})

	githubHandler, err := di.GetGitHubHandler()
	if err != nil {
		log.Fatalf("Failed to get github handler: %v", err)
	}
	// Github repository statistics route
	router.GET("/api/github/stats", githubHandler.GetGitHubStats)

	// Setup all route groups
	SetupAuthRoutes(router)
	SetupChatRoutes(router)
	SetupWaitlistRoutes(router)
	SetupUploadRoutes(router)
	SetupGoogleOAuthRoutes(router)
}
