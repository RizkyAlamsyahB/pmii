package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// TestimonialRepository interface untuk data access testimonial
type TestimonialRepository interface {
	Create(testimonial *domain.Testimonial) error
	FindAll(page, limit int, search string) ([]domain.Testimonial, int64, error)
	FindByID(id int) (*domain.Testimonial, error)
	Update(testimonial *domain.Testimonial) error
	Delete(id int) error
}

type testimonialRepository struct {
	db *gorm.DB
}

// NewTestimonialRepository constructor untuk TestimonialRepository
func NewTestimonialRepository(db *gorm.DB) TestimonialRepository {
	return &testimonialRepository{db: db}
}

// Create membuat testimonial baru
func (r *testimonialRepository) Create(testimonial *domain.Testimonial) error {
	return r.db.Create(testimonial).Error
}

// FindAll mengambil semua testimonial dengan pagination dan search
func (r *testimonialRepository) FindAll(page, limit int, search string) ([]domain.Testimonial, int64, error) {
	var testimonials []domain.Testimonial
	var total int64

	query := r.db.Model(&domain.Testimonial{})

	// Search by name, organization, position, atau content
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"name ILIKE ? OR organization ILIKE ? OR position ILIKE ? OR content ILIKE ?",
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
	err := query.Order("id ASC").Limit(limit).Offset(offset).Find(&testimonials).Error
	return testimonials, total, err
}

// FindByID mengambil testimonial berdasarkan ID
func (r *testimonialRepository) FindByID(id int) (*domain.Testimonial, error) {
	var testimonial domain.Testimonial
	err := r.db.First(&testimonial, id).Error
	if err != nil {
		return nil, err
	}
	return &testimonial, nil
}

// Update mengupdate testimonial
func (r *testimonialRepository) Update(testimonial *domain.Testimonial) error {
	return r.db.Save(testimonial).Error
}

// Delete menghapus testimonial (hard delete)
func (r *testimonialRepository) Delete(id int) error {
	return r.db.Delete(&domain.Testimonial{}, id).Error
}
