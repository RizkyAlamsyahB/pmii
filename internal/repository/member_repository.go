package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// MemberRepository interface untuk data access member
type MemberRepository interface {
	Create(member *domain.Member) error
	FindAll(page, limit int, search string) ([]domain.Member, int64, error)
	FindByID(id int) (*domain.Member, error)
	Update(member *domain.Member) error
	Delete(id int) error
	// Public methods
	FindActiveWithPagination(page, limit int, search string) ([]domain.Member, int64, error)
	FindActiveByDepartment(department string, page, limit int, search string) ([]domain.Member, int64, error)
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

// FindAll mengambil semua member dengan pagination dan search
func (r *memberRepository) FindAll(page, limit int, search string) ([]domain.Member, int64, error) {
	var members []domain.Member
	var total int64

	query := r.db.Model(&domain.Member{})

	// Search by full_name, position, department, atau social_links
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"full_name ILIKE ? OR position ILIKE ? OR department::text ILIKE ? OR social_links::text ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data (ASC - oldest first)
	err := query.Order("id ASC").Limit(limit).Offset(offset).Find(&members).Error
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

// FindActiveWithPagination mengambil member aktif dengan pagination dan search
func (r *memberRepository) FindActiveWithPagination(page, limit int, search string) ([]domain.Member, int64, error) {
	var members []domain.Member
	var total int64

	query := r.db.Model(&domain.Member{}).Where("is_active = ?", true)

	// Search by name or position
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR position ILIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data - ORDER BY id ASC as secondary sort for consistent pagination
	err := query.Order("created_at ASC, id ASC").Limit(limit).Offset(offset).Find(&members).Error
	return members, total, err
}

// FindActiveByDepartment mengambil member aktif berdasarkan department dengan pagination
func (r *memberRepository) FindActiveByDepartment(department string, page, limit int, search string) ([]domain.Member, int64, error) {
	var members []domain.Member
	var total int64

	query := r.db.Model(&domain.Member{}).Where("is_active = ? AND department = ?", true, department)

	// Search by name or position
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR position ILIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data - ORDER BY id ASC as secondary sort for consistent pagination
	err := query.Order("created_at ASC, id ASC").Limit(limit).Offset(offset).Find(&members).Error
	return members, total, err
}
