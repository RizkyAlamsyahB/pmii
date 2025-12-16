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

// MockAboutRepository adalah mock untuk AboutRepository
type MockAboutRepository struct {
	GetFunc    func() (*domain.About, error)
	UpsertFunc func(about *domain.About) error
}

func (m *MockAboutRepository) Get() (*domain.About, error) {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockAboutRepository) Upsert(about *domain.About) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(about)
	}
	return errors.New("mock not configured")
}

// MockAboutCloudinaryService adalah mock untuk CloudinaryService (untuk about tests)
type MockAboutCloudinaryService struct {
	UploadImageFunc func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteImageFunc func(ctx context.Context, folder string, filename string) error
	GetImageURLFunc func(folder string, filename string) string
}

func (m *MockAboutCloudinaryService) UploadImage(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadImageFunc != nil {
		return m.UploadImageFunc(ctx, folder, file)
	}
	return "", errors.New("mock not configured")
}

func (m *MockAboutCloudinaryService) DeleteImage(ctx context.Context, folder string, filename string) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, folder, filename)
	}
	return errors.New("mock not configured")
}

func (m *MockAboutCloudinaryService) GetImageURL(folder string, filename string) string {
	if m.GetImageURLFunc != nil {
		return m.GetImageURLFunc(folder, filename)
	}
	return ""
}

// ==================== GET TESTS ====================

// Test: Get berhasil dengan data existing
func TestAboutGet_Success(t *testing.T) {
	history := "Sejarah PMII"
	vision := "Visi PMII"
	mission := "Misi PMII"
	imageURI := "about123.jpg"
	videoURL := "https://youtube.com/watch?v=123"

	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{
				ID:        1,
				History:   &history,
				Vision:    &vision,
				Mission:   &mission,
				ImageURI:  &imageURI,
				VideoURL:  &videoURL,
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/" + folder + "/" + filename
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	result, err := service.Get(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.ID != 1 {
		t.Errorf("Expected ID 1, got: %d", result.ID)
	}
	if *result.History != history {
		t.Errorf("Expected History '%s', got: '%s'", history, *result.History)
	}
	if *result.Vision != vision {
		t.Errorf("Expected Vision '%s', got: '%s'", vision, *result.Vision)
	}
	if *result.Mission != mission {
		t.Errorf("Expected Mission '%s', got: '%s'", mission, *result.Mission)
	}
	if result.ImageUrl != "https://cloudinary.com/about/about123.jpg" {
		t.Errorf("Expected ImageUrl, got: '%s'", result.ImageUrl)
	}
}

// Test: Get ketika belum ada data (return empty response)
func TestAboutGet_NoData(t *testing.T) {
	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{}

	service := NewAboutService(mockRepo, mockCloudinary)
	result, err := service.Get(context.Background())

	if err != nil {
		t.Errorf("Expected no error for empty data, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected empty result, got nil")
	}
	if result.ID != 0 {
		t.Errorf("Expected ID 0 for empty, got: %d", result.ID)
	}
}

// ==================== UPDATE TESTS ====================

// Test: Update berhasil tanpa upload image
func TestAboutUpdate_SuccessWithoutImage(t *testing.T) {
	history := "Sejarah Lama"
	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{
				ID:      1,
				History: &history,
			}, nil
		},
		UpsertFunc: func(about *domain.About) error {
			return nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{
		History: "Sejarah Baru",
		Vision:  "Visi Baru",
	}

	result, err := service.Update(context.Background(), req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if *result.History != "Sejarah Baru" {
		t.Errorf("Expected History 'Sejarah Baru', got: '%s'", *result.History)
	}
	if *result.Vision != "Visi Baru" {
		t.Errorf("Expected Vision 'Visi Baru', got: '%s'", *result.Vision)
	}
}

// Test: Update berhasil dengan upload image baru
func TestAboutUpdate_SuccessWithImage(t *testing.T) {
	oldImage := "old-image.jpg"
	deleteCalled := false

	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{
				ID:       1,
				ImageURI: &oldImage,
			}, nil
		},
		UpsertFunc: func(about *domain.About) error {
			return nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			if folder != "about" {
				t.Errorf("Expected folder 'about', got: %s", folder)
			}
			return "new-image.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			deleteCalled = true
			if filename != "old-image.jpg" {
				t.Errorf("Expected to delete old image 'old-image.jpg', got: %s", filename)
			}
			return nil
		},
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/" + folder + "/" + filename
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{History: "Updated History"}
	mockFile := &multipart.FileHeader{Filename: "new-photo.jpg"}

	result, err := service.Update(context.Background(), req, mockFile)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if !deleteCalled {
		t.Error("Expected old image to be deleted after successful update")
	}
	if result.ImageUrl != "https://cloudinary.com/about/new-image.jpg" {
		t.Errorf("Expected new image URL, got: %s", result.ImageUrl)
	}
}

// Test: Update ketika belum ada data (create baru)
func TestAboutUpdate_CreateNew(t *testing.T) {
	upsertCalled := false

	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return nil, errors.New("record not found")
		},
		UpsertFunc: func(about *domain.About) error {
			upsertCalled = true
			if about.History == nil || *about.History != "Sejarah Baru" {
				t.Error("Expected History to be set")
			}
			return nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{
		History: "Sejarah Baru",
	}

	_, err := service.Update(context.Background(), req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !upsertCalled {
		t.Error("Expected Upsert to be called")
	}
}

// Test: Upload error harus return error
func TestAboutUpdate_ErrorUploadFailed(t *testing.T) {
	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{ID: 1}, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{History: "Test"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Update(context.Background(), req, mockFile)

	if err == nil || err.Error() != "gagal mengupload gambar" {
		t.Errorf("Expected 'gagal mengupload gambar' error, got: %v", err)
	}
}

// Test: Database error harus rollback (hapus image yang sudah diupload)
func TestAboutUpdate_ErrorDatabaseWithRollback(t *testing.T) {
	deleteCalled := false

	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{ID: 1}, nil
		},
		UpsertFunc: func(about *domain.About) error {
			return errors.New("database error")
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "new-image.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			deleteCalled = true
			if filename != "new-image.jpg" {
				t.Errorf("Expected to rollback 'new-image.jpg', got: %s", filename)
			}
			return nil
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{History: "Test"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Update(context.Background(), req, mockFile)

	if err == nil || err.Error() != "gagal menyimpan about" {
		t.Errorf("Expected 'gagal menyimpan about' error, got: %v", err)
	}
	if !deleteCalled {
		t.Error("Expected rollback: image baru harus dihapus ketika database error")
	}
}

// Test: Update dengan semua fields
func TestAboutUpdate_AllFields(t *testing.T) {
	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{ID: 1}, nil
		},
		UpsertFunc: func(about *domain.About) error {
			// Verify semua field di-update
			if about.History == nil || *about.History != "History" {
				t.Error("Expected History to be updated")
			}
			if about.Vision == nil || *about.Vision != "Vision" {
				t.Error("Expected Vision to be updated")
			}
			if about.Mission == nil || *about.Mission != "Mission" {
				t.Error("Expected Mission to be updated")
			}
			if about.VideoURL == nil || *about.VideoURL != "https://youtube.com/test" {
				t.Error("Expected VideoURL to be updated")
			}
			return nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{
		History:  "History",
		Vision:   "Vision",
		Mission:  "Mission",
		VideoURL: "https://youtube.com/test",
	}

	result, err := service.Update(context.Background(), req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

// Test: Update partial fields (hanya update yang dikirim)
func TestAboutUpdate_PartialFields(t *testing.T) {
	existingHistory := "Existing History"
	existingVision := "Existing Vision"

	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{
				ID:      1,
				History: &existingHistory,
				Vision:  &existingVision,
			}, nil
		},
		UpsertFunc: func(about *domain.About) error {
			// History harus tetap karena tidak dikirim (empty string)
			// Vision harus update karena dikirim
			if *about.Vision != "New Vision" {
				t.Errorf("Expected Vision to be updated to 'New Vision', got: %s", *about.Vision)
			}
			return nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{
		Vision: "New Vision",
		// History tidak dikirim (empty)
	}

	_, err := service.Update(context.Background(), req, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test: Update tanpa image lama (tidak perlu delete)
func TestAboutUpdate_NoOldImage(t *testing.T) {
	deleteCalled := false

	mockRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{
				ID:       1,
				ImageURI: nil, // Tidak ada image lama
			}, nil
		},
		UpsertFunc: func(about *domain.About) error {
			return nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "new-image.jpg", nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, filename string) error {
			deleteCalled = true
			return nil
		},
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/" + folder + "/" + filename
		},
	}

	service := NewAboutService(mockRepo, mockCloudinary)
	req := requests.UpdateAboutRequest{History: "Test"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Update(context.Background(), req, mockFile)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if deleteCalled {
		t.Error("Expected no delete call when there's no old image")
	}
}
