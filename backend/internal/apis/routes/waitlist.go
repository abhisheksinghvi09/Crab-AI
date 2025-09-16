package routes

import (
	"crab-ai/internal/di"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupWaitlistRoutes(router *gin.Engine) {
	waitlistHandler, err := di.GetWaitlistHandler()
	if err != nil {
		log.Fatalf("Failed to get waitlist handler: %v", err)
	}
	enterprise := router.Group("/api/enterprise")
	{
		enterprise.POST("/waitlist", waitlistHandler.AddToWaitlist)
	}
}
