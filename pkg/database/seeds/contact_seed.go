package seeds

import (
	"log"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// ContactSeedData contains contact seeding data
type ContactSeedData struct {
	Email         string
	Phone         string
	Address       string
	GoogleMapsURL string
}

// SeedContact seeds contact data
func (s *Seeder) SeedContact() error {
	logSeederStart("Contact")

	var count int64
	s.db.Model(&domain.Contact{}).Count(&count)
	if count > 0 {
		logSeederSkip("Contact")
		return nil
	}

	data := getContactData()

	contact := domain.Contact{
		Email:         &data.Email,
		Phone:         &data.Phone,
		Address:       &data.Address,
		GoogleMapsURL: &data.GoogleMapsURL,
	}

	if err := s.db.Create(&contact).Error; err != nil {
		log.Printf("❌ Failed to create contact: %v", err)
		return err
	}

	log.Println("✅ Contact seeded successfully")
	return nil
}

// getContactData returns contact seed data
func getContactData() ContactSeedData {
	return ContactSeedData{
		Email:         "bidangmediapbpmii@gmail.com",
		Phone:         "62213920047",
		Address:       "Jl. Salemba Tengah No.57, RT.10/RW.8, Paseban, Kec. Senen, Kota Jakarta Pusat, DKI Jakarta 10440",
		GoogleMapsURL: "https://www.google.com/maps?ll=-6.193005,106.855055&z=16&t=m&hl=en&gl=ID&mapclient=embed&cid=12235338042217647956",
	}
}
