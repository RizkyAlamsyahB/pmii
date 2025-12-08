package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// UserRepository interface untuk abstraksi data layer
// Memudahkan mocking saat unit testing
type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository constructor untuk UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByEmail mencari user berdasarkan email
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	// GORM otomatis melakukan parameter binding (Anti SQL Injection)
	// Query ke kolom user_email di database (bukan email)
	if err := r.db.Where("user_email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID mencari user berdasarkan ID
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create membuat user baru
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// Update mengupdate data user
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete menghapus user berdasarkan ID
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
