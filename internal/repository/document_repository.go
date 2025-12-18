package repository

import (
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// DocumentRepository interface untuk data access document
type DocumentRepository interface {
	Create(document *domain.Document) error
	FindAll(page, limit int, fileType, search string) ([]domain.Document, int64, error)
	FindByID(id int) (*domain.Document, error)
	Update(document *domain.Document) error
	Delete(id int) error
	// Public methods
	FindAllActive(fileType string) ([]domain.Document, error)
	FindActiveByType(fileType string) ([]domain.Document, error)
}

type documentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository constructor untuk DocumentRepository
func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

// Create membuat document baru
func (r *documentRepository) Create(document *domain.Document) error {
	return r.db.Create(document).Error
}

// FindAll mengambil semua document dengan pagination, filter by type, dan search
func (r *documentRepository) FindAll(page, limit int, fileType, search string) ([]domain.Document, int64, error) {
	var documents []domain.Document
	var total int64

	query := r.db.Model(&domain.Document{}).Where("deleted_at IS NULL")

	// Filter by file type if provided
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	// Search by name
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ?", searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&documents).Error
	return documents, total, err
}

// FindByID mengambil document berdasarkan ID (tidak termasuk soft deleted)
func (r *documentRepository) FindByID(id int) (*domain.Document, error) {
	var document domain.Document
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// Update mengupdate document
func (r *documentRepository) Update(document *domain.Document) error {
	return r.db.Save(document).Error
}

// Delete melakukan soft delete pada document
func (r *documentRepository) Delete(id int) error {
	return r.db.Model(&domain.Document{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindAllActive mengambil semua document aktif (untuk public)
func (r *documentRepository) FindAllActive(fileType string) ([]domain.Document, error) {
	var documents []domain.Document
	query := r.db.Where("deleted_at IS NULL")

	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	err := query.Order("created_at DESC").Find(&documents).Error
	return documents, err
}

// FindActiveByType mengambil document aktif berdasarkan type (untuk public)
func (r *documentRepository) FindActiveByType(fileType string) ([]domain.Document, error) {
	var documents []domain.Document
	err := r.db.Where("deleted_at IS NULL AND file_type = ?", fileType).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}
