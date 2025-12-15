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
	testimonialHandler *handlers.TestimonialHandler,
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

			// Testimonial Routes - Admin Only
			adminRoutes.POST("/testimonials", testimonialHandler.Create)       // POST /v1/admin/testimonials
			adminRoutes.GET("/testimonials", testimonialHandler.GetAll)        // GET /v1/admin/testimonials
			adminRoutes.GET("/testimonials/:id", testimonialHandler.GetByID)   // GET /v1/admin/testimonials/:id
			adminRoutes.PUT("/testimonials/:id", testimonialHandler.Update)    // PUT /v1/admin/testimonials/:id
			adminRoutes.DELETE("/testimonials/:id", testimonialHandler.Delete) // DELETE /v1/admin/testimonials/:id
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

		posts := v1.Group("/posts")
		{
			posts.POST("", handlers.CreatePost)       // Create
			posts.GET("", handlers.GetPosts)          // Read All
			posts.GET("/:id", handlers.GetPost)       // Read One
			posts.PUT("/:id", handlers.UpdatePost)    // Update
			posts.DELETE("/:id", handlers.DeletePost) // Delete
		}
		categories := v1.Group("/categories")
		{
			// Create (POST /v1/categories)
			categories.POST("", handlers.CreateCategory)

			// Read All (GET /v1/categories)
			categories.GET("", handlers.GetCategories)

			// Update (PUT /v1/categories/:id) -> Jika ingin pakai
			categories.PUT("/:id", handlers.UpdateCategory)

			// Delete (DELETE /v1/categories/:id)
			categories.DELETE("/:id", handlers.DeleteCategory)
		}

		tags := v1.Group("/tags")
		{
			// Create (POST /v1/tags)
			tags.POST("", handlers.CreateTag)

			// Read All (GET /v1/tags)
			tags.GET("", handlers.GetTags)

			// Update (PUT /v1/tags/:id)
			tags.PUT("/:id", handlers.UpdateTag)

			// Delete (DELETE /v1/tags/:id)
			tags.DELETE("/:id", handlers.DeleteTag)
		}

	}

}
