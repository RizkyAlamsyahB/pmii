package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/pkg/logger"
)

// init menginisialisasi logger sebelum test dijalankan
func init() {
	logger.Init()
}

// ============================================================
// Mock Repositories untuk HomeService
// ============================================================

// MockHomeRepository adalah mock untuk HomeRepository
type MockHomeRepository struct {
	GetHeroSectionFunc       func() ([]responses.HeroSectionResponse, error)
	GetLatestNewsSectionFunc func() ([]responses.LatestNewsSectionResponse, error)
	GetAboutUsSectionFunc    func() (*responses.AboutUsSectionResponse, error)
	GetWhySectionFunc        func() (*responses.WhySectionResponse, error)
	GetFaqSectionFunc        func() (*responses.FaqSectionResponse, error)
	GetCtaSectionFunc        func() (*responses.CtaSectionResponse, error)
}

func (m *MockHomeRepository) GetHeroSection() ([]responses.HeroSectionResponse, error) {
	if m.GetHeroSectionFunc != nil {
		return m.GetHeroSectionFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockHomeRepository) GetLatestNewsSection() ([]responses.LatestNewsSectionResponse, error) {
	if m.GetLatestNewsSectionFunc != nil {
		return m.GetLatestNewsSectionFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockHomeRepository) GetAboutUsSection() (*responses.AboutUsSectionResponse, error) {
	if m.GetAboutUsSectionFunc != nil {
		return m.GetAboutUsSectionFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockHomeRepository) GetWhySection() (*responses.WhySectionResponse, error) {
	if m.GetWhySectionFunc != nil {
		return m.GetWhySectionFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockHomeRepository) GetFaqSection() (*responses.FaqSectionResponse, error) {
	if m.GetFaqSectionFunc != nil {
		return m.GetFaqSectionFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockHomeRepository) GetCtaSection() (*responses.CtaSectionResponse, error) {
	if m.GetCtaSectionFunc != nil {
		return m.GetCtaSectionFunc()
	}
	return nil, errors.New("mock not configured")
}

// MockTestimonialRepository adalah mock untuk TestimonialRepository (untuk HomeService)
type MockTestimonialRepositoryForHomeService struct {
	CreateFunc   func(testimonial *domain.Testimonial) error
	FindAllFunc  func(page, limit int, search string) ([]domain.Testimonial, int64, error)
	FindByIDFunc func(id int) (*domain.Testimonial, error)
	UpdateFunc   func(testimonial *domain.Testimonial) error
	DeleteFunc   func(id int) error
}

func (m *MockTestimonialRepositoryForHomeService) Create(testimonial *domain.Testimonial) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(testimonial)
	}
	return nil
}

func (m *MockTestimonialRepositoryForHomeService) FindAll(page, limit int, search string) ([]domain.Testimonial, int64, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(page, limit, search)
	}
	return nil, 0, errors.New("mock not configured")
}

func (m *MockTestimonialRepositoryForHomeService) FindByID(id int) (*domain.Testimonial, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockTestimonialRepositoryForHomeService) Update(testimonial *domain.Testimonial) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(testimonial)
	}
	return nil
}

func (m *MockTestimonialRepositoryForHomeService) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockCloudinaryServiceForHomeService adalah mock untuk CloudinaryService
type MockCloudinaryServiceForHomeService struct {
	GetImageURLFunc    func(folder string, fileName string) string
	UploadImageFunc    func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteImageFunc    func(ctx context.Context, folder string, fileName string) error
	UploadFileFunc     func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteFileFunc     func(ctx context.Context, folder string, filename string) error
	GetFileURLFunc     func(folder string, filename string) string
	GetDownloadURLFunc func(folder string, filename string) string
}

func (m *MockCloudinaryServiceForHomeService) GetImageURL(folder string, fileName string) string {
	if m.GetImageURLFunc != nil {
		return m.GetImageURLFunc(folder, fileName)
	}
	return "https://res.cloudinary.com/test/" + folder + "/" + fileName
}

func (m *MockCloudinaryServiceForHomeService) UploadImage(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadImageFunc != nil {
		return m.UploadImageFunc(ctx, folder, file)
	}
	return "", nil
}

func (m *MockCloudinaryServiceForHomeService) DeleteImage(ctx context.Context, folder string, fileName string) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, folder, fileName)
	}
	return nil
}

func (m *MockCloudinaryServiceForHomeService) UploadFile(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(ctx, folder, file)
	}
	return "", nil
}

func (m *MockCloudinaryServiceForHomeService) DeleteFile(ctx context.Context, folder string, filename string) error {
	if m.DeleteFileFunc != nil {
		return m.DeleteFileFunc(ctx, folder, filename)
	}
	return nil
}

func (m *MockCloudinaryServiceForHomeService) GetFileURL(folder string, filename string) string {
	if m.GetFileURLFunc != nil {
		return m.GetFileURLFunc(folder, filename)
	}
	return ""
}

func (m *MockCloudinaryServiceForHomeService) GetDownloadURL(folder string, filename string) string {
	if m.GetDownloadURLFunc != nil {
		return m.GetDownloadURLFunc(folder, filename)
	}
	return ""
}

// ============================================================
// Test Cases untuk GetHeroSection
// ============================================================

func TestGetHeroSection(t *testing.T) {
	mockHeroSections := []responses.HeroSectionResponse{
		{
			ID:       1,
			Title:    "Post 1",
			Slug:     "post-1",
			ImageURL: "image1.jpg",
		},
		{
			ID:       2,
			Title:    "Post 2",
			Slug:     "post-2",
			ImageURL: "image2.jpg",
		},
	}

	tests := []struct {
		name           string
		mockHeroData   []responses.HeroSectionResponse
		mockErr        error
		expectedLen    int
		expectedErr    bool
		validateImages bool
	}{
		{
			name:           "success dengan data",
			mockHeroData:   mockHeroSections,
			expectedLen:    2,
			validateImages: true,
		},
		{
			name:         "empty list",
			mockHeroData: []responses.HeroSectionResponse{},
			expectedLen:  0,
		},
		{
			name:        "repository error",
			mockErr:     errors.New("database connection failed"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHomeRepo := &MockHomeRepository{
				GetHeroSectionFunc: func() ([]responses.HeroSectionResponse, error) {
					return tt.mockHeroData, tt.mockErr
				},
			}
			mockCloudinary := &MockCloudinaryServiceForHomeService{
				GetImageURLFunc: func(folder, fileName string) string {
					return "https://res.cloudinary.com/test/" + folder + "/" + fileName
				},
			}

			service := NewPublicHomeService(mockHomeRepo, &MockTestimonialRepositoryForHomeService{}, mockCloudinary)
			heroes, err := service.GetHeroSection()

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(heroes) != tt.expectedLen {
				t.Errorf("expected %d heroes, got %d", tt.expectedLen, len(heroes))
			}
			if tt.validateImages && len(heroes) > 0 {
				for _, hero := range heroes {
					if hero.ImageURL == "" {
						t.Error("expected image URL to be transformed")
					}
				}
			}
		})
	}
}

// ============================================================
// Test Cases untuk GetLatestNewsSection
// ============================================================

func TestGetLatestNewsSection(t *testing.T) {
	mockLatestNews := []responses.LatestNewsSectionResponse{
		{ID: 1, Title: "News 1", FeaturedImage: "news1.jpg"},
		{ID: 2, Title: "News 2", FeaturedImage: "news2.jpg"},
	}

	tests := []struct {
		name        string
		mockData    []responses.LatestNewsSectionResponse
		mockErr     error
		expectedLen int
		expectedErr bool
	}{
		{
			name:        "success dengan data",
			mockData:    mockLatestNews,
			expectedLen: 2,
		},
		{
			name:        "empty list",
			mockData:    []responses.LatestNewsSectionResponse{},
			expectedLen: 0,
		},
		{
			name:        "repository error",
			mockErr:     errors.New("database connection failed"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHomeRepo := &MockHomeRepository{
				GetLatestNewsSectionFunc: func() ([]responses.LatestNewsSectionResponse, error) {
					return tt.mockData, tt.mockErr
				},
			}
			mockCloudinary := &MockCloudinaryServiceForHomeService{
				GetImageURLFunc: func(folder, fileName string) string {
					return "https://res.cloudinary.com/test/" + folder + "/" + fileName
				},
			}

			service := NewPublicHomeService(mockHomeRepo, &MockTestimonialRepositoryForHomeService{}, mockCloudinary)
			news, err := service.GetLatestNewsSection()

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(news) != tt.expectedLen {
				t.Errorf("expected %d news, got %d", tt.expectedLen, len(news))
			}
		})
	}
}

// ============================================================
// Test Cases untuk GetAboutUsSection
// ============================================================

func TestGetAboutUsSection(t *testing.T) {
	mockAboutUs := &responses.AboutUsSectionResponse{
		Title:       "Tentang PMII",
		Subtitle:    "Subtitle",
		Description: "Description",
		ImageURI:    "about.jpg",
	}

	tests := []struct {
		name        string
		mockData    *responses.AboutUsSectionResponse
		mockErr     error
		expectedErr bool
	}{
		{
			name:     "success",
			mockData: mockAboutUs,
		},
		{
			name:        "repository error",
			mockErr:     errors.New("database connection failed"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHomeRepo := &MockHomeRepository{
				GetAboutUsSectionFunc: func() (*responses.AboutUsSectionResponse, error) {
					return tt.mockData, tt.mockErr
				},
			}
			mockCloudinary := &MockCloudinaryServiceForHomeService{
				GetImageURLFunc: func(folder, fileName string) string {
					return "https://res.cloudinary.com/test/" + folder + "/" + fileName
				},
			}

			service := NewPublicHomeService(mockHomeRepo, &MockTestimonialRepositoryForHomeService{}, mockCloudinary)
			aboutUs, err := service.GetAboutUsSection()

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if aboutUs.Title != tt.mockData.Title {
				t.Errorf("expected title %s, got %s", tt.mockData.Title, aboutUs.Title)
			}
		})
	}
}

// ============================================================
// Test Cases untuk GetWhySection
// ============================================================

func TestGetWhySection(t *testing.T) {
	mockWhySection := &responses.WhySectionResponse{
		Title:    "Mengapa PMII?",
		Subtitle: "Subtitle",
		Data: []responses.WhyItem{
			{Title: "Item 1", Description: "Desc 1", IconURI: "icon1.png"},
			{Title: "Item 2", Description: "Desc 2", IconURI: "icon2.png"},
		},
	}

	tests := []struct {
		name        string
		mockData    *responses.WhySectionResponse
		mockErr     error
		expectedErr bool
	}{
		{
			name:     "success",
			mockData: mockWhySection,
		},
		{
			name:        "repository error",
			mockErr:     errors.New("database connection failed"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHomeRepo := &MockHomeRepository{
				GetWhySectionFunc: func() (*responses.WhySectionResponse, error) {
					return tt.mockData, tt.mockErr
				},
			}
			mockCloudinary := &MockCloudinaryServiceForHomeService{
				GetImageURLFunc: func(folder, fileName string) string {
					return "https://res.cloudinary.com/test/" + folder + "/" + fileName
				},
			}

			service := NewPublicHomeService(mockHomeRepo, &MockTestimonialRepositoryForHomeService{}, mockCloudinary)
			whySection, err := service.GetWhySection()

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if whySection.Title != tt.mockData.Title {
				t.Errorf("expected title %s, got %s", tt.mockData.Title, whySection.Title)
			}
			if len(whySection.Data) != len(tt.mockData.Data) {
				t.Errorf("expected %d items, got %d", len(tt.mockData.Data), len(whySection.Data))
			}
		})
	}
}

// ============================================================
// Test Cases untuk GetTestimonialSection
// ============================================================

func TestGetTestimonialSection(t *testing.T) {
	position := "Ketua"
	organization := "PMII Cabang Jakarta"
	photoURI := "photo.jpg"

	mockTestimonials := []domain.Testimonial{
		{ID: 1, Name: "John", Content: "Great!", Position: &position, Organization: &organization, PhotoURI: &photoURI},
		{ID: 2, Name: "Jane", Content: "Amazing!", Position: nil, Organization: nil, PhotoURI: nil},
	}

	tests := []struct {
		name          string
		mockData      []domain.Testimonial
		mockTotal     int64
		mockErr       error
		expectedLen   int
		expectedErr   bool
		checkNilField bool
	}{
		{
			name:          "success dengan data lengkap dan nil fields",
			mockData:      mockTestimonials,
			mockTotal:     2,
			expectedLen:   2,
			checkNilField: true,
		},
		{
			name:        "empty list",
			mockData:    []domain.Testimonial{},
			mockTotal:   0,
			expectedLen: 0,
		},
		{
			name:        "repository error",
			mockErr:     errors.New("database connection failed"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTestimonialRepo := &MockTestimonialRepositoryForHomeService{
				FindAllFunc: func(page, limit int, search string) ([]domain.Testimonial, int64, error) {
					return tt.mockData, tt.mockTotal, tt.mockErr
				},
			}
			mockCloudinary := &MockCloudinaryServiceForHomeService{
				GetImageURLFunc: func(folder, fileName string) string {
					return "https://res.cloudinary.com/test/" + folder + "/" + fileName
				},
			}

			service := NewPublicHomeService(&MockHomeRepository{}, mockTestimonialRepo, mockCloudinary)
			testimonials, err := service.GetTestimonialSection()

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(testimonials) != tt.expectedLen {
				t.Errorf("expected %d testimonials, got %d", tt.expectedLen, len(testimonials))
			}
			// Check nil field handling
			if tt.checkNilField && len(testimonials) >= 2 {
				// First testimonial should have values
				if testimonials[0].Status != position {
					t.Errorf("expected status %s, got %s", position, testimonials[0].Status)
				}
				// Second testimonial should have empty strings for nil fields
				if testimonials[1].Status != "" {
					t.Errorf("expected empty status for nil field, got %s", testimonials[1].Status)
				}
				if testimonials[1].Career != "" {
					t.Errorf("expected empty career for nil field, got %s", testimonials[1].Career)
				}
			}
		})
	}
}

// ============================================================
// Test Cases untuk GetFaqSection
// ============================================================

func TestGetFaqSection(t *testing.T) {
	mockFaqSection := &responses.FaqSectionResponse{
		Title:    "FAQ",
		Subtitle: "Pertanyaan Umum",
		Data: []responses.FaqItem{
			{Question: "Q1?", Answer: "A1"},
			{Question: "Q2?", Answer: "A2"},
		},
	}

	tests := []struct {
		name        string
		mockData    *responses.FaqSectionResponse
		mockErr     error
		expectedErr bool
	}{
		{
			name:     "success",
			mockData: mockFaqSection,
		},
		{
			name:        "repository error",
			mockErr:     errors.New("database connection failed"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHomeRepo := &MockHomeRepository{
				GetFaqSectionFunc: func() (*responses.FaqSectionResponse, error) {
					return tt.mockData, tt.mockErr
				},
			}

			service := NewPublicHomeService(mockHomeRepo, &MockTestimonialRepositoryForHomeService{}, &MockCloudinaryServiceForHomeService{})
			faqSection, err := service.GetFaqSection()

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if faqSection.Title != tt.mockData.Title {
				t.Errorf("expected title %s, got %s", tt.mockData.Title, faqSection.Title)
			}
			if len(faqSection.Data) != len(tt.mockData.Data) {
				t.Errorf("expected %d items, got %d", len(tt.mockData.Data), len(faqSection.Data))
			}
		})
	}
}

// ============================================================
// Test Cases untuk GetCtaSection
// ============================================================

func TestGetCtaSection(t *testing.T) {
	mockCtaSection := &responses.CtaSectionResponse{
		Title:    "Siap Bergabung?",
		Subtitle: "CTA Subtitle",
	}

	tests := []struct {
		name        string
		mockData    *responses.CtaSectionResponse
		mockErr     error
		expectedErr bool
	}{
		{
			name:     "success",
			mockData: mockCtaSection,
		},
		{
			name:        "repository error",
			mockErr:     errors.New("database connection failed"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHomeRepo := &MockHomeRepository{
				GetCtaSectionFunc: func() (*responses.CtaSectionResponse, error) {
					return tt.mockData, tt.mockErr
				},
			}

			service := NewPublicHomeService(mockHomeRepo, &MockTestimonialRepositoryForHomeService{}, &MockCloudinaryServiceForHomeService{})
			ctaSection, err := service.GetCtaSection()

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if ctaSection.Title != tt.mockData.Title {
				t.Errorf("expected title %s, got %s", tt.mockData.Title, ctaSection.Title)
			}
			if ctaSection.Subtitle != tt.mockData.Subtitle {
				t.Errorf("expected subtitle %s, got %s", tt.mockData.Subtitle, ctaSection.Subtitle)
			}
		})
	}
}
