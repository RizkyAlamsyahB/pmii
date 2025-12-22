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
	type heroPostResult struct {
		ID            int
		Title         string
		Slug          string
		Excerpt       *string
		FeaturedImage *string
		PublishedAt   *string
		CategoryID    int
		CategoryName  string
		UserID        int
		TotalViews    int64
	}

	var results []heroPostResult

	// Query only published posts for hero section with category info
	err := r.db.Table("posts").
		Select(`
			posts.id,
			posts.title,
			posts.slug,
			posts.excerpt,
			posts.featured_image,
			posts.published_at,
			posts.category_id,
			categories.name AS category_name,
			posts.user_id,
			COUNT(post_views.id) AS total_views
		`).
		Joins("LEFT JOIN categories ON categories.id = posts.category_id").
		Joins("LEFT JOIN post_views ON post_views.post_id = posts.id").
		Where("posts.status = ?", "published").
		Where("posts.deleted_at IS NULL").
		Group("posts.id, categories.name").
		Order("total_views DESC").
		Limit(5).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// Build the response with tags for each post
	var heroSections []responses.HeroSectionResponse
	for _, result := range results {
		// Get tags for this post
		var tagNames []string
		r.db.Table("tags").
			Select("tags.name").
			Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
			Where("post_tags.post_id = ?", result.ID).
			Pluck("name", &tagNames)

		excerpt := ""
		if result.Excerpt != nil {
			excerpt = *result.Excerpt
		}

		imageURL := ""
		if result.FeaturedImage != nil {
			imageURL = *result.FeaturedImage
		}

		publishedAt := ""
		if result.PublishedAt != nil {
			publishedAt = *result.PublishedAt
		}

		heroSection := responses.HeroSectionResponse{
			ID:          int64(result.ID),
			Title:       result.Title,
			Slug:        result.Slug,
			Excerpt:     excerpt,
			ImageURL:    imageURL,
			PublishedAt: publishedAt,
			Category: responses.HeroCategoryResponse{
				ID:   int64(result.CategoryID),
				Name: result.CategoryName,
			},
			AuthorID: int64(result.UserID),
			Tags:     tagNames,
		}

		heroSections = append(heroSections, heroSection)
	}

	// Return empty slice if no data
	if heroSections == nil {
		heroSections = []responses.HeroSectionResponse{}
	}

	return heroSections, nil
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

func (r *homeRepository) GetAboutUsSection() (*responses.AboutUsSectionResponse, error) {
	return &responses.AboutUsSectionResponse{
		Title:       "Sekilas tentang sejarah, nilai, dan arah gerakan PMII.",
		Subtitle:    "Tentang PMII",
		Description: "Organisasi mahasiswa berbasis nilai keislaman dan keindonesiaan yang telah bergerak sejak 1960 untuk mencetak kader bangsa berkarakter dan berwawasan luas.",
		ImageURI:    "about-image.jpg",
	}, nil
}

func (r *homeRepository) GetWhySection() (*responses.WhySectionResponse, error) {
	whyItems := []responses.WhyItem{
		{
			Title:       "Perkaderan Benjenjang",
			Description: "Pembentukan karakter dan kepemimpinan.",
			IconURI:     "icon1.png",
		},
		{
			Title:       "Komunitas Nasional",
			Description: "Hadir di berbagai kampus di Indonesia",
			IconURI:     "icon2.png",
		},
		{
			Title:       "Kegiatan Intelektual & Sosial",
			Description: "Pelatihan, kajian, dan aksi sosial",
			IconURI:     "icon3.png",
		},
		{
			Title:       "Akses Jaringan Alumni",
			Description: "Dukungan karier dan mentorship",
			IconURI:     "icon4.png",
		},
	}

	return &responses.WhySectionResponse{
		Title:       "Nilai dan keunggulan yang membuat PMII menjadi ruang tumbuh bagi mahasiswa",
		Subtitle:    "Mengapa PMII?",
		Description: nil,
		Data:        whyItems,
	}, nil
}

func (r *homeRepository) GetFaqSection() (*responses.FaqSectionResponse, error) {
	faqItems := []responses.FaqItem{
		{
			Question: "Apa itu PMII?",
			Answer:   "PMII adalah organisasi mahasiswa berbasis nilai keislaman dan keindonesiaan yang telah bergerak sejak 1960 untuk mencetak kader bangsa berkarakter dan berwawasan luas.",
		},
		{
			Question: "Siapa yang bisa bergabung?",
			Answer:   "Semua mahasiswa yang memiliki nilai keislaman dan keindonesiaan dapat bergabung dengan PMII.",
		},
		{
			Question: "Bagaimana proses kaderisasi?",
			Answer:   "Proses kaderisasi melibatkan beberapa tahap, termasuk pendaftaran, wawancara, dan pengujian karakter.",
		},
		{
			Question: "Apakah PMII ada di semua kampus?",
			Answer:   "PMII ada di berbagai kampus di Indonesia.",
		},
		{
			Question: "Apakah menjadi anggota dikenakan biaya?",
			Answer:   "Tidak, menjadi anggota PMII tidak dikenakan biaya.",
		},
	}

	return &responses.FaqSectionResponse{
		Title:       "Jawaban atas pertanyaan umum dari calon anggota dan pengunjung.",
		Subtitle:    "FAQ",
		Description: nil,
		Data:        faqItems,
	}, nil
}

func (r *homeRepository) GetCtaSection() (*responses.CtaSectionResponse, error) {
	return &responses.CtaSectionResponse{
		Title:    "Ambil langkah pertama untuk menjadi bagian dari gerakan perubahan.",
		Subtitle: "Siap Bergabung dengan PMII?",
	}, nil
}
