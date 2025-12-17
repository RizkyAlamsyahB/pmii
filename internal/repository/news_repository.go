package repository

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

type NewsRepository interface {
	GetPublishedNews(offset, limit int, search string) ([]domain.Post, int64, error)
	GetNewsBySlug(slug string) (domain.Post, error)
}

type newsRepository struct{}

func NewNewsRepository() NewsRepository {
	return &newsRepository{}
}

func (r *newsRepository) GetPublishedNews(offset, limit int, search string) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	db := config.DB.Model(&domain.Post{}).
		Preload("Category").
		Preload("Tags").
		Where("status = ?", "published") // Filter hanya yang sudah publish

	if search != "" {
		db = db.Where("title ILIKE ?", "%"+search+"%")
	}

	db.Count(&total)
	err := db.Limit(limit).Offset(offset).Order("published_at DESC").Find(&posts).Error

	return posts, total, err
}

func (r *newsRepository) GetNewsBySlug(slug string) (domain.Post, error) {
	var post domain.Post
	err := config.DB.Preload("Category").
		Preload("Tags").
		Where("slug = ? AND status = ?", slug, "published").
		First(&post).Error
	return post, err
}
