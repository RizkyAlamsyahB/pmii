package main

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/handlers"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/internal/routes"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/garuda-labs-1/pmii-be/pkg/cloudinary"
	"github.com/garuda-labs-1/pmii-be/pkg/database"
	"github.com/garuda-labs-1/pmii-be/pkg/logger"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Logger
	logger.Init()
	logger.Info.Println("🚀 Starting PMII Backend API...")

	// 2. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error.Fatalf("Failed to load configuration: %v", err)
	}
	logger.Info.Printf("Environment: %s", cfg.Server.Environment)

	// 3. Initialize JWT Secret
	utils.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	logger.Info.Printf("✅ JWT initialized (expiration: %d hours)", cfg.JWT.ExpirationHours)

	// 3a. Initialize Token Blacklist
	utils.InitBlacklist()
	logger.Info.Println("✅ Token blacklist initialized")

	// 4. Initialize Database Connection
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}

	// 4a. Run Database Migrations
	if err := database.RunMigrations(dbConfig); err != nil {
		logger.Error.Fatalf("Failed to run migrations: %v", err)
	}

	// 4b. Connect to Database
	db, err := database.InitDB(dbConfig)
	if err != nil {
		logger.Error.Fatalf("Failed to connect to database: %v", err)
	}

	// 4c. Seed Default Users (Auto-run on startup)
	if err := database.SeedDefaultUsers(db); err != nil {
		logger.Error.Fatalf("Failed to seed default users: %v", err)
	}

	// 5. Initialize Cloudinary Service
	cloudinaryService, err := cloudinary.NewService(cfg.Cloudinary.URL)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	logger.Info.Println("✅ Cloudinary service initialized")

	// 6. Initialize Repositories (Data Layer)
	userRepo := repository.NewUserRepository(db)
	testimonialRepo := repository.NewTestimonialRepository(db)

	// 7. Initialize Services (Business Logic Layer)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	testimonialService := service.NewTestimonialService(testimonialRepo, cloudinaryService)

	// 8. Initialize Handlers (Transport Layer)
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	testimonialHandler := handlers.NewTestimonialHandler(testimonialService)

	// 9. Setup Gin Router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 10. Setup Routes (dari internal/routes)
	routes.SetupRoutes(r, authHandler, adminHandler, userHandler, testimonialHandler, cfg.Server.AllowedOrigins, cfg.Server.Environment)

	// 11. Start Server
	serverAddr := ":" + cfg.Server.Port
	logger.Info.Printf("🌐 Server running on http://localhost%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		logger.Error.Fatalf("Failed to start server: %v", err)
	}
}
