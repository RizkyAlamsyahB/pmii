package seeds

import (
	"log"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// SiteSettingsSeedData contains site settings seeding data
type SiteSettingsSeedData struct {
	SiteName        string
	SiteTitle       string
	SiteDescription string
	LogoHeaderFile  string
	LogoBigFile     string
	FaviconFile     string
	FacebookURL     string
	TwitterURL      string
	InstagramURL    string
	YoutubeURL      string
	LinkedinURL     string
	GithubURL       string
}

// SeedSiteSettings seeds site settings data
func (s *Seeder) SeedSiteSettings() error {
	logSeederStart("Site Settings")

	var count int64
	s.db.Model(&domain.SiteSetting{}).Count(&count)
	if count > 0 {
		logSeederSkip("Site Settings")
		return nil
	}

	data := getSiteSettingsData()

	// Upload logo header
	logoHeaderPath := getFilePath(s.seedsPath, "site_settings", data.LogoHeaderFile)
	logoHeaderURL, err := uploadFile(s.uploader, logoHeaderPath, "settings")
	if err != nil {
		log.Printf("⚠️ Warning: Failed to upload logo header: %v", err)
		return err
	}

	// Upload logo big
	logoBigPath := getFilePath(s.seedsPath, "site_settings", data.LogoBigFile)
	logoBigURL, err := uploadFile(s.uploader, logoBigPath, "settings")
	if err != nil {
		log.Printf("⚠️ Warning: Failed to upload logo big: %v", err)
		return err
	}

	// Upload favicon
	faviconPath := getFilePath(s.seedsPath, "site_settings", data.FaviconFile)
	faviconURL, err := uploadFile(s.uploader, faviconPath, "settings")
	if err != nil {
		log.Printf("⚠️ Warning: Failed to upload favicon: %v", err)
		return err
	}

	settings := domain.SiteSetting{
		SiteName:        &data.SiteName,
		SiteTitle:       &data.SiteTitle,
		SiteDescription: &data.SiteDescription,
		LogoHeader:      &logoHeaderURL,
		LogoBig:         &logoBigURL,
		Favicon:         &faviconURL,
		FacebookURL:     &data.FacebookURL,
		TwitterURL:      &data.TwitterURL,
		InstagramURL:    &data.InstagramURL,
		YoutubeURL:      &data.YoutubeURL,
		LinkedinURL:     &data.LinkedinURL,
		GithubURL:       &data.GithubURL,
	}

	if err := s.db.Create(&settings).Error; err != nil {
		log.Printf("❌ Failed to create site settings: %v", err)
		return err
	}

	log.Println("✅ Site Settings seeded successfully")
	return nil
}

// getSiteSettingsData returns site settings seed data
func getSiteSettingsData() SiteSettingsSeedData {
	return SiteSettingsSeedData{
		SiteName:        "PB PMII",
		SiteTitle:       "PB PMII | Pengurus Besar Pergerakan Mahasiswa Islam Indonesia",
		SiteDescription: "Website Resmi Pengurus Besar Pergerakan Mahasiswa Islam Indonesia",
		LogoHeaderFile:  "1766972876.webp",
		LogoBigFile:     "1766972752.jpg",
		FaviconFile:     "1766972620.png",
		FacebookURL:     "https://www.facebook.com/PMIIOFFICIAL17",
		TwitterURL:      "https://twitter.com/pmiiofficial",
		InstagramURL:    "https://www.instagram.com/pmiiofficial",
		YoutubeURL:      "https://www.youtube.com/c/PMIIOFFICIAL",
		LinkedinURL:     "https://id.linkedin.com/in/ircham-ali",
		GithubURL:       "",
	}
}
