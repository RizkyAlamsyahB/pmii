package seeds

import (
	"fmt"
	"log"

	"github.com/garuda-labs-1/pmii-be/pkg/cloudinary"
	"gorm.io/gorm"
)

// Seeder orchestrates all content seeding
type Seeder struct {
	db        *gorm.DB
	uploader  *cloudinary.Service
	seedsPath string
}

// NewSeeder creates a new Seeder instance
func NewSeeder(db *gorm.DB, uploader *cloudinary.Service, seedsPath string) *Seeder {
	return &Seeder{
		db:        db,
		uploader:  uploader,
		seedsPath: seedsPath,
	}
}

// SeedAll runs all content seeders (only runs once)
// Uses site_settings table as marker - if exists, skip all seeding
func (s *Seeder) SeedAll() error {
	// Check if already seeded by looking at site_settings (singleton table)
	var count int64
	s.db.Table("site_settings").Count(&count)
	if count > 0 {
		log.Println("ℹ️  Content already seeded, skipping...")
		return nil
	}

	log.Println("🌱 Starting content seeding...")

	seeders := []struct {
		name string
		fn   func() error
	}{
		{"About", s.SeedAbout},
		{"Site Settings", s.SeedSiteSettings},
		{"Contact", s.SeedContact},
		{"Members", s.SeedMembers},
		{"Testimonials", s.SeedTestimonials},
		{"Documents", s.SeedDocuments},
	}

	var failedSeeders []string
	for _, seeder := range seeders {
		if err := seeder.fn(); err != nil {
			log.Printf("⚠️ Warning: Failed to seed %s: %v", seeder.name, err)
			failedSeeders = append(failedSeeders, seeder.name)
		}
	}

	if len(failedSeeders) > 0 {
		log.Printf("⚠️ Content seeding completed with %d failures: %v", len(failedSeeders), failedSeeders)
		return fmt.Errorf("failed to seed: %v", failedSeeders)
	}

	log.Println("✅ Content seeding completed successfully!")
	return nil
}

// SeedSelected runs specific seeders by name
func (s *Seeder) SeedSelected(names ...string) error {
	log.Printf("🌱 Starting selected seeding: %v", names)

	seederMap := map[string]func() error{
		"about":        s.SeedAbout,
		"settings":     s.SeedSiteSettings,
		"contact":      s.SeedContact,
		"members":      s.SeedMembers,
		"testimonials": s.SeedTestimonials,
		"documents":    s.SeedDocuments,
	}

	for _, name := range names {
		if fn, ok := seederMap[name]; ok {
			if err := fn(); err != nil {
				log.Printf("⚠️ Warning: Failed to seed %s: %v", name, err)
			}
		} else {
			log.Printf("⚠️ Unknown seeder: %s", name)
		}
	}

	log.Println("✅ Selected seeding completed!")
	return nil
}
