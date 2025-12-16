package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// AboutRepository interface untuk data access about
type AboutRepository interface {
	Get() (*domain.About, error)
	Upsert(about *domain.About) error
}

type aboutRepository struct {
	db *gorm.DB
}

// NewAboutRepository constructor untuk AboutRepository
func NewAboutRepository(db *gorm.DB) AboutRepository {
	return &aboutRepository{db: db}
}

// Get mengambil data about (singleton - hanya ada 1 record)
func (r *aboutRepository) Get() (*domain.About, error) {
	var about domain.About
	err := r.db.First(&about).Error
	if err != nil {
		return nil, err
	}
	return &about, nil
}

// Upsert membuat atau mengupdate data about
// Jika sudah ada record, update. Jika belum ada, create.
func (r *aboutRepository) Upsert(about *domain.About) error {
	// Cek apakah sudah ada record
	var existing domain.About
	err := r.db.First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Belum ada record, create baru
		return r.db.Create(about).Error
	}

	if err != nil {
		return err
	}

	// Sudah ada record, update dengan ID yang sama
	about.ID = existing.ID
	return r.db.Save(about).Error
}
