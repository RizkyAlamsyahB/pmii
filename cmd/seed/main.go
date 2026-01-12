package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/pkg/cloudinary"
	"github.com/garuda-labs-1/pmii-be/pkg/database"
	"github.com/garuda-labs-1/pmii-be/pkg/database/seeds"
)

func main() {
	// Parse command line flags
	seedAll := flag.Bool("all", false, "Seed all content data")
	seedAbout := flag.Bool("about", false, "Seed about data")
	seedSettings := flag.Bool("settings", false, "Seed site settings")
	seedContact := flag.Bool("contact", false, "Seed contact data")
	seedMembers := flag.Bool("members", false, "Seed members data")
	seedTestimonials := flag.Bool("testimonials", false, "Seed testimonials data")
	seedDocuments := flag.Bool("documents", false, "Seed documents data")
	flag.Parse()

	// Check if any flag is set
	if !*seedAll && !*seedAbout && !*seedSettings && !*seedContact && !*seedMembers && !*seedTestimonials && !*seedDocuments {
		fmt.Println("Usage: go run cmd/seed/main.go [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -all          Seed all content data")
		fmt.Println("  -about        Seed about data")
		fmt.Println("  -settings     Seed site settings")
		fmt.Println("  -contact      Seed contact data")
		fmt.Println("  -members      Seed members data")
		fmt.Println("  -testimonials Seed testimonials data")
		fmt.Println("  -documents    Seed documents data")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run cmd/seed/main.go -all")
		fmt.Println("  go run cmd/seed/main.go -members -testimonials")
		os.Exit(0)
	}

	log.Println("🌱 Starting manual seeder...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}

	db, err := database.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected")

	// Initialize Cloudinary
	cloudinaryService, err := cloudinary.NewService(cfg.Cloudinary.URL)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	log.Println("✅ Cloudinary initialized")

	// Create seeder
	seeder := seeds.NewSeeder(db, cloudinaryService, "./seeds")

	// Run seeders based on flags
	if *seedAll {
		log.Println("📦 Running all seeders...")
		if err := seeder.SeedAll(); err != nil {
			log.Printf("⚠️ Seeding completed with errors: %v", err)
		}
	} else {
		// Run selected seeders
		selectedSeeders := []string{}

		if *seedAbout {
			selectedSeeders = append(selectedSeeders, "about")
		}
		if *seedSettings {
			selectedSeeders = append(selectedSeeders, "settings")
		}
		if *seedContact {
			selectedSeeders = append(selectedSeeders, "contact")
		}
		if *seedMembers {
			selectedSeeders = append(selectedSeeders, "members")
		}
		if *seedTestimonials {
			selectedSeeders = append(selectedSeeders, "testimonials")
		}
		if *seedDocuments {
			selectedSeeders = append(selectedSeeders, "documents")
		}

		if len(selectedSeeders) > 0 {
			log.Printf("📦 Running selected seeders: %v", selectedSeeders)
			if err := seeder.SeedSelected(selectedSeeders...); err != nil {
				log.Printf("⚠️ Seeding completed with errors: %v", err)
			}
		}
	}

	log.Println("✅ Seeding completed!")
}
