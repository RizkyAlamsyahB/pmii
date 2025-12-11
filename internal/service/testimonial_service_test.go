package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
)

// MockTestimonialRepository adalah mock untuk TestimonialRepository
type MockTestimonialRepository struct {
	CreateFunc   func(testimonial *domain.Testimonial) error
	FindAllFunc  func() ([]domain.Testimonial, error)
	FindByIDFunc func(id int) (*domain.Testimonial, error)
	UpdateFunc   func(testimonial *domain.Testimonial) error
	DeleteFunc   func(id int) error
}

func (m *MockTestimonialRepository) Create(testimonial *domain.Testimonial) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(testimonial)
	}
	return errors.New("mock not configured")
}

func (m *MockTestimonialRepository) FindAll() ([]domain.Testimonial, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockTestimonialRepository) FindByID(id int) (*domain.Testimonial, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockTestimonialRepository) Update(testimonial *domain.Testimonial) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(testimonial)
	}
	return errors.New("mock not configured")
}

func (m *MockTestimonialRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return errors.New("mock not configured")
}

// MockCloudinaryService adalah mock untuk Cloudinary Service
type MockCloudinaryService struct {
	UploadImageFunc func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteImageFunc func(ctx context.Context, folder string, filename string) error
	GetImageURLFunc func(folder string, filename string) string
}

func (m *MockCloudinaryService) UploadImage(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadImageFunc != nil {
		return m.UploadImageFunc(ctx, folder, file)
	}
	return "", errors.New("mock not configured")
}

func (m *MockCloudinaryService) DeleteImage(ctx context.Context, folder string, filename string) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, folder, filename)
	}
	return errors.New("mock not configured")
}

func (m *MockCloudinaryService) GetImageURL(folder string, filename string) string {
	if m.GetImageURLFunc != nil {
		return m.GetImageURLFunc(folder, filename)
	}
	return ""
}

// Helper function untuk membuat pointer string
func stringPtr(s string) *string {
	return &s
}

// Helper function untuk membuat pointer bool
func boolPtr(b bool) *bool {
	return &b
}

// ==================== CREATE TESTS ====================

func TestCreate_Success_WithPhoto(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		CreateFunc: func(testimonial *domain.Testimonial) error {
			testimonial.ID = 1
			testimonial.CreatedAt = time.Now()
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "photo123.jpg", nil
		},
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/testimonials/photo123.jpg"
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.CreateTestimonialRequest{
		Name:         "John Doe",
		Organization: "PT ABC",
		Position:     "CEO",
		Content:      "Great service!",
	}

	// Mock file upload
	mockFile := &multipart.FileHeader{
		Filename: "photo.jpg",
		Size:     1024,
	}

	result, err := service.Create(context.Background(), req, mockFile)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.TestimonialName != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", result.TestimonialName)
	}
	if result.TestimonialImage != "https://cloudinary.com/testimonials/photo123.jpg" {
		t.Errorf("Expected image URL, got '%s'", result.TestimonialImage)
	}
}

func TestCreate_Success_WithoutPhoto(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		CreateFunc: func(testimonial *domain.Testimonial) error {
			testimonial.ID = 1
			testimonial.CreatedAt = time.Now()
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.CreateTestimonialRequest{
		Name:         "Jane Doe",
		Organization: "PT XYZ",
		Position:     "Manager",
		Content:      "Excellent!",
	}

	result, err := service.Create(context.Background(), req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.TestimonialName != "Jane Doe" {
		t.Errorf("Expected name 'Jane Doe', got '%s'", result.TestimonialName)
	}
	if result.TestimonialImage != "" {
		t.Errorf("Expected empty image URL, got '%s'", result.TestimonialImage)
	}
}

func TestCreate_Error_UploadFailed(t *testing.T) {
	mockRepo := &MockTestimonialRepository{}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.CreateTestimonialRequest{
		Name:    "John Doe",
		Content: "Great service!",
	}

	mockFile := &multipart.FileHeader{
		Filename: "photo.jpg",
		Size:     1024,
	}

	result, err := service.Create(context.Background(), req, mockFile)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
	if err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected 'gagal mengupload foto', got: %v", err.Error())
	}
}

func TestCreate_Error_DatabaseFailed_WithRollback(t *testing.T) {
	uploadCalled := false
	deleteCalled := false

	mockRepo := &MockTestimonialRepository{
		CreateFunc: func(testimonial *domain.Testimonial) error {
			return errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			uploadCalled = true
			return "photo123.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			deleteCalled = true
			if filename != "photo123.jpg" {
				t.Errorf("Expected to delete 'photo123.jpg', got '%s'", filename)
			}
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.CreateTestimonialRequest{
		Name:    "John Doe",
		Content: "Great service!",
	}

	mockFile := &multipart.FileHeader{
		Filename: "photo.jpg",
		Size:     1024,
	}

	result, err := service.Create(context.Background(), req, mockFile)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
	if err.Error() != "gagal menyimpan testimonial" {
		t.Errorf("Expected 'gagal menyimpan testimonial', got: %v", err.Error())
	}
	if !uploadCalled {
		t.Error("Expected upload to be called")
	}
	if !deleteCalled {
		t.Error("Expected delete to be called for rollback")
	}
}

func TestCreate_Success_EmptyOptionalFields(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		CreateFunc: func(testimonial *domain.Testimonial) error {
			// Validate that optional fields are nil when empty
			if testimonial.Organization != nil {
				t.Error("Expected Organization to be nil")
			}
			if testimonial.Position != nil {
				t.Error("Expected Position to be nil")
			}
			testimonial.ID = 1
			testimonial.CreatedAt = time.Now()
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.CreateTestimonialRequest{
		Name:    "John Doe",
		Content: "Great service!",
		// Organization and Position are empty strings
	}

	result, err := service.Create(context.Background(), req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

// ==================== GET ALL TESTS ====================

func TestGetAll_Success(t *testing.T) {
	now := time.Now()
	org := "PT ABC"
	pos := "CEO"
	photo := "photo1.jpg"

	mockRepo := &MockTestimonialRepository{
		FindAllFunc: func() ([]domain.Testimonial, error) {
			return []domain.Testimonial{
				{
					ID:           1,
					Name:         "John Doe",
					Organization: &org,
					Position:     &pos,
					Content:      "Great!",
					PhotoURI:     &photo,
					IsActive:     true,
					CreatedAt:    now,
				},
				{
					ID:           2,
					Name:         "Jane Doe",
					Organization: nil,
					Position:     nil,
					Content:      "Excellent!",
					PhotoURI:     nil,
					IsActive:     false,
					CreatedAt:    now,
				},
			}, nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			if filename == "photo1.jpg" {
				return "https://cloudinary.com/testimonials/photo1.jpg"
			}
			return ""
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	results, err := service.GetAll(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Validate first testimonial
	if results[0].TestimonialName != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", results[0].TestimonialName)
	}
	if results[0].TestimonialImage != "https://cloudinary.com/testimonials/photo1.jpg" {
		t.Errorf("Expected image URL, got '%s'", results[0].TestimonialImage)
	}
	if !results[0].TestimonialIsActive {
		t.Error("Expected first testimonial to be active")
	}

	// Validate second testimonial
	if results[1].TestimonialName != "Jane Doe" {
		t.Errorf("Expected name 'Jane Doe', got '%s'", results[1].TestimonialName)
	}
	if results[1].TestimonialImage != "" {
		t.Errorf("Expected empty image URL, got '%s'", results[1].TestimonialImage)
	}
	if results[1].TestimonialIsActive {
		t.Error("Expected second testimonial to be inactive")
	}
}

func TestGetAll_Error_DatabaseFailed(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindAllFunc: func() ([]domain.Testimonial, error) {
			return nil, errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	results, err := service.GetAll(context.Background())

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if results != nil {
		t.Errorf("Expected nil results, got: %v", results)
	}
	if err.Error() != "gagal mengambil data testimonial" {
		t.Errorf("Expected 'gagal mengambil data testimonial', got: %v", err.Error())
	}
}

func TestGetAll_Success_EmptyList(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindAllFunc: func() ([]domain.Testimonial, error) {
			return []domain.Testimonial{}, nil
		},
	}

	mockCloudinary := &MockCloudinaryService{}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	results, err := service.GetAll(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected empty list, got %d items", len(results))
	}
}

// ==================== GET BY ID TESTS ====================

func TestGetByID_Success(t *testing.T) {
	now := time.Now()
	org := "PT ABC"
	pos := "CEO"
	photo := "photo1.jpg"

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			if id == 1 {
				return &domain.Testimonial{
					ID:           1,
					Name:         "John Doe",
					Organization: &org,
					Position:     &pos,
					Content:      "Great!",
					PhotoURI:     &photo,
					IsActive:     true,
					CreatedAt:    now,
				}, nil
			}
			return nil, errors.New("not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/testimonials/photo1.jpg"
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	result, err := service.GetByID(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.TestimonialID != 1 {
		t.Errorf("Expected ID 1, got %d", result.TestimonialID)
	}
	if result.TestimonialName != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", result.TestimonialName)
	}
}

func TestGetByID_Error_NotFound(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	result, err := service.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
	if err.Error() != "testimonial tidak ditemukan" {
		t.Errorf("Expected 'testimonial tidak ditemukan', got: %v", err.Error())
	}
}

// ==================== UPDATE TESTS ====================

func TestUpdate_Success_WithNewPhoto(t *testing.T) {
	now := time.Now()
	oldPhoto := "old_photo.jpg"
	org := "PT ABC"
	pos := "CEO"

	deleteOldPhotoCalled := false

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:           1,
				Name:         "John Doe",
				Organization: &org,
				Position:     &pos,
				Content:      "Great!",
				PhotoURI:     &oldPhoto,
				IsActive:     true,
				CreatedAt:    now,
			}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "new_photo.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			deleteOldPhotoCalled = true
			if filename != "old_photo.jpg" {
				t.Errorf("Expected to delete 'old_photo.jpg', got '%s'", filename)
			}
			return nil
		},
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/testimonials/new_photo.jpg"
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.UpdateTestimonialRequest{
		Name:    "John Updated",
		Content: "Updated content",
	}

	mockFile := &multipart.FileHeader{
		Filename: "new_photo.jpg",
		Size:     2048,
	}

	result, err := service.Update(context.Background(), 1, req, mockFile)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.TestimonialName != "John Updated" {
		t.Errorf("Expected name 'John Updated', got '%s'", result.TestimonialName)
	}
	if !deleteOldPhotoCalled {
		t.Error("Expected old photo to be deleted")
	}
}

func TestUpdate_Success_WithoutPhoto(t *testing.T) {
	now := time.Now()
	photo := "photo.jpg"
	org := "PT ABC"

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:           1,
				Name:         "John Doe",
				Organization: &org,
				Position:     nil,
				Content:      "Great!",
				PhotoURI:     &photo,
				IsActive:     true,
				CreatedAt:    now,
			}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/testimonials/photo.jpg"
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	isActive := false
	req := requests.UpdateTestimonialRequest{
		Name:     "John Updated",
		Content:  "Updated content",
		IsActive: &isActive,
	}

	result, err := service.Update(context.Background(), 1, req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.TestimonialName != "John Updated" {
		t.Errorf("Expected name 'John Updated', got '%s'", result.TestimonialName)
	}
	if result.TestimonialContent != "Updated content" {
		t.Errorf("Expected content 'Updated content', got '%s'", result.TestimonialContent)
	}
	if result.TestimonialIsActive {
		t.Error("Expected IsActive to be false")
	}
}

func TestUpdate_Success_PartialUpdate(t *testing.T) {
	now := time.Now()
	org := "PT ABC"
	pos := "CEO"

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:           1,
				Name:         "John Doe",
				Organization: &org,
				Position:     &pos,
				Content:      "Original content",
				PhotoURI:     nil,
				IsActive:     true,
				CreatedAt:    now,
			}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			// Verify only Name is updated, others remain the same
			if testimonial.Name != "Updated Name" {
				t.Errorf("Expected name to be updated to 'Updated Name', got '%s'", testimonial.Name)
			}
			if testimonial.Content != "Original content" {
				t.Errorf("Expected content to remain 'Original content', got '%s'", testimonial.Content)
			}
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	// Only update Name field
	req := requests.UpdateTestimonialRequest{
		Name: "Updated Name",
	}

	result, err := service.Update(context.Background(), 1, req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.TestimonialName != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", result.TestimonialName)
	}
	if result.TestimonialContent != "Original content" {
		t.Errorf("Expected content to remain 'Original content', got '%s'", result.TestimonialContent)
	}
}

func TestUpdate_Error_NotFound(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.UpdateTestimonialRequest{
		Name: "Updated Name",
	}

	result, err := service.Update(context.Background(), 999, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
	if err.Error() != "testimonial tidak ditemukan" {
		t.Errorf("Expected 'testimonial tidak ditemukan', got: %v", err.Error())
	}
}

func TestUpdate_Error_UploadFailed(t *testing.T) {
	now := time.Now()
	org := "PT ABC"

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:           1,
				Name:         "John Doe",
				Organization: &org,
				Content:      "Great!",
				PhotoURI:     nil,
				IsActive:     true,
				CreatedAt:    now,
			}, nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.UpdateTestimonialRequest{
		Name: "Updated Name",
	}

	mockFile := &multipart.FileHeader{
		Filename: "photo.jpg",
		Size:     1024,
	}

	result, err := service.Update(context.Background(), 1, req, mockFile)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
	if err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected 'gagal mengupload foto', got: %v", err.Error())
	}
}

func TestUpdate_Error_DatabaseFailed_WithRollback(t *testing.T) {
	now := time.Now()
	oldPhoto := "old_photo.jpg"

	deleteNewPhotoCalled := false
	deleteOldPhotoCalled := false

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:        1,
				Name:      "John Doe",
				Content:   "Great!",
				PhotoURI:  &oldPhoto,
				IsActive:  true,
				CreatedAt: now,
			}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			return errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "new_photo.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			if filename == "new_photo.jpg" {
				deleteNewPhotoCalled = true
			}
			if filename == "old_photo.jpg" {
				deleteOldPhotoCalled = true
			}
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	req := requests.UpdateTestimonialRequest{
		Name: "Updated Name",
	}

	mockFile := &multipart.FileHeader{
		Filename: "new_photo.jpg",
		Size:     2048,
	}

	result, err := service.Update(context.Background(), 1, req, mockFile)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
	if err.Error() != "gagal mengupdate testimonial" {
		t.Errorf("Expected 'gagal mengupdate testimonial', got: %v", err.Error())
	}
	if !deleteNewPhotoCalled {
		t.Error("Expected new photo to be deleted for rollback")
	}
	if deleteOldPhotoCalled {
		t.Error("Old photo should NOT be deleted when database fails")
	}
}

// ==================== DELETE TESTS ====================

func TestDelete_Success_WithPhoto(t *testing.T) {
	photo := "photo.jpg"
	deletePhotoCalled := false

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:       1,
				Name:     "John Doe",
				Content:  "Great!",
				PhotoURI: &photo,
				IsActive: true,
			}, nil
		},
		DeleteFunc: func(id int) error {
			if id != 1 {
				t.Errorf("Expected to delete ID 1, got %d", id)
			}
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			deletePhotoCalled = true
			if filename != "photo.jpg" {
				t.Errorf("Expected to delete 'photo.jpg', got '%s'", filename)
			}
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !deletePhotoCalled {
		t.Error("Expected photo to be deleted")
	}
}

func TestDelete_Success_WithoutPhoto(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:       1,
				Name:     "John Doe",
				Content:  "Great!",
				PhotoURI: nil,
				IsActive: true,
			}, nil
		},
		DeleteFunc: func(id int) error {
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			t.Error("DeleteImage should not be called when PhotoURI is nil")
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDelete_Error_NotFound(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 999)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "testimonial tidak ditemukan" {
		t.Errorf("Expected 'testimonial tidak ditemukan', got: %v", err.Error())
	}
}

func TestDelete_Error_DatabaseFailed(t *testing.T) {
	photo := "photo.jpg"

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:       1,
				Name:     "John Doe",
				Content:  "Great!",
				PhotoURI: &photo,
				IsActive: true,
			}, nil
		},
		DeleteFunc: func(id int) error {
			return errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			t.Error("DeleteImage should not be called when database delete fails")
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 1)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "gagal menghapus testimonial" {
		t.Errorf("Expected 'gagal menghapus testimonial', got: %v", err.Error())
	}
}

// ==================== EDGE CASES ====================

func TestUpdate_Success_UpdateIsActiveToTrue(t *testing.T) {
	now := time.Now()

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:        1,
				Name:      "John Doe",
				Content:   "Great!",
				PhotoURI:  nil,
				IsActive:  false,
				CreatedAt: now,
			}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			if !testimonial.IsActive {
				t.Error("Expected IsActive to be true")
			}
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	isActive := true
	req := requests.UpdateTestimonialRequest{
		IsActive: &isActive,
	}

	result, err := service.Update(context.Background(), 1, req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if !result.TestimonialIsActive {
		t.Error("Expected IsActive to be true in response")
	}
}

func TestUpdate_Success_UpdateAllFields(t *testing.T) {
	now := time.Now()

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:           1,
				Name:         "Old Name",
				Organization: stringPtr("Old Org"),
				Position:     stringPtr("Old Pos"),
				Content:      "Old Content",
				PhotoURI:     nil,
				IsActive:     true,
				CreatedAt:    now,
			}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			if testimonial.Name != "New Name" {
				t.Errorf("Expected Name 'New Name', got '%s'", testimonial.Name)
			}
			if *testimonial.Organization != "New Org" {
				t.Errorf("Expected Organization 'New Org', got '%s'", *testimonial.Organization)
			}
			if *testimonial.Position != "New Pos" {
				t.Errorf("Expected Position 'New Pos', got '%s'", *testimonial.Position)
			}
			if testimonial.Content != "New Content" {
				t.Errorf("Expected Content 'New Content', got '%s'", testimonial.Content)
			}
			if testimonial.IsActive {
				t.Error("Expected IsActive to be false")
			}
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	isActive := false
	req := requests.UpdateTestimonialRequest{
		Name:         "New Name",
		Organization: "New Org",
		Position:     "New Pos",
		Content:      "New Content",
		IsActive:     &isActive,
	}

	result, err := service.Update(context.Background(), 1, req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}
