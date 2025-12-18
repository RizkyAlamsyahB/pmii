package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"gorm.io/gorm"
)

type HomeRepository interface {
	GetHeroSection() ([]responses.HeroSectionResponse, error)
	GetLatestNewsSection() ([]responses.LatestNewsSectionResponse, error)
	GetAboutUsSection() (*responses.AboutUsSectionResponse, error)
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
		Order("RANDOM()").
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
