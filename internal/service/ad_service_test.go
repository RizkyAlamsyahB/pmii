package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/stretchr/testify/assert"
)

// MockAdRepository adalah mock untuk AdRepository
type MockAdRepository struct {
	FindAllFunc    func() ([]domain.Ad, error)
	FindByIDFunc   func(id int) (*domain.Ad, error)
	FindByPageFunc func(page domain.AdPage) ([]domain.Ad, error)
	UpdateFunc     func(ad *domain.Ad) error
}

func (m *MockAdRepository) FindAll() ([]domain.Ad, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockAdRepository) FindByID(id int) (*domain.Ad, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockAdRepository) FindByPage(page domain.AdPage) ([]domain.Ad, error) {
	if m.FindByPageFunc != nil {
		return m.FindByPageFunc(page)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockAdRepository) FindByPageAndSlot(page domain.AdPage, slot int) (*domain.Ad, error) {
	return nil, errors.New("mock not configured")
}

func (m *MockAdRepository) Update(ad *domain.Ad) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ad)
	}
	return nil
}

func (m *MockAdRepository) UpdateImage(id int, imageURL string) error {
	return nil
}

// MockCloudinaryServiceForAd adalah mock untuk CloudinaryService
type MockCloudinaryServiceForAd struct {
	UploadImageFunc func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteImageFunc func(ctx context.Context, folder string, filename string) error
	GetImageURLFunc func(folder string, filename string) string
}

func (m *MockCloudinaryServiceForAd) UploadImage(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadImageFunc != nil {
		return m.UploadImageFunc(ctx, folder, file)
	}
	return "uploaded_image.jpg", nil
}

func (m *MockCloudinaryServiceForAd) DeleteImage(ctx context.Context, folder string, filename string) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, folder, filename)
	}
	return nil
}

func (m *MockCloudinaryServiceForAd) GetImageURL(folder string, filename string) string {
	if m.GetImageURLFunc != nil {
		return m.GetImageURLFunc(folder, filename)
	}
	return "https://example.com/" + folder + "/" + filename
}

func (m *MockCloudinaryServiceForAd) UploadFile(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	return "uploaded_file.pdf", nil
}

func (m *MockCloudinaryServiceForAd) DeleteFile(ctx context.Context, folder string, filename string) error {
	return nil
}

func (m *MockCloudinaryServiceForAd) GetFileURL(folder string, filename string) string {
	return "https://example.com/" + folder + "/" + filename
}

func (m *MockCloudinaryServiceForAd) GetDownloadURL(folder string, filename string) string {
	return "https://example.com/download/" + folder + "/" + filename
}

// MockActivityLogRepoForAd adalah mock untuk ActivityLogRepository
type MockActivityLogRepoForAd struct{}

func (m *MockActivityLogRepoForAd) Create(log *domain.ActivityLog) error {
	return nil
}

func (m *MockActivityLogRepoForAd) GetActivityLogs(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
	return nil, 0, nil
}

// TestGetAllAds_Success tests successful retrieval of all ads
func TestGetAllAds_Success(t *testing.T) {
	imageURL := "test_image.jpg"
	mockAds := []domain.Ad{
		{ID: 1, Page: domain.AdPageLanding, Slot: 1, ImageURL: &imageURL, Resolution: "728x90"},
		{ID: 2, Page: domain.AdPageLanding, Slot: 2, ImageURL: nil, Resolution: "16x9"},
		{ID: 3, Page: domain.AdPageNews, Slot: 1, ImageURL: &imageURL, Resolution: "16x9"},
	}

	mockRepo := &MockAdRepository{
		FindAllFunc: func() ([]domain.Ad, error) {
			return mockAds, nil
		},
	}

	adService := NewAdService(mockRepo, &MockCloudinaryServiceForAd{}, &MockActivityLogRepoForAd{})
	result, err := adService.GetAllAds(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result)) // 2 pages (landing, news)
}

// TestGetAdByID_Success tests successful retrieval of ad by ID
func TestGetAdByID_Success(t *testing.T) {
	imageURL := "test_image.jpg"
	mockAd := &domain.Ad{
		ID:         1,
		Page:       domain.AdPageLanding,
		Slot:       1,
		ImageURL:   &imageURL,
		Resolution: "728x90",
	}

	mockRepo := &MockAdRepository{
		FindByIDFunc: func(id int) (*domain.Ad, error) {
			if id == 1 {
				return mockAd, nil
			}
			return nil, errors.New("not found")
		},
	}

	adService := NewAdService(mockRepo, &MockCloudinaryServiceForAd{}, &MockActivityLogRepoForAd{})
	result, err := adService.GetAdByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ADS 1 Landing Page", result.SlotName)
}

// TestGetAdByID_NotFound tests ad not found scenario
func TestGetAdByID_NotFound(t *testing.T) {
	mockRepo := &MockAdRepository{
		FindByIDFunc: func(id int) (*domain.Ad, error) {
			return nil, errors.New("not found")
		},
	}

	adService := NewAdService(mockRepo, &MockCloudinaryServiceForAd{}, &MockActivityLogRepoForAd{})
	result, err := adService.GetAdByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ad tidak ditemukan", err.Error())
}

// TestGetAdsByPage_Success tests successful retrieval of ads by page
func TestGetAdsByPage_Success(t *testing.T) {
	imageURL := "test_image.jpg"
	mockAds := []domain.Ad{
		{ID: 1, Page: domain.AdPageLanding, Slot: 1, ImageURL: &imageURL, Resolution: "728x90"},
		{ID: 2, Page: domain.AdPageLanding, Slot: 2, ImageURL: nil, Resolution: "16x9"},
	}

	mockRepo := &MockAdRepository{
		FindByPageFunc: func(page domain.AdPage) ([]domain.Ad, error) {
			if page == domain.AdPageLanding {
				return mockAds, nil
			}
			return nil, errors.New("not found")
		},
	}

	adService := NewAdService(mockRepo, &MockCloudinaryServiceForAd{}, &MockActivityLogRepoForAd{})
	result, err := adService.GetAdsByPage(context.Background(), "landing")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
}

// TestGetAdsByPage_InvalidPage tests invalid page name
func TestGetAdsByPage_InvalidPage(t *testing.T) {
	mockRepo := &MockAdRepository{}

	adService := NewAdService(mockRepo, &MockCloudinaryServiceForAd{}, &MockActivityLogRepoForAd{})
	result, err := adService.GetAdsByPage(context.Background(), "invalid_page")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "halaman tidak valid", err.Error())
}

// TestUpdateAd_Success tests successful ad update
func TestUpdateAd_Success(t *testing.T) {
	imageURL := "test_image.jpg"
	mockAd := &domain.Ad{
		ID:         1,
		Page:       domain.AdPageLanding,
		Slot:       1,
		ImageURL:   &imageURL,
		Resolution: "728x90",
	}

	mockRepo := &MockAdRepository{
		FindByIDFunc: func(id int) (*domain.Ad, error) {
			return mockAd, nil
		},
		UpdateFunc: func(ad *domain.Ad) error {
			return nil
		},
	}

	adService := NewAdService(mockRepo, &MockCloudinaryServiceForAd{}, &MockActivityLogRepoForAd{})

	result, err := adService.UpdateAd(context.Background(), 1, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestDeleteAdImage_Success tests successful image deletion
func TestDeleteAdImage_Success(t *testing.T) {
	imageURL := "test_image.jpg"
	mockAd := &domain.Ad{
		ID:         1,
		Page:       domain.AdPageLanding,
		Slot:       1,
		ImageURL:   &imageURL,
		Resolution: "728x90",
	}

	mockRepo := &MockAdRepository{
		FindByIDFunc: func(id int) (*domain.Ad, error) {
			return mockAd, nil
		},
		UpdateFunc: func(ad *domain.Ad) error {
			return nil
		},
	}

	adService := NewAdService(mockRepo, &MockCloudinaryServiceForAd{}, &MockActivityLogRepoForAd{})
	result, err := adService.DeleteAdImage(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.ImageURL)
}

// TestIsValidAdPage tests page validation
func TestIsValidAdPage(t *testing.T) {
	testCases := []struct {
		page     string
		expected bool
	}{
		{"landing", true},
		{"news", true},
		{"opini", true},
		{"life_at_pmii", true},
		{"islamic", true},
		{"detail_article", true},
		{"invalid", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.page, func(t *testing.T) {
			result := domain.IsValidAdPage(tc.page)
			assert.Equal(t, tc.expected, result)
		})
	}
}
