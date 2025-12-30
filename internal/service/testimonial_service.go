package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// CloudinaryService interface untuk operasi cloudinary
type CloudinaryService interface {
	// Image-specific methods (for backward compatibility)
	UploadImage(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteImage(ctx context.Context, folder string, filename string) error
	GetImageURL(folder string, filename string) string
	// Generic file methods (auto-detect resource type: image/video/raw)
	UploadFile(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteFile(ctx context.Context, folder string, filename string) error
	GetFileURL(folder string, filename string) string
	GetDownloadURL(folder string, filename string) string
}

// TestimonialService interface untuk business logic testimonial
type TestimonialService interface {
	Create(ctx context.Context, req requests.CreateTestimonialRequest, photoFile *multipart.FileHeader) (*responses.TestimonialResponse, error)
	GetAll(ctx context.Context, page, limit int, search string) ([]responses.TestimonialResponse, int, int, int64, error)
	GetByID(ctx context.Context, id int) (*responses.TestimonialResponse, error)
	Update(ctx context.Context, id int, req requests.UpdateTestimonialRequest, photoFile *multipart.FileHeader) (*responses.TestimonialResponse, error)
	Delete(ctx context.Context, id int) error
}

type testimonialService struct {
	testimonialRepo   repository.TestimonialRepository
	cloudinaryService CloudinaryService
	activityLogRepo   repository.ActivityLogRepository
}

// NewTestimonialService constructor untuk TestimonialService
func NewTestimonialService(testimonialRepo repository.TestimonialRepository, cloudinaryService CloudinaryService, activityLogRepo repository.ActivityLogRepository) TestimonialService {
	return &testimonialService{
		testimonialRepo:   testimonialRepo,
		cloudinaryService: cloudinaryService,
		activityLogRepo:   activityLogRepo,
	}
}

// Create membuat testimonial baru dengan upload foto ke Cloudinary
func (s *testimonialService) Create(ctx context.Context, req requests.CreateTestimonialRequest, photoFile *multipart.FileHeader) (*responses.TestimonialResponse, error) {
	// Upload photo ke Cloudinary (jika ada)
	var photoFilename *string
	if photoFile != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "testimonials", photoFile)
		if err != nil {
			return nil, errors.New("gagal mengupload foto")
		}
		photoFilename = &filename
	}

	// Prepare domain entity
	var org *string
	if req.Organization != "" {
		org = &req.Organization
	}

	var pos *string
	if req.Position != "" {
		pos = &req.Position
	}

	testimonial := &domain.Testimonial{
		Name:         req.Name,
		Organization: org,
		Position:     pos,
		Content:      req.Content,
		PhotoURI:     photoFilename,
		IsActive:     true,
	}

	// Save ke database
	if err := s.testimonialRepo.Create(testimonial); err != nil {
		// Rollback: hapus foto dari Cloudinary jika save gagal
		if photoFilename != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "testimonials", *photoFilename)
		}
		return nil, errors.New("gagal menyimpan testimonial")
	}

	// Convert to response DTO
	resp := s.toResponseDTO(testimonial)

	// Log activity - Create Testimonial
	s.logActivity(ctx, domain.ActionCreate, domain.ModuleTestimoni, "Membuat testimonial baru: "+testimonial.Name, nil, map[string]any{
		"id":   testimonial.ID,
		"name": testimonial.Name,
	}, &testimonial.ID)

	return resp, nil
}

// GetAll mengambil semua testimonial dengan pagination dan search
func (s *testimonialService) GetAll(ctx context.Context, page, limit int, search string) ([]responses.TestimonialResponse, int, int, int64, error) {
	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	testimonials, total, err := s.testimonialRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, 0, 0, 0, errors.New("gagal mengambil data testimonial")
	}

	// Calculate last page
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	// Auto-clamp: jika page melebihi lastPage dan ada data, clamp ke lastPage
	if page > lastPage && lastPage > 0 {
		page = lastPage
		// Re-fetch dengan page yang sudah di-clamp
		testimonials, _, err = s.testimonialRepo.FindAll(page, limit, search)
		if err != nil {
			return nil, 0, 0, 0, errors.New("gagal mengambil data testimonial")
		}
	}

	// Convert to response DTOs
	result := make([]responses.TestimonialResponse, len(testimonials))
	for i, t := range testimonials {
		result[i] = *s.toResponseDTO(&t)
	}

	return result, page, lastPage, total, nil
}

// GetByID mengambil testimonial berdasarkan ID
func (s *testimonialService) GetByID(ctx context.Context, id int) (*responses.TestimonialResponse, error) {
	testimonial, err := s.testimonialRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("testimonial tidak ditemukan")
	}

	return s.toResponseDTO(testimonial), nil
}

// Update mengupdate testimonial dengan optional upload foto baru
func (s *testimonialService) Update(ctx context.Context, id int, req requests.UpdateTestimonialRequest, photoFile *multipart.FileHeader) (*responses.TestimonialResponse, error) {
	// Ambil testimonial existing
	testimonial, err := s.testimonialRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("testimonial tidak ditemukan")
	}

	// Simpan foto lama untuk rollback
	oldPhotoURI := testimonial.PhotoURI

	// Upload foto baru ke Cloudinary (jika ada)
	var newPhotoFilename *string
	if photoFile != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "testimonials", photoFile)
		if err != nil {
			return nil, errors.New("gagal mengupload foto")
		}
		newPhotoFilename = &filename
		testimonial.PhotoURI = &filename
	}

	// Update fields yang dikirim
	if req.Name != "" {
		testimonial.Name = req.Name
	}
	if req.Organization != "" {
		testimonial.Organization = &req.Organization
	}
	if req.Position != "" {
		testimonial.Position = &req.Position
	}
	if req.Content != "" {
		testimonial.Content = req.Content
	}
	if req.IsActive != nil {
		testimonial.IsActive = *req.IsActive
	}

	// Save ke database
	if err := s.testimonialRepo.Update(testimonial); err != nil {
		// Rollback: hapus foto baru jika update gagal
		if newPhotoFilename != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "testimonials", *newPhotoFilename)
		}
		return nil, errors.New("gagal mengupdate testimonial")
	}

	// Hapus foto lama SETELAH database update berhasil
	if newPhotoFilename != nil && oldPhotoURI != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "testimonials", *oldPhotoURI)
	}

	// Log activity - Update Testimonial
	s.logActivity(ctx, domain.ActionUpdate, domain.ModuleTestimoni, "Mengupdate testimonial: "+testimonial.Name, nil, map[string]any{
		"id":   testimonial.ID,
		"name": testimonial.Name,
	}, &testimonial.ID)

	return s.toResponseDTO(testimonial), nil
}

// Delete menghapus testimonial dan foto dari Cloudinary
func (s *testimonialService) Delete(ctx context.Context, id int) error {
	// Ambil testimonial untuk mendapatkan info foto
	testimonial, err := s.testimonialRepo.FindByID(id)
	if err != nil {
		return errors.New("testimonial tidak ditemukan")
	}

	// Log activity sebelum delete
	s.logActivity(ctx, domain.ActionDelete, domain.ModuleTestimoni, "Menghapus testimonial: "+testimonial.Name, map[string]any{
		"id":   testimonial.ID,
		"name": testimonial.Name,
	}, nil, &testimonial.ID)

	// Hapus dari database
	if err := s.testimonialRepo.Delete(id); err != nil {
		return errors.New("gagal menghapus testimonial")
	}

	// Hapus foto dari Cloudinary (jika ada)
	if testimonial.PhotoURI != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "testimonials", *testimonial.PhotoURI)
	}

	return nil
}

// toResponseDTO converts domain.Testimonial to responses.TestimonialResponse
func (s *testimonialService) toResponseDTO(t *domain.Testimonial) *responses.TestimonialResponse {
	var imageURL string
	if t.PhotoURI != nil {
		imageURL = s.cloudinaryService.GetImageURL("testimonials", *t.PhotoURI)
	}

	return &responses.TestimonialResponse{
		ID:           t.ID,
		Name:         t.Name,
		Organization: t.Organization,
		Position:     t.Position,
		Content:      t.Content,
		ImageUrl:     imageURL,
		IsActive:     t.IsActive,
		CreatedAt:    t.CreatedAt,
	}
}

// logActivity helper untuk mencatat activity log
func (s *testimonialService) logActivity(ctx context.Context, actionType domain.ActivityActionType, module domain.ActivityModuleType, description string, oldValue, newValue map[string]any, targetID *int) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return // Skip if no user in context
	}

	ipAddress := utils.GetIPAddress(ctx)
	userAgent := utils.GetUserAgent(ctx)

	var ipPtr, uaPtr *string
	if ipAddress != "" {
		ipPtr = &ipAddress
	}
	if userAgent != "" {
		uaPtr = &userAgent
	}

	log := &domain.ActivityLog{
		UserID:      userID,
		ActionType:  actionType,
		Module:      module,
		Description: &description,
		TargetID:    targetID,
		OldValue:    oldValue,
		NewValue:    newValue,
		IPAddress:   ipPtr,
		UserAgent:   uaPtr,
	}

	// Ignore error - logging should not affect main operation
	_ = s.activityLogRepo.Create(log)
}
