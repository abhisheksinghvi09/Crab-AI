package main

import (
	"context"
	"crab-ai/config"
	"crab-ai/internal/apis/routes"
	"crab-ai/internal/di"
	"crab-ai/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Initialize logger
	config.InitLogger()

	// Initialize dependencies
	di.Initialize()

	// Setup Gin
	ginApp := gin.New() // Use gin.New() instead of gin.Default()

	// Add enhanced middleware stack
	setupMiddleware(ginApp)

	// Setup routes
	routes.SetupDefaultRoutes(ginApp)

	// Create server
	srv := &http.Server{
		Addr:         ":" + config.Env.Port,
		Handler:      ginApp,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		config.Info("Starting server on port", config.Env.Port)
		fmt.Printf("✨ Welcome to NeoBase! Running in %s Mode\n", config.Env.Environment)
		fmt.Printf("📱 Client UI: %s\n", config.Env.CorsAllowedOrigin)
		if config.Env.LandingPageCorsAllowedOrigin != "" {
			fmt.Printf("🌐 Landing Page: %s\n", config.Env.LandingPageCorsAllowedOrigin)
		}

		allowedOrigins := []string{config.Env.CorsAllowedOrigin}
		if config.Env.LandingPageCorsAllowedOrigin != "" {
			allowedOrigins = append(allowedOrigins, config.Env.LandingPageCorsAllowedOrigin)
		}
		fmt.Printf("🚀 CORS Origins: %v\n", allowedOrigins)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Fatal("NeoBase failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	config.Info("🔻 NeoBase is shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		config.Fatal("NeoBase forced to shutdown:", err)
	}

	config.Info("👋 NeoBase has been shut down successfully")
}

func setupMiddleware(ginApp *gin.Engine) {
	// Request ID middleware (first, for tracing)
	ginApp.Use(middleware.RequestIDMiddleware())

	// Enhanced recovery middleware
	ginApp.Use(middleware.EnhancedRecoveryMiddleware())

	// Error handler middleware
	ginApp.Use(middleware.ErrorHandlerMiddleware())

	// Security headers
	ginApp.Use(middleware.SecurityHeadersMiddleware())

	// Request validation
	ginApp.Use(middleware.RequestValidationMiddleware())

	// Request sanitization
	ginApp.Use(middleware.RequestSanitizationMiddleware())

	// Rate limiting (100 requests per minute per IP)
	ginApp.Use(middleware.RateLimitMiddleware(100, time.Minute))

	// Add logging middleware
	ginApp.Use(gin.Logger())

	// Build allowed origins list
	allowedOrigins := []string{config.Env.CorsAllowedOrigin}
	if config.Env.LandingPageCorsAllowedOrigin != "" {
		allowedOrigins = append(allowedOrigins, config.Env.LandingPageCorsAllowedOrigin)
	}

	// Enhanced CORS
	ginApp.Use(middleware.EnhancedCORSMiddleware(allowedOrigins))

	// Fallback to original CORS configuration for compatibility
	ginApp.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"User-Agent",
			"Referer",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
			"X-Request-ID",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
