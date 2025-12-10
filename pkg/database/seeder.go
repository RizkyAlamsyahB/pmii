package database

import (
	"log"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
	"gorm.io/gorm"
)

// SeedDefaultUsers membuat default admin & user jika belum ada
// Auto-run saat aplikasi start untuk memastikan ada user default
func SeedDefaultUsers(db *gorm.DB) error {
	log.Println("🌱 Checking default users...")

	// 1. Seed Admin User
	var adminCount int64
	db.Model(&domain.User{}).Where("email = ?", "admin@pmii.id").Count(&adminCount)

	if adminCount == 0 {
		// Hash password admin123
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			return err
		}

		admin := domain.User{
			Role:         1, // Admin
			FullName:     "Administrator",
			Email:        "admin@pmii.id",
			PasswordHash: hashedPassword,
			PhotoURI:     nil,
			IsActive:     true,
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		log.Println("✅ Default admin user created: admin@pmii.id / admin123")
	} else {
		log.Println("ℹ️  Admin user already exists")
	}

	// 2. Seed Regular User (Author)
	var userCount int64
	db.Model(&domain.User{}).Where("email = ?", "user@pmii.id").Count(&userCount)

	if userCount == 0 {
		// Hash password user123
		hashedPassword, err := utils.HashPassword("user123")
		if err != nil {
			return err
		}

		user := domain.User{
			Role:         2, // Author
			FullName:     "Regular User",
			Email:        "user@pmii.id",
			PasswordHash: hashedPassword,
			PhotoURI:     nil,
			IsActive:     true,
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}
		log.Println("✅ Default user created: user@pmii.id / user123")
	} else {
		log.Println("ℹ️  Regular user already exists")
	}

	log.Println("✅ User seeding completed")
	return nil
}
