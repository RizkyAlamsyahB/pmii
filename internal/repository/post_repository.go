package repository

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

type PostRepository interface {
	FindAll(offset, limit int, search string) ([]domain.Post, int64, error)
	FindByID(id int) (domain.Post, error)
	FindBySlugOrID(identifier string) (domain.Post, error)
	Create(post *domain.Post) error
	Update(post *domain.Post) error
	Delete(post *domain.Post, unscoped bool) error
	GetTagBySlug(slug string, name string) (domain.Tag, error)
}

type postRepository struct{}

func NewPostRepository() PostRepository {
	return &postRepository{}
}

func (r *postRepository) FindAll(offset, limit int, search string) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64
	query := config.DB.Model(&domain.Post{}).Preload("Tags").Preload("Category")
	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchKeyword, searchKeyword)
	}
	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Order("published_at DESC").Find(&posts).Error
	return posts, total, err
}

func (r *postRepository) FindByID(id int) (domain.Post, error) {
	var post domain.Post
	err := config.DB.Preload("Category").Preload("Tags").First(&post, id).Error
	return post, err
}

func (r *postRepository) FindBySlugOrID(identifier string) (domain.Post, error) {
	var post domain.Post
	query := config.DB.Preload("Category").Preload("Tags")
	err := query.Where("id = ? OR slug = ?", identifier, identifier).First(&post).Error
	return post, err
}

func (r *postRepository) Create(post *domain.Post) error {
	return config.DB.Create(post).Error
}

func (r *postRepository) Update(post *domain.Post) error {
	return config.DB.Save(post).Error
}

func (r *postRepository) Delete(post *domain.Post, unscoped bool) error {
	db := config.DB
	if unscoped {
		db = db.Unscoped()
	}
	return db.Delete(post).Error
}

func (r *postRepository) GetTagBySlug(slug string, name string) (domain.Tag, error) {
	var tag domain.Tag
	err := config.DB.Where(domain.Tag{Slug: slug}).Attrs(domain.Tag{Name: name}).FirstOrCreate(&tag).Error
	return tag, err
}
