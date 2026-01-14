package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/pkg/database"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
	"gorm.io/gorm"
)

// Development user data
var devUsers = []struct {
	ID       int
	FullName string
	Email    string
}{
	{3, "Ahmad Fauzi", "ahmad.fauzi@pmii.id"},
	{4, "Siti Nurhaliza", "siti.nurhaliza@pmii.id"},
	{5, "Budi Santoso", "budi.santoso@pmii.id"},
	{6, "Dewi Lestari", "dewi.lestari@pmii.id"},
	{7, "Rizki Pratama", "rizki.pratama@pmii.id"},
	{8, "Nurul Hidayah", "nurul.hidayah@pmii.id"},
	{9, "Hendra Wijaya", "hendra.wijaya@pmii.id"},
	{10, "Fatimah Zahra", "fatimah.zahra@pmii.id"},
	{11, "Agus Setiawan", "agus.setiawan@pmii.id"},
	{12, "Putri Amelia", "putri.amelia@pmii.id"},
}

// Activity log templates
var activityTemplates = []struct {
	ActionType  domain.ActivityActionType
	Module      domain.ActivityModuleType
	Description string
	AdminOnly   bool // true if only admin can perform this action
}{
	{domain.ActionCreate, domain.ModulePost, "Membuat artikel baru", false},
	{domain.ActionUpdate, domain.ModulePost, "Mengupdate artikel", false},
	{domain.ActionDelete, domain.ModulePost, "Menghapus artikel", false},
	{domain.ActionCreate, domain.ModuleCategory, "Membuat kategori baru", false},
	{domain.ActionUpdate, domain.ModuleCategory, "Mengupdate kategori", false},
	{domain.ActionDelete, domain.ModuleCategory, "Menghapus kategori", false},
	{domain.ActionCreate, domain.ModuleTags, "Membuat tag baru", false},
	{domain.ActionUpdate, domain.ModuleTags, "Mengupdate tag", false},
	{domain.ActionDelete, domain.ModuleTags, "Menghapus tag", false},
	{domain.ActionLogin, domain.ModuleAuth, "Login ke sistem", false},
	{domain.ActionCreate, domain.ModuleUser, "Membuat user baru", true},
	{domain.ActionUpdate, domain.ModuleUser, "Mengupdate data user", true},
	{domain.ActionDelete, domain.ModuleUser, "Menghapus user", true},
}

var ipAddresses = []string{
	"192.168.1.100",
	"192.168.1.101",
	"192.168.1.102",
	"10.0.0.50",
	"10.0.0.51",
	"10.0.1.100",
	"172.16.0.10",
	"172.16.0.11",
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/119.0.0.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15",
	"Mozilla/5.0 (Android 14; Mobile) AppleWebKit/537.36 Chrome/120.0.0.0",
}

func main() {
	log.Println("🌱 Starting Development Seeder...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}

	db, err := database.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Seed development users
	if err := seedDevUsers(db); err != nil {
		log.Fatalf("❌ Failed to seed users: %v", err)
	}

	// Seed development activity logs
	if err := seedDevActivityLogs(db); err != nil {
		log.Fatalf("❌ Failed to seed activity logs: %v", err)
	}

	log.Println("✅ Development seeding completed successfully!")
}

func seedDevUsers(db *gorm.DB) error {
	log.Println("🌱 Seeding development users...")

	// Hash password once for all users
	hashedPassword, err := utils.HashPassword("author123")
	if err != nil {
		return err
	}

	seededCount := 0
	skippedCount := 0

	for _, userData := range devUsers {
		// Check if user already exists by ID or email (including soft-deleted)
		var count int64
		db.Unscoped().Model(&domain.User{}).Where("id = ? OR email = ?", userData.ID, userData.Email).Count(&count)

		if count > 0 {
			skippedCount++
			continue
		}

		user := domain.User{
			ID:           userData.ID,
			Role:         2, // Author
			FullName:     userData.FullName,
			Email:        userData.Email,
			PasswordHash: hashedPassword,
			IsActive:     true,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("⚠️  Failed to create user %s: %v", userData.Email, err)
			continue
		}

		seededCount++
	}

	log.Printf("✅ Users seeded: %d created, %d skipped (already exist)", seededCount, skippedCount)
	return nil
}

func seedDevActivityLogs(db *gorm.DB) error {
	log.Println("🌱 Seeding development activity logs...")

	// Check if activity logs already exist
	var count int64
	db.Model(&domain.ActivityLog{}).Count(&count)

	if count >= 20 {
		log.Printf("ℹ️  Activity logs already exist (%d logs), skipping...", count)
		return nil
	}

	// Seed random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// All user IDs (admin + default user + dev users)
	allUserIDs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	adminUserID := 1

	seededCount := 0
	targetCount := 20

	for i := 0; i < targetCount; i++ {
		// Pick a random activity template
		template := activityTemplates[r.Intn(len(activityTemplates))]

		// Determine user ID based on activity type
		var userID int
		if template.AdminOnly {
			userID = adminUserID
		} else {
			userID = allUserIDs[r.Intn(len(allUserIDs))]
		}

		// Random IP and user agent
		ipAddress := ipAddresses[r.Intn(len(ipAddresses))]
		userAgent := userAgents[r.Intn(len(userAgents))]

		// Random target ID (1-100)
		targetID := r.Intn(100) + 1

		activityLog := domain.ActivityLog{
			UserID:      userID,
			ActionType:  template.ActionType,
			Module:      template.Module,
			Description: &template.Description,
			TargetID:    &targetID,
			IPAddress:   &ipAddress,
			UserAgent:   &userAgent,
			CreatedAt:   time.Now().Add(-time.Duration(r.Intn(720)) * time.Hour), // Random time in last 30 days
		}

		if err := db.Create(&activityLog).Error; err != nil {
			log.Printf("⚠️  Failed to create activity log: %v", err)
			continue
		}

		seededCount++
	}

	log.Printf("✅ Activity logs seeded: %d created", seededCount)
	return nil
}
