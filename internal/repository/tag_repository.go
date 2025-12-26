package repository

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

type TagRepository interface {
	FindAll(offset, limit int, search string) ([]domain.Tag, int64, error)
	FindByID(id string) (domain.Tag, error)
	FindBySlug(slug string) (domain.Tag, error)
	Create(tag *domain.Tag) error
	Update(tag *domain.Tag) error
	Delete(tag *domain.Tag) error
}

type tagRepository struct{}

func NewTagRepository() TagRepository {
	return &tagRepository{}
}

func (r *tagRepository) FindAll(offset, limit int, search string) ([]domain.Tag, int64, error) {
	var tags []domain.Tag
	var total int64
	db := config.DB.Model(&domain.Tag{})

	// Implementasi Search menggunakan ILIKE
	if search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}

	db.Count(&total)
	err := db.Limit(limit).Offset(offset).Order("name ASC").Find(&tags).Error
	return tags, total, err
}

func (r *tagRepository) FindByID(id string) (domain.Tag, error) {
	var tag domain.Tag
	err := config.DB.First(&tag, id).Error
	return tag, err
}

func (r *tagRepository) FindBySlug(slug string) (domain.Tag, error) {
	var tag domain.Tag
	err := config.DB.Where("slug = ?", slug).First(&tag).Error
	return tag, err
}

func (r *tagRepository) Create(tag *domain.Tag) error {
	return config.DB.Create(tag).Error
}

func (r *tagRepository) Update(tag *domain.Tag) error {
	return config.DB.Save(tag).Error
}

func (r *tagRepository) Delete(tag *domain.Tag) error {
	return config.DB.Delete(tag).Error
}
