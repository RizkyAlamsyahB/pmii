package seeds

import (
	"log"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// DocumentSeedData contains document seeding data
type DocumentSeedData struct {
	Name     string
	FileType domain.DocumentType
	DocFile  string
}

// SeedDocuments seeds document data
func (s *Seeder) SeedDocuments() error {
	logSeederStart("Documents")

	// Clear existing documents and reset sequence for clean re-seeding
	s.db.Exec("TRUNCATE TABLE documents RESTART IDENTITY CASCADE")

	documents := getDocumentsData()
	successCount := 0

	for i, d := range documents {
		// Upload document - files are in subfolders by type (produk_hukum, lagu_organisasi, logo_organisasi)
		docPath := getFilePath(s.seedsPath, "documents", string(d.FileType), d.DocFile)
		docURL, err := uploadFile(s.uploader, docPath, d.FileType.GetCloudinaryFolder())
		if err != nil {
			log.Printf("⚠️ Warning: Failed to upload document %s: %v", d.Name, err)
			continue
		}

		document := domain.Document{
			Name:     d.Name,
			FileType: d.FileType,
			FileURI:  docURL,
		}

		if err := s.db.Create(&document).Error; err != nil {
			log.Printf("⚠️ Warning: Failed to create document %s: %v", d.Name, err)
			continue
		}

		successCount++
		logSeederProgress(i+1, len(documents), d.Name)
	}

	logSeederResult("Documents", successCount, len(documents))
	return nil
}

// getDocumentsData returns all document seed data
func getDocumentsData() []DocumentSeedData {
	return []DocumentSeedData{
		// Produk Hukum
		{Name: "AD-ART PMII KONGRES 2021 BALIKPAPAN", FileType: domain.DocumentTypeProdukHukum, DocFile: "1766022783.pdf"},
		{Name: "AKTA & SK Kemenkumham PB PMII 2022", FileType: domain.DocumentTypeProdukHukum, DocFile: "1766022790.pdf"},
		{Name: "HASIL MUSPIMNAS PMII TULUNGAGUNG 2022", FileType: domain.DocumentTypeProdukHukum, DocFile: "1766022813.pdf"},

		// Lagu Organisasi
		{Name: "Mars PMII (Bahasa Arab)", FileType: domain.DocumentTypeLaguOrganisasi, DocFile: "1766022845.wma"},
		{Name: "Darah Juang PMII", FileType: domain.DocumentTypeLaguOrganisasi, DocFile: "1766022855.mp3"},
		{Name: "Hymne PMII", FileType: domain.DocumentTypeLaguOrganisasi, DocFile: "1766022884.mp3"},
		{Name: "Mars PMII (Bahasa Indonesia)", FileType: domain.DocumentTypeLaguOrganisasi, DocFile: "1766022891.wma"},
		{Name: "Mars PMII (Bahasa Inggris)", FileType: domain.DocumentTypeLaguOrganisasi, DocFile: "1766022901.wma"},
		{Name: "Mars PMII (Orchestra)", FileType: domain.DocumentTypeLaguOrganisasi, DocFile: "1766024088.wma"},

		// Logo Organisasi
		{Name: "Logo PMII Resmi", FileType: domain.DocumentTypeLogoOrganisasi, DocFile: "1766022913.png"},
	}
}
