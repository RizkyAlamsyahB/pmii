package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
)

// MockTestimonialRepository adalah mock untuk TestimonialRepository
type MockTestimonialRepository struct {
	CreateFunc   func(testimonial *domain.Testimonial) error
	FindAllFunc  func(page, limit int, search string) ([]domain.Testimonial, int64, error)
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

func (m *MockTestimonialRepository) FindAll(page, limit int, search string) ([]domain.Testimonial, int64, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(page, limit, search)
	}
	return nil, 0, errors.New("mock not configured")
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
	UploadImageFunc    func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteImageFunc    func(ctx context.Context, folder string, filename string) error
	GetImageURLFunc    func(folder string, filename string) string
	UploadFileFunc     func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteFileFunc     func(ctx context.Context, folder string, filename string) error
	GetFileURLFunc     func(folder string, filename string) string
	GetDownloadURLFunc func(folder string, filename string) string
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

func (m *MockCloudinaryService) UploadFile(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(ctx, folder, file)
	}
	return "", nil
}

func (m *MockCloudinaryService) DeleteFile(ctx context.Context, folder string, filename string) error {
	if m.DeleteFileFunc != nil {
		return m.DeleteFileFunc(ctx, folder, filename)
	}
	return nil
}

func (m *MockCloudinaryService) GetFileURL(folder string, filename string) string {
	if m.GetFileURLFunc != nil {
		return m.GetFileURLFunc(folder, filename)
	}
	return ""
}

func (m *MockCloudinaryService) GetDownloadURL(folder string, filename string) string {
	if m.GetDownloadURLFunc != nil {
		return m.GetDownloadURLFunc(folder, filename)
	}
	return ""
}

// ==================== CREATE TESTS ====================

// Test: Upload error harus return error
func TestCreate_ErrorUploadFailed(t *testing.T) {
	mockRepo := &MockTestimonialRepository{}
	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)
	req := requests.CreateTestimonialRequest{Name: "Test", Content: "Content"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Create(context.Background(), req, mockFile)

	if err == nil || err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected 'gagal mengupload foto' error, got: %v", err)
	}
}

// Test: Database error harus rollback (hapus foto yang sudah diupload)
func TestCreate_ErrorDatabaseWithRollback(t *testing.T) {
	deleteCalled := false

	mockRepo := &MockTestimonialRepository{
		CreateFunc: func(testimonial *domain.Testimonial) error {
			return errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "photo123.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			deleteCalled = true
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)
	req := requests.CreateTestimonialRequest{Name: "Test", Content: "Content"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Create(context.Background(), req, mockFile)

	if err == nil || err.Error() != "gagal menyimpan testimonial" {
		t.Errorf("Expected 'gagal menyimpan testimonial' error, got: %v", err)
	}
	if !deleteCalled {
		t.Error("Expected rollback: foto harus dihapus ketika database error")
	}
}

// Test: Optional fields (organization, position) harus jadi nil kalau kosong
func TestCreate_ValidationEmptyOptionalFields(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		CreateFunc: func(testimonial *domain.Testimonial) error {
			if testimonial.Organization != nil {
				t.Error("Expected Organization nil when empty")
			}
			if testimonial.Position != nil {
				t.Error("Expected Position nil when empty")
			}
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string { return "" },
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)
	req := requests.CreateTestimonialRequest{
		Name:         "Test",
		Content:      "Content",
		Organization: "", // empty
		Position:     "", // empty
	}

	_, err := service.Create(context.Background(), req, nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// ==================== GET ALL TESTS ====================

// Test: Pagination default values harus diset (page=1, limit=10)
func TestGetAll_ValidationDefaultPagination(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindAllFunc: func(page, limit int, search string) ([]domain.Testimonial, int64, error) {
			if page != 1 || limit != 10 {
				t.Errorf("Expected default page=1, limit=10, got page=%d, limit=%d", page, limit)
			}
			return []domain.Testimonial{}, 0, nil
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewTestimonialService(mockRepo, mockCloudinary)

	// Test dengan invalid values
	_, _, _, _, err := service.GetAll(context.Background(), 0, -5, "")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test: Database error harus return error
func TestGetAll_ErrorDatabaseFailed(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindAllFunc: func(page, limit int, search string) ([]domain.Testimonial, int64, error) {
			return nil, 0, errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewTestimonialService(mockRepo, mockCloudinary)

	_, _, _, _, err := service.GetAll(context.Background(), 1, 10, "")

	if err == nil || err.Error() != "gagal mengambil data testimonial" {
		t.Errorf("Expected 'gagal mengambil data testimonial' error, got: %v", err)
	}
}

// Test: Kalkulasi lastPage harus benar
func TestGetAll_LogicLastPageCalculation(t *testing.T) {
	tests := []struct {
		total        int64
		limit        int
		expectedLast int
	}{
		{20, 10, 2}, // 20/10 = 2
		{25, 10, 3}, // 25/10 = 2.5 -> 3
		{5, 10, 1},  // 5/10 = 0.5 -> 1
		{0, 10, 0},  // 0/10 = 0
		{21, 10, 3}, // 21/10 = 2.1 -> 3
	}

	for _, tt := range tests {
		mockRepo := &MockTestimonialRepository{
			FindAllFunc: func(page, limit int, search string) ([]domain.Testimonial, int64, error) {
				return []domain.Testimonial{}, tt.total, nil
			},
		}

		mockCloudinary := &MockCloudinaryService{}
		service := NewTestimonialService(mockRepo, mockCloudinary)

		_, _, lastPage, _, err := service.GetAll(context.Background(), 1, tt.limit, "")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if lastPage != tt.expectedLast {
			t.Errorf("Total=%d, Limit=%d: Expected lastPage=%d, got=%d", tt.total, tt.limit, tt.expectedLast, lastPage)
		}
	}
}

// ==================== GET BY ID TESTS ====================

// Test: Not found harus return error
func TestGetByID_ErrorNotFound(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewTestimonialService(mockRepo, mockCloudinary)

	_, err := service.GetByID(context.Background(), 999)

	if err == nil || err.Error() != "testimonial tidak ditemukan" {
		t.Errorf("Expected 'testimonial tidak ditemukan' error, got: %v", err)
	}
}

// ==================== UPDATE TESTS ====================

// Test: Update dengan foto baru harus hapus foto lama SETELAH database berhasil
func TestUpdate_LogicDeleteOldPhotoAfterSuccess(t *testing.T) {
	oldPhoto := "old.jpg"
	deleteOldPhotoCalled := false
	updateCalled := false

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{
				ID:       1,
				Name:     "Test",
				Content:  "Content",
				PhotoURI: &oldPhoto,
			}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			updateCalled = true
			if deleteOldPhotoCalled {
				t.Error("BUG: Foto lama dihapus SEBELUM update database!")
			}
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "new.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			if !updateCalled {
				t.Error("BUG: Foto lama dihapus SEBELUM update database!")
			}
			deleteOldPhotoCalled = true
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)
	req := requests.UpdateTestimonialRequest{Name: "Updated"}
	mockFile := &multipart.FileHeader{Filename: "new.jpg"}

	_, err := service.Update(context.Background(), 1, req, mockFile)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !deleteOldPhotoCalled {
		t.Error("Expected old photo to be deleted")
	}
}

// Test: Upload error harus return error
func TestUpdate_ErrorUploadFailed(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{ID: 1, Name: "Test", Content: "Content"}, nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)
	req := requests.UpdateTestimonialRequest{Name: "Updated"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Update(context.Background(), 1, req, mockFile)

	if err == nil || err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected 'gagal mengupload foto' error, got: %v", err)
	}
}

// Test: Database error harus rollback (hapus foto baru, JANGAN hapus foto lama)
func TestUpdate_ErrorDatabaseWithRollback(t *testing.T) {
	oldPhoto := "old.jpg"
	deleteNewPhotoCalled := false
	deleteOldPhotoCalled := false

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{ID: 1, Name: "Test", Content: "Content", PhotoURI: &oldPhoto}, nil
		},
		UpdateFunc: func(testimonial *domain.Testimonial) error {
			return errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "new.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			if filename == "new.jpg" {
				deleteNewPhotoCalled = true
			}
			if filename == "old.jpg" {
				deleteOldPhotoCalled = true
			}
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)
	req := requests.UpdateTestimonialRequest{Name: "Updated"}
	mockFile := &multipart.FileHeader{Filename: "new.jpg"}

	_, err := service.Update(context.Background(), 1, req, mockFile)

	if err == nil || err.Error() != "gagal mengupdate testimonial" {
		t.Errorf("Expected 'gagal mengupdate testimonial' error, got: %v", err)
	}
	if !deleteNewPhotoCalled {
		t.Error("Expected rollback: foto baru harus dihapus")
	}
	if deleteOldPhotoCalled {
		t.Error("BUG: Foto lama TIDAK BOLEH dihapus ketika database error!")
	}
}

// Test: Not found harus return error
func TestUpdate_ErrorNotFound(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewTestimonialService(mockRepo, mockCloudinary)

	_, err := service.Update(context.Background(), 999, requests.UpdateTestimonialRequest{}, nil)

	if err == nil || err.Error() != "testimonial tidak ditemukan" {
		t.Errorf("Expected 'testimonial tidak ditemukan' error, got: %v", err)
	}
}

// ==================== DELETE TESTS ====================

// Test: Delete harus hapus record database dulu, BARU hapus foto
func TestDelete_LogicDeleteDatabaseFirst(t *testing.T) {
	photo := "photo.jpg"
	databaseDeleteCalled := false

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{ID: 1, Name: "Test", Content: "Content", PhotoURI: &photo}, nil
		},
		DeleteFunc: func(id int) error {
			databaseDeleteCalled = true
			return nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			if !databaseDeleteCalled {
				t.Error("BUG: Foto dihapus SEBELUM database delete!")
			}
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test: Database error harus return error (foto TIDAK dihapus)
func TestDelete_ErrorDatabaseFailed(t *testing.T) {
	photo := "photo.jpg"
	photoDeleteCalled := false

	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return &domain.Testimonial{ID: 1, Name: "Test", Content: "Content", PhotoURI: &photo}, nil
		},
		DeleteFunc: func(id int) error {
			return errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			photoDeleteCalled = true
			return nil
		},
	}

	service := NewTestimonialService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 1)

	if err == nil || err.Error() != "gagal menghapus testimonial" {
		t.Errorf("Expected 'gagal menghapus testimonial' error, got: %v", err)
	}
	if photoDeleteCalled {
		t.Error("BUG: Foto TIDAK BOLEH dihapus ketika database error!")
	}
}

// Test: Not found harus return error
func TestDelete_ErrorNotFound(t *testing.T) {
	mockRepo := &MockTestimonialRepository{
		FindByIDFunc: func(id int) (*domain.Testimonial, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewTestimonialService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 999)

	if err == nil || err.Error() != "testimonial tidak ditemukan" {
		t.Errorf("Expected 'testimonial tidak ditemukan' error, got: %v", err)
	}
}
