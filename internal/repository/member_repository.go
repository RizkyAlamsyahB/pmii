package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// MemberRepository interface untuk data access member
type MemberRepository interface {
	Create(member *domain.Member) error
	FindAll(page, limit int) ([]domain.Member, int64, error)
	FindByID(id int) (*domain.Member, error)
	Update(member *domain.Member) error
	Delete(id int) error
}

type memberRepository struct {
	db *gorm.DB
}

// NewMemberRepository constructor untuk MemberRepository
func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

// Create membuat member baru
func (r *memberRepository) Create(member *domain.Member) error {
	return r.db.Create(member).Error
}

// FindAll mengambil semua member dengan pagination
func (r *memberRepository) FindAll(page, limit int) ([]domain.Member, int64, error) {
	var members []domain.Member
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Member{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data (no ordering, natural database order)
	err := r.db.Limit(limit).Offset(offset).Find(&members).Error
	return members, total, err
}

// FindByID mengambil member berdasarkan ID
func (r *memberRepository) FindByID(id int) (*domain.Member, error) {
	var member domain.Member
	err := r.db.First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// Update mengupdate member
func (r *memberRepository) Update(member *domain.Member) error {
	return r.db.Save(member).Error
}

// Delete menghapus member (hard delete)
func (r *memberRepository) Delete(id int) error {
	return r.db.Delete(&domain.Member{}, id).Error
}
