package repository

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// AdRepository interface untuk abstraksi data layer ads
type AdRepository interface {
	FindAll() ([]domain.Ad, error)
	FindByID(id int) (*domain.Ad, error)
	FindByPage(page domain.AdPage) ([]domain.Ad, error)
	FindByPageAndSlot(page domain.AdPage, slot int) (*domain.Ad, error)
	Update(ad *domain.Ad) error
	UpdateImage(id int, imageURL string) error
}

type adRepository struct {
	db *gorm.DB
}

// NewAdRepository constructor untuk AdRepository
func NewAdRepository() AdRepository {
	return &adRepository{db: config.DB}
}

// FindAll mengambil semua ads
func (r *adRepository) FindAll() ([]domain.Ad, error) {
	var ads []domain.Ad
	if err := r.db.Order("page, slot").Find(&ads).Error; err != nil {
		return nil, err
	}
	return ads, nil
}

// FindByID mencari ad berdasarkan ID
func (r *adRepository) FindByID(id int) (*domain.Ad, error) {
	var ad domain.Ad
	if err := r.db.First(&ad, id).Error; err != nil {
		return nil, err
	}
	return &ad, nil
}

// FindByPage mengambil semua ads untuk page tertentu
func (r *adRepository) FindByPage(page domain.AdPage) ([]domain.Ad, error) {
	var ads []domain.Ad
	if err := r.db.Where("page = ?", page).Order("slot").Find(&ads).Error; err != nil {
		return nil, err
	}
	return ads, nil
}

// FindByPageAndSlot mencari ad berdasarkan page dan slot
func (r *adRepository) FindByPageAndSlot(page domain.AdPage, slot int) (*domain.Ad, error) {
	var ad domain.Ad
	if err := r.db.Where("page = ? AND slot = ?", page, slot).First(&ad).Error; err != nil {
		return nil, err
	}
	return &ad, nil
}

// Update mengupdate data ad
func (r *adRepository) Update(ad *domain.Ad) error {
	return r.db.Save(ad).Error
}

// UpdateImage mengupdate image_url ad
func (r *adRepository) UpdateImage(id int, imageURL string) error {
	return r.db.Model(&domain.Ad{}).Where("id = ?", id).Update("image_url", imageURL).Error
}
