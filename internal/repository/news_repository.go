package repository

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

type NewsRepository interface {
	GetPublishedNews(offset, limit int, search string) ([]domain.Post, int64, error)
	GetNewsBySlug(slug string) (domain.Post, error)
	//metod untuk mendapatkan berita berdasarkan kategori
	GetNewsByCategorySlug(categorySlug string, offset, limit int) ([]domain.Post, int64, error)
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) GetPublishedNews(offset, limit int, search string) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// Tambahkan Select subquery agar angka views muncul di list /v1/news
	query := r.db.Model(&domain.Post{}).
		Select("posts.*, (SELECT COUNT(*) FROM post_views WHERE post_views.post_id = posts.id) as views_count").
		Where("status = ?", "published").
		Preload("Tags").Preload("Category")

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchKeyword, searchKeyword)
	}

	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Order("published_at DESC").Find(&posts).Error
	return posts, total, err
}

func (r *newsRepository) GetNewsByCategorySlug(categorySlug string, offset, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// Query harus join dengan kategori dan menghitung views_count
	query := r.db.Model(&domain.Post{}).
		Select("posts.*, (SELECT COUNT(*) FROM post_views WHERE post_views.post_id = posts.id) as views_count").
		Joins("JOIN categories ON categories.id = posts.category_id").
		Where("categories.slug = ? AND posts.status = ?", categorySlug, "published").
		Preload("Tags").Preload("Category")

	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Order("published_at DESC").Find(&posts).Error

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
