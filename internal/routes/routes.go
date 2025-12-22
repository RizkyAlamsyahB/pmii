package routes

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/handlers"
	"github.com/garuda-labs-1/pmii-be/internal/middleware"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/internal/service"
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
	memberHandler *handlers.MemberHandler,
	aboutHandler *handlers.AboutHandler,
	siteSettingHandler *handlers.SiteSettingHandler,
	contactHandler *handlers.ContactHandler,
	publicAboutHandler *handlers.PublicAboutHandler,
	allowedOrigins string,
	environment string,
) {

	// --- Inisialisasi Modul Post (Clean Architecture) ---
	postRepo := repository.NewPostRepository()
	postSvc := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postSvc)

	catRepo := repository.NewCategoryRepository()
	catSvc := service.NewCategoryService(catRepo)
	catHandler := handlers.NewCategoryHandler(catSvc)

	tagRepo := repository.NewTagRepository()
	tagSvc := service.NewTagService(tagRepo)
	tagHandler := handlers.NewTagHandler(tagSvc)

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

			// Ubah/ganti password
			auth.POST("/change-password", middleware.AuthMiddleware(), authHandler.ChangePassword)
		}

		// Public Routes - About Page (No Authentication Required)
		v1.GET("/about", publicAboutHandler.GetAboutPage)                               // GET /v1/about
		v1.GET("/about/departments", publicAboutHandler.GetDepartments)                 // GET /v1/about/departments
		v1.GET("/about/members/:department", publicAboutHandler.GetMembersByDepartment) // GET /v1/about/members/:department

		// Admin Routes - Requires Admin Role (Level 1)
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole("1"))
		{
			// GET /v1/admin/dashboard - Admin dashboard
			adminRoutes.GET("/dashboard", adminHandler.GetDashboard)

			// Testimonial Routes - Admin Only
			adminRoutes.POST("/testimonials", testimonialHandler.Create)       // POST /v1/admin/testimonials
			adminRoutes.GET("/testimonials", testimonialHandler.GetAll)        // GET /v1/admin/testimonials
			adminRoutes.GET("/testimonials/:id", testimonialHandler.GetByID)   // GET /v1/admin/testimonials/:id
			adminRoutes.PUT("/testimonials/:id", testimonialHandler.Update)    // PUT /v1/admin/testimonials/:id
			adminRoutes.DELETE("/testimonials/:id", testimonialHandler.Delete) // DELETE /v1/admin/testimonials/:id

			// Member Routes - Admin Only
			adminRoutes.POST("/members", memberHandler.Create)       // POST /v1/admin/members
			adminRoutes.GET("/members", memberHandler.GetAll)        // GET /v1/admin/members
			adminRoutes.GET("/members/:id", memberHandler.GetByID)   // GET /v1/admin/members/:id
			adminRoutes.PUT("/members/:id", memberHandler.Update)    // PUT /v1/admin/members/:id
			adminRoutes.DELETE("/members/:id", memberHandler.Delete) // DELETE /v1/admin/members/:id

			// User Management Routes - Admin Only
			adminRoutes.GET("/users", userHandler.GetAllUsers)           // GET /v1/admin/users
			adminRoutes.GET("/users/:id", userHandler.GetUserByID)       // GET /v1/admin/users/:id
			adminRoutes.POST("/users", userHandler.CreateUser)           // POST /v1/admin/users
			adminRoutes.PUT("/users/:id", userHandler.UpdateUserByID)    // PUT /v1/admin/users/:id
			adminRoutes.DELETE("/users/:id", userHandler.DeleteUserByID) // DELETE /v1/admin/users/:id
			// About Routes - Admin Only (singleton - only GET and PUT)
			adminRoutes.GET("/about", aboutHandler.Get)    // GET /v1/admin/about
			adminRoutes.PUT("/about", aboutHandler.Update) // PUT /v1/admin/about

			// Site Settings Routes - Admin Only (singleton - only GET and PUT)
			adminRoutes.GET("/settings", siteSettingHandler.Get)    // GET /v1/admin/settings
			adminRoutes.PUT("/settings", siteSettingHandler.Update) // PUT /v1/admin/settings

			// Contact Routes - Admin Only (singleton - only GET and PUT)
			adminRoutes.GET("/contact", contactHandler.Get)    // GET /v1/admin/contact
			adminRoutes.PUT("/contact", contactHandler.Update) // PUT /v1/admin/contact
		}

		// User Routes - Requires Authentication (Any authenticated user)
		userRoutes := v1.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			// GET /v1/users/me - Get own profile
			userRoutes.GET("/me", userHandler.GetMyProfile)
		}

		posts := v1.Group("/posts")
		{
			posts.GET("", postHandler.GetPosts)
			posts.POST("", postHandler.CreatePost)
			posts.GET("/:id", postHandler.GetPost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
		}

		categories := v1.Group("/categories")
		{
			categories.GET("", catHandler.GetCategories)
			categories.POST("", catHandler.CreateCategory)
			categories.PUT("/:id", catHandler.UpdateCategory)
			categories.DELETE("/:id", catHandler.DeleteCategory)
		}

		tags := v1.Group("/tags")
		{
			tags.GET("", tagHandler.GetTags)
			tags.POST("", tagHandler.CreateTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.DELETE("/:id", tagHandler.DeleteTag)
		}

	}

}
