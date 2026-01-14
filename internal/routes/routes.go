package routes

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/config"
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
	publicHomeHandler *handlers.PublicHomeHandler,
	documentHandler *handlers.DocumentHandler,
	publicDocumentHandler *handlers.PublicDocumentHandler,
	dashboardHandler *handlers.DashboardHandler,
	publicSiteSettingHandler *handlers.PublicSiteSettingHandler,
	allowedOrigins string,
	environment string,
) {

	// --- Activity Log untuk Audit (inisialisasi pertama karena digunakan banyak service) ---
	activityLogRepo := repository.NewActivityLogRepository()
	activityLogSvc := service.NewActivityLogService(activityLogRepo)
	activityLogHandler := handlers.NewActivityLogHandler(activityLogSvc)

	// Inisialisasi Dependency untuk News Publik
	newsRepo := repository.NewNewsRepository(config.DB)
	newsSvc := service.NewNewsService(newsRepo)
	newsHandler := handlers.NewNewsHandler(newsSvc)

	postRepo := repository.NewPostRepository(config.DB)
	postSvc := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postSvc)

	catRepo := repository.NewCategoryRepository()
	catSvc := service.NewCategoryService(catRepo, activityLogRepo)
	catHandler := handlers.NewCategoryHandler(catSvc)

	tagRepo := repository.NewTagRepository()
	tagSvc := service.NewTagService(tagRepo, activityLogRepo)
	tagHandler := handlers.NewTagHandler(tagSvc)

	// Global Middlewares
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(allowedOrigins))
	r.Use(middleware.RequestInfoMiddleware()) // Inject IP and User-Agent for activity logging

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

		v1.GET("/news", newsHandler.GetNewsList)                   // GET /v1/news
		v1.GET("/news/:slug", postHandler.GetPost)                 // GET /v1/news/:slug
		v1.GET("/categories/:slug", newsHandler.GetNewsByCategory) // GET /v1/categories/:slug

		// Public Routes - About Page (No Authentication Required)
		v1.GET("/about", publicAboutHandler.GetAboutPage)                               // GET /v1/about
		v1.GET("/about/departments", publicAboutHandler.GetDepartments)                 // GET /v1/about/departments
		v1.GET("/about/members/:department", publicAboutHandler.GetMembersByDepartment) // GET /v1/about/members/:department

		// Public Routes - Home Page (No Authentication Required)
		v1.GET("/home/hero", publicHomeHandler.GetHeroSection)               // GET /v1/home/hero
		v1.GET("/home/latest-news", publicHomeHandler.GetLatestNewsSection)  // GET /v1/home/latest-news
		v1.GET("/home/about-us", publicHomeHandler.GetAboutUsSection)        // GET /v1/home/about-us
		v1.GET("/home/why", publicHomeHandler.GetWhySection)                 // GET /v1/home/why
		v1.GET("/home/testimonial", publicHomeHandler.GetTestimonialSection) // GET /v1/home/testimonial
		v1.GET("/home/faq", publicHomeHandler.GetFaqSection)                 // GET /v1/home/faq
		v1.GET("/home/cta", publicHomeHandler.GetCtaSection)                 // GET /v1/home/cta
		v1.GET("/documents", publicDocumentHandler.GetAllPublic)             // GET /v1/documents
		v1.GET("/documents/:type", publicDocumentHandler.GetByTypePublic)    // GET /v1/documents/:type

		// Admin Routes - Requires Admin Role (Level 1)
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole("1"))
		{
			// Dashboard Routes - Admin with Activity Logs
			adminRoutes.GET("/dashboard", dashboardHandler.GetDashboard)                // GET /v1/admin/dashboard?year=2026&month=1
			adminRoutes.GET("/dashboard/periods", dashboardHandler.GetAvailablePeriods) // GET /v1/admin/dashboard/periods

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

			// Document Routes - Admin Only
			adminRoutes.GET("/documents/types", documentHandler.GetTypes) // GET /v1/admin/documents/types
			adminRoutes.POST("/documents", documentHandler.Create)        // POST /v1/admin/documents
			adminRoutes.GET("/documents", documentHandler.GetAll)         // GET /v1/admin/documents
			adminRoutes.GET("/documents/:id", documentHandler.GetByID)    // GET /v1/admin/documents/:id
			adminRoutes.PUT("/documents/:id", documentHandler.Update)     // PUT /v1/admin/documents/:id
			adminRoutes.DELETE("/documents/:id", documentHandler.Delete)  // DELETE /v1/admin/documents/:id

			// Activity Log Routes - Admin Only
			adminRoutes.GET("/activity-logs", activityLogHandler.GetActivityLogs) // GET /v1/admin/activity-logs
		}

		// User Routes - Requires Authentication (Any authenticated user)
		userRoutes := v1.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			// GET /v1/users/me - Get own profile
			userRoutes.GET("/me", userHandler.GetMyProfile)

			// Dashboard Routes - Author without Activity Logs
			userRoutes.GET("/dashboard", dashboardHandler.GetDashboard)                // GET /v1/users/dashboard?year=2026&month=1
			userRoutes.GET("/dashboard/periods", dashboardHandler.GetAvailablePeriods) // GET /v1/users/dashboard/periods
		}

		// Public Posts Routes - Untuk pengunjung melihat postingan
		posts := v1.Group("/posts")
		{
			posts.GET("", postHandler.GetPosts)
			posts.GET("/:id", postHandler.GetPost)
		}

		// Protected Posts Routes - Requires Admin or Author
		postsProtected := v1.Group("/posts")
		postsProtected.Use(middleware.AuthMiddleware(), middleware.RequireAnyRole("1", "2"))
		{
			postsProtected.POST("", postHandler.CreatePost)
			postsProtected.PUT("/:id", postHandler.UpdatePost)
			postsProtected.DELETE("/:id", postHandler.DeletePost)
		}

		// Public Categories Routes
		categories := v1.Group("/categories")
		{
			categories.GET("", catHandler.GetCategories)
		}

		// Protected Categories Routes - Requires Admin
		categoriesProtected := v1.Group("/categories")
		categoriesProtected.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
		{
			categoriesProtected.POST("", catHandler.CreateCategory)
			categoriesProtected.PUT("/:id", catHandler.UpdateCategory)
			categoriesProtected.DELETE("/:id", catHandler.DeleteCategory)
		}

		// Public Tags Routes
		tags := v1.Group("/tags")
		{
			tags.GET("", tagHandler.GetTags)
		}

		// Protected Tags Routes - Requires Admin
		tagsProtected := v1.Group("/tags")
		tagsProtected.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
		{
			tagsProtected.POST("", tagHandler.CreateTag)
			tagsProtected.PUT("/:id", tagHandler.UpdateTag)
			tagsProtected.DELETE("/:id", tagHandler.DeleteTag)
		}

	}
}
