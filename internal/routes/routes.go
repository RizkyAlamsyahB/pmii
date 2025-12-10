package routes

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/handlers"
	"github.com/garuda-labs-1/pmii-be/internal/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SetupRoutes mengatur semua routing untuk aplikasi
func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	adminHandler *handlers.AdminHandler,
	userHandler *handlers.UserHandler,
	allowedOrigins string,
	environment string,
) {
	// Global Middlewares
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(allowedOrigins))

	// Rate Limiter untuk login endpoint (60 request per menit = 1 req/s, burst 60)
	loginLimiter := middleware.NewRateLimiter(rate.Limit(1), 60)

	// Health Check Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"message":     "🚀 PMII Backend API is running!",
			"service":     "Backend Go Gin - Clean Architecture",
			"environment": environment,
			"version":     "1.0.0",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// API Routes - Base URL: /v1 (Development)
	v1 := r.Group("/v1")
	{
		// Public Routes - Authentication
		auth := v1.Group("/auth")
		{
			// Login dengan rate limiter (60 req/menit per IP)
			auth.POST("/login", loginLimiter.Limit(), authHandler.Login)

			// Logout (butuh auth)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
		}

		// Admin Routes - Requires Admin Role (Level 1)
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole("1"))
		{
			// GET /v1/admin/dashboard - Admin dashboard
			adminRoutes.GET("/dashboard", adminHandler.GetDashboard)

			// GET /v1/admin/users - List all users (Admin only)
			adminRoutes.GET("/users", adminHandler.GetAllUsers)
		}

		// User Routes - Requires Authentication (Any authenticated user)
		userRoutes := v1.Group("/user")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			// GET /v1/user/dashboard - User dashboard
			userRoutes.GET("/dashboard", userHandler.GetDashboard)

			// GET /v1/user/profile - Get own profile
			userRoutes.GET("/profile", userHandler.GetProfile)
		}
	}
}
