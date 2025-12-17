package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"gorm.io/gorm"
)

type HomeRepository interface {
	GetHeroSection() ([]responses.HeroSectionResponse, error)
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
