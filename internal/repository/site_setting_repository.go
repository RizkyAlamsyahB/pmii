package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// SiteSettingRepository interface untuk data access site settings
type SiteSettingRepository interface {
	Get() (*domain.SiteSetting, error)
	Update(setting *domain.SiteSetting) error
}

type siteSettingRepository struct {
	db *gorm.DB
}

// NewSiteSettingRepository constructor untuk SiteSettingRepository
func NewSiteSettingRepository(db *gorm.DB) SiteSettingRepository {
	return &siteSettingRepository{db: db}
}

// Get mengambil site settings (singleton - always ID 1)
func (r *siteSettingRepository) Get() (*domain.SiteSetting, error) {
	var setting domain.SiteSetting
	err := r.db.First(&setting, 1).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// Update mengupdate atau create site settings
func (r *siteSettingRepository) Update(setting *domain.SiteSetting) error {
	// Check if record exists
	var existing domain.SiteSetting
	err := r.db.First(&existing, 1).Error

	if err != nil {
		// Create new record with ID 1
		setting.ID = 1
		return r.db.Create(setting).Error
	}

	// Update existing record
	setting.ID = 1
	return r.db.Save(setting).Error
}
