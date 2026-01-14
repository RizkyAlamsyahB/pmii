package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// ContactRepository interface untuk data access contact
type ContactRepository interface {
	Get() (*domain.Contact, error)
	Update(contact *domain.Contact) error
}

type contactRepository struct {
	db *gorm.DB
}

// NewContactRepository constructor untuk ContactRepository
func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

// Get mengambil contact info (singleton - mengambil record pertama yang ada)
func (r *contactRepository) Get() (*domain.Contact, error) {
	var contact domain.Contact
	err := r.db.First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update mengupdate atau create contact info
func (r *contactRepository) Update(contact *domain.Contact) error {
	// Check if record exists
	var existing domain.Contact
	err := r.db.First(&existing, 1).Error

	if err != nil {
		// Create new record with ID 1
		contact.ID = 1
		return r.db.Create(contact).Error
	}

	// Update existing record
	contact.ID = 1
	return r.db.Save(contact).Error
}
