package seeds

import (
	"log"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// TestimonialSeedData contains testimonial seeding data
type TestimonialSeedData struct {
	Name         string
	Organization string
	Position     string
	Content      string
	ImageFile    string
}

// SeedTestimonials seeds testimonial data
func (s *Seeder) SeedTestimonials() error {
	logSeederStart("Testimonials")

	// Clear existing testimonials and reset sequence for clean re-seeding
	s.db.Exec("TRUNCATE TABLE testimonials RESTART IDENTITY CASCADE")

	testimonials := getTestimonialsData()
	successCount := 0

	for i, t := range testimonials {
		// Upload image
		imagePath := getFilePath(s.seedsPath, "testimonials", t.ImageFile)
		imageURL, err := uploadFile(s.uploader, imagePath, "testimonials")
		if err != nil {
			log.Printf("⚠️ Warning: Failed to upload image for %s: %v", t.Name, err)
			continue
		}

		testimonial := domain.Testimonial{
			Name:         t.Name,
			Organization: &t.Organization,
			Position:     &t.Position,
			Content:      t.Content,
			PhotoURI:     &imageURL,
			IsActive:     true,
		}

		if err := s.db.Create(&testimonial).Error; err != nil {
			log.Printf("⚠️ Warning: Failed to create testimonial %s: %v", t.Name, err)
			continue
		}

		successCount++
		logSeederProgress(i+1, len(testimonials), t.Name)
	}

	logSeederResult("Testimonials", successCount, len(testimonials))
	return nil
}

// getTestimonialsData returns all testimonial seed data
func getTestimonialsData() []TestimonialSeedData {
	return []TestimonialSeedData{
		{
			Name:         "Ahmad Fauzi",
			Organization: "Universitas Indonesia",
			Position:     "Ketua BEM UI",
			Content:      "PMII telah membentuk karakter kepemimpinan saya. Melalui kaderisasi yang intensif, saya belajar bagaimana menjadi pemimpin yang amanah dan berintegritas. Organisasi ini mengajarkan nilai-nilai Ahlussunnah wal Jamaah yang menjadi pegangan hidup saya.",
			ImageFile:    "1766113265.jpg",
		},
		{
			Name:         "Siti Nurhaliza",
			Organization: "UIN Syarif Hidayatullah Jakarta",
			Position:     "Pengurus PC PMII Jakarta",
			Content:      "Bergabung dengan PMII adalah keputusan terbaik dalam hidup saya. Di sini saya menemukan keluarga baru yang peduli terhadap sesama dan berkomitmen untuk perubahan sosial. Program-program pemberdayaan masyarakat yang kami jalankan sangat berdampak.",
			ImageFile:    "1766113267.jpg",
		},
		{
			Name:         "Muhammad Rizki",
			Organization: "Universitas Brawijaya",
			Position:     "Alumni PMII Malang",
			Content:      "PMII mengajarkan saya untuk berpikir kritis dan tidak mudah terpengaruh propaganda. Kajian-kajian keislaman yang moderat membuka wawasan saya tentang Islam yang rahmatan lil alamin. Terima kasih PMII!",
			ImageFile:    "1766113269.jpg",
		},
		{
			Name:         "Dewi Kartika",
			Organization: "Universitas Gadjah Mada",
			Position:     "Koordinator Bidang Perempuan",
			Content:      "Sebagai perempuan, saya merasa sangat diberdayakan di PMII. Organisasi ini memberikan ruang yang luas bagi perempuan untuk berkontribusi dan memimpin. Nilai kesetaraan gender yang dijunjung tinggi sangat menginspirasi.",
			ImageFile:    "1766113270.jpg",
		},
	}
}
