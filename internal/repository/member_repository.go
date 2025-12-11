package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// MemberRepository interface untuk data access member
type MemberRepository interface {
	Create(member *domain.Member) error
	FindAll() ([]domain.Member, error)
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

// FindAll mengambil semua member
func (r *memberRepository) FindAll() ([]domain.Member, error) {
	var members []domain.Member
	err := r.db.Order("created_at DESC").Find(&members).Error
	return members, err
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
