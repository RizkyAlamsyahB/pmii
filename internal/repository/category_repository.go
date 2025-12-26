package repository

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

type CategoryRepository interface {
	FindAll(offset, limit int, search string) ([]domain.Category, int64, error)
	FindByID(id string) (domain.Category, error)
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(category *domain.Category) error
}

type categoryRepository struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

func (r *categoryRepository) FindAll(offset, limit int, search string) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64
	db := config.DB.Model(&domain.Category{})

	if search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}

	db.Count(&total)
	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&categories).Error
	return categories, total, err
}

func (r *categoryRepository) FindByID(id string) (domain.Category, error) {
	var category domain.Category
	err := config.DB.First(&category, id).Error
	return category, err
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return config.DB.Create(category).Error
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return config.DB.Save(category).Error
}

func (r *categoryRepository) Delete(category *domain.Category) error {
	return config.DB.Delete(category).Error
}
