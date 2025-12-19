package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"gorm.io/gorm"
)

type HomeRepository interface {
	GetHeroSection() ([]responses.HeroSectionResponse, error)
	GetLatestNewsSection() ([]responses.LatestNewsSectionResponse, error)
	GetAboutUsSection() (*responses.AboutUsSectionResponse, error)
	GetWhySection() (*responses.WhySectionResponse, error)
	GetFaqSection() (*responses.FaqSectionResponse, error)
	GetCtaSection() (*responses.CtaSectionResponse, error)
}

type homeRepository struct {
	db *gorm.DB
}

func NewHomeRepository(db *gorm.DB) HomeRepository {
	return &homeRepository{db: db}
}

func (r *homeRepository) GetHeroSection() ([]responses.HeroSectionResponse, error) {
	var posts []responses.HeroSectionResponse

	// Query only published posts for hero section
	// GORM returns empty slice [] if no data found (not an error)
	err := r.db.Table("posts").
		Select(`
			posts.id,
			posts.title,
			posts.featured_image,
			posts.created_at,
			posts.updated_at,
			COUNT(post_views.id) AS total_views
		`).
		Where("posts.status = ?", "published").
		Joins("LEFT JOIN post_views ON post_views.post_id = posts.id").
		Group("posts.id").
		Order("total_views DESC").
		Limit(5).
		Find(&posts).Error
	if err != nil {
		// Only database errors (connection issues, query errors, etc.)
		return nil, err
	}

	// Returns posts slice (can be empty []) with no error
	return posts, nil
}

func (r *homeRepository) GetLatestNewsSection() ([]responses.LatestNewsSectionResponse, error) {
	var news []responses.LatestNewsSectionResponse

	// Query only published posts for latest news section
	// GORM returns empty slice [] if no data found (not an error)
	err := r.db.Table("posts").
		Select(`
			posts.id,
			posts.title,
			posts.featured_image,
			posts.created_at,
			posts.updated_at,
			COUNT(post_views.id) AS total_views
		`).
		Where("posts.status = ?", "published").
		Joins("LEFT JOIN post_views ON post_views.post_id = posts.id").
		Joins("INNER JOIN categories ON categories.id = posts.category_id").
		Where("categories.slug = ?", "news").
		Group("posts.id").
		Order("posts.created_at DESC").
		Limit(5).
		Find(&news).Error
	if err != nil {
		// Only database errors (connection issues, query errors, etc.)
		return nil, err
	}

	// Returns news slice (can be empty []) with no error
	return news, nil
}

// todo: this is temporarily harcoded
func (r *homeRepository) GetAboutUsSection() (*responses.AboutUsSectionResponse, error) {
	return &responses.AboutUsSectionResponse{
		Title:       "Sekilas tentang sejarah, nilai, dan arah gerakan PMII.",
		Subtitle:    "Tentang PMII",
		Description: "Organisasi mahasiswa berbasis nilai keislaman dan keindonesiaan yang telah bergerak sejak 1960 untuk mencetak kader bangsa berkarakter dan berwawasan luas.",
		ImageURI:    "about-image.jpg",
	}, nil
}

// todo: this is temporarily hardcoded
func (r *homeRepository) GetWhySection() (*responses.WhySectionResponse, error) {
	var whyItem []map[string]string

	whyItem = append(whyItem, map[string]string{
		"title":       "Perkaderan Benjenjang",
		"description": "Pembentukan karakter dan kepemimpinan.",
		"iconURI":     "icon1.png",
	})

	whyItem = append(whyItem, map[string]string{
		"title":       "Komunitas Nasional",
		"description": "Hadir di berbagai kampus di Indonesia",
		"iconURI":     "icon2.png",
	})

	whyItem = append(whyItem, map[string]string{
		"title":       "Kegiatan Intelektual & Sosial",
		"description": "Pelatihan, kajian, dan aksi sosial",
		"iconURI":     "icon3.png",
	})

	whyItem = append(whyItem, map[string]string{
		"title":       "Akses Jaringan Alumni",
		"description": "Dukungan karier dan mentorship",
		"iconURI":     "icon4.png",
	})

	return &responses.WhySectionResponse{
		Title:       "Nilai dan keunggulan yang membuat PMII menjadi ruang tumbuh bagi mahasiswa",
		Subtitle:    "Mengapa PMII?",
		Description: nil,
		Data:        whyItem,
	}, nil
}

// todo: this is temporarily hardcoded
func (r *homeRepository) GetFaqSection() (*responses.FaqSectionResponse, error) {
	var faqItem []map[string]string

	faqItem = append(faqItem, map[string]string{
		"title":       "Apa itu PMII?",
		"description": "PMII adalah organisasi mahasiswa berbasis nilai keislaman dan keindonesiaan yang telah bergerak sejak 1960 untuk mencetak kader bangsa berkarakter dan berwawasan luas.",
	})

	faqItem = append(faqItem, map[string]string{
		"title":       "Apa tujuan PMII?",
		"description": "Tujuan PMII adalah untuk mencetak kader bangsa berkarakter dan berwawasan luas.",
	})

	faqItem = append(faqItem, map[string]string{
		"title":       "Bagaimana saya bisa menjadi anggota PMII?",
		"description": "Untuk menjadi anggota PMII, Anda perlu mengikuti proses pendaftaran yang dapat Anda lakukan di website resmi PMII.",
	})

	faqItem = append(faqItem, map[string]string{
		"title":       "Apakah PMII ada di semua kampus?",
		"description": "PMII ada di berbagai kampus di Indonesia.",
	})

	faqItem = append(faqItem, map[string]string{
		"title":       "Apakah menjadi anggota dikenakan biaya?",
		"description": "Tidak, menjadi anggota PMII tidak dikenakan biaya.",
	})

	return &responses.FaqSectionResponse{
		Title:       "Pertanyaan sering diajukan",
		Subtitle:    "FAQ",
		Description: nil,
		Data:        faqItem,
	}, nil
}

// todo: this is temporarily hardcoded
func (r *homeRepository) GetCtaSection() (*responses.CtaSectionResponse, error) {
	return &responses.CtaSectionResponse{
		Title:    "Ayo Bergabung",
		Subtitle: "Bergabung dengan PMII",
	}, nil
}
