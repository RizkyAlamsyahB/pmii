package seeds

import (
	"log"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// AboutSeedData contains about seeding data
type AboutSeedData struct {
	Title    string
	History  string
	Vision   string
	Mission  string
	VideoURL string
}

// SeedAbout seeds about data
func (s *Seeder) SeedAbout() error {
	logSeederStart("About")

	var count int64
	s.db.Model(&domain.About{}).Count(&count)
	if count > 0 {
		logSeederSkip("About")
		return nil
	}

	aboutData := getAboutData()

	about := domain.About{
		Title:    &aboutData.Title,
		History:  &aboutData.History,
		Vision:   &aboutData.Vision,
		Mission:  &aboutData.Mission,
		VideoURL: &aboutData.VideoURL,
	}

	if err := s.db.Create(&about).Error; err != nil {
		log.Printf("❌ Failed to create about: %v", err)
		return err
	}

	log.Println("✅ About seeded successfully")
	return nil
}

// getAboutData returns about seed data
func getAboutData() AboutSeedData {
	return AboutSeedData{
		Title:    "PERGERAKAN MAHASISWA ISLAM INDONESIA",
		History:  "PMII merupakan organisasi gerakan dan kaderisasi yang berlandaskan islam ahlussunah waljamaah. Berdiri sejak tanggal 17 April 1960 di Surabaya dan hingga lebih dari setengah abad kini PMII terus eksis untuk memberikan kontribusi bagi kemajuan bangsa dan negara.",
		Vision:   "Terwujudnya kader PMII yang berilmu, berakhlak...",
		Mission:  "Menguatkan Profesionalitas Organisasi Menuju Era Baru PMII",
		VideoURL: "https://youtu.be/zFN7dJa4niw",
	}
}
