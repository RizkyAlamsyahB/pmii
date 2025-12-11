package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// TestimonialRepository interface untuk data access testimonial
type TestimonialRepository interface {
	Create(testimonial *domain.Testimonial) error
	FindAll() ([]domain.Testimonial, error)
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

// FindAll mengambil semua testimonial
func (r *testimonialRepository) FindAll() ([]domain.Testimonial, error) {
	var testimonials []domain.Testimonial
	err := r.db.Order("created_at DESC").Find(&testimonials).Error
	return testimonials, err
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
