package service

import (
	"context"
	"errors"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// MockPublicMemberRepository adalah mock untuk MemberRepository (public methods)
type MockPublicMemberRepository struct {
	CreateFunc                   func(member *domain.Member) error
	FindAllFunc                  func(page, limit int, search string) ([]domain.Member, int64, error)
	FindByIDFunc                 func(id int) (*domain.Member, error)
	UpdateFunc                   func(member *domain.Member) error
	DeleteFunc                   func(id int) error
	FindActiveWithPaginationFunc func(page, limit int, search string) ([]domain.Member, int64, error)
	FindActiveByDepartmentFunc   func(department string, page, limit int, search string) ([]domain.Member, int64, error)
}

func (m *MockPublicMemberRepository) Create(member *domain.Member) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(member)
	}
	return errors.New("mock not configured")
}

func (m *MockPublicMemberRepository) FindAll(page, limit int, search string) ([]domain.Member, int64, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(page, limit, search)
	}
	return nil, 0, errors.New("mock not configured")
}

func (m *MockPublicMemberRepository) FindByID(id int) (*domain.Member, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockPublicMemberRepository) Update(member *domain.Member) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(member)
	}
	return errors.New("mock not configured")
}

func (m *MockPublicMemberRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return errors.New("mock not configured")
}

func (m *MockPublicMemberRepository) FindActiveWithPagination(page, limit int, search string) ([]domain.Member, int64, error) {
	if m.FindActiveWithPaginationFunc != nil {
		return m.FindActiveWithPaginationFunc(page, limit, search)
	}
	return nil, 0, errors.New("mock not configured")
}

func (m *MockPublicMemberRepository) FindActiveByDepartment(department string, page, limit int, search string) ([]domain.Member, int64, error) {
	if m.FindActiveByDepartmentFunc != nil {
		return m.FindActiveByDepartmentFunc(department, page, limit, search)
	}
	return nil, 0, errors.New("mock not configured")
}

// MockContactRepository adalah mock untuk ContactRepository
type MockContactRepository struct {
	GetFunc    func() (*domain.Contact, error)
	UpdateFunc func(contact *domain.Contact) error
}

func (m *MockContactRepository) Get() (*domain.Contact, error) {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	return nil, errors.New("mock not configured")
}

func (m *MockContactRepository) Update(contact *domain.Contact) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(contact)
	}
	return errors.New("mock not configured")
}

// ==================== GET ABOUT PAGE TESTS ====================

// Test: GetAboutPage berhasil dengan data lengkap
func TestPublicAboutGetAboutPage_Success(t *testing.T) {
	history := "Sejarah PMII"
	vision := "Visi PMII"
	mission := "Misi PMII"
	imageURI := "about.jpg"
	photoURI := "member.jpg"
	address := "Jl. Test No. 123"
	email := "contact@pmii.id"

	mockAboutRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{
				ID:       1,
				History:  &history,
				Vision:   &vision,
				Mission:  &mission,
				ImageURI: &imageURI,
			}, nil
		},
	}

	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			dept := domain.MemberDepartment(department)
			return []domain.Member{
				{ID: 1, FullName: "John Doe", Position: "Ketua", Department: dept, PhotoURI: &photoURI, IsActive: true},
				{ID: 2, FullName: "Jane Doe", Position: "Sekretaris", Department: dept, IsActive: true},
			}, 2, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/" + folder + "/" + filename
		},
	}

	mockContactRepo := &MockContactRepository{
		GetFunc: func() (*domain.Contact, error) {
			return &domain.Contact{
				ID:      1,
				Address: &address,
				Email:   &email,
			}, nil
		},
	}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, mockContactRepo, mockCloudinary)
	result, err := service.GetAboutPage(context.Background(), 8, "")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Check about data
	if *result.About.History != history {
		t.Errorf("Expected History '%s', got: '%s'", history, *result.About.History)
	}
	if result.About.ImageUrl != "https://cloudinary.com/about/about.jpg" {
		t.Errorf("Expected ImageUrl, got: '%s'", result.About.ImageUrl)
	}

	// Check departments data (should have 4 departments)
	if len(result.Departments) != 4 {
		t.Errorf("Expected 4 departments, got: %d", len(result.Departments))
	}

	// Check first department has members
	if len(result.Departments[0].Members.Data) != 2 {
		t.Errorf("Expected 2 members in first department, got: %d", len(result.Departments[0].Members.Data))
	}
}

// Test: GetAboutPage ketika about belum ada data
func TestPublicAboutGetAboutPage_NoAboutData(t *testing.T) {
	mockAboutRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return nil, errors.New("record not found")
		},
	}

	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			return []domain.Member{}, 0, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	result, err := service.GetAboutPage(context.Background(), 8, "")

	if err != nil {
		t.Errorf("Expected no error for empty about, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	// About should be empty but not error
	if result.About.History != nil {
		t.Error("Expected nil History for empty about")
	}
}

// Test: GetAboutPage dengan search
func TestPublicAboutGetAboutPage_WithSearch(t *testing.T) {
	searchQuery := "Ketua"

	mockAboutRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{ID: 1}, nil
		},
	}

	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			if search != searchQuery {
				t.Errorf("Expected search '%s', got: '%s'", searchQuery, search)
			}
			return []domain.Member{
				{ID: 1, FullName: "John Doe", Position: "Ketua Bidang", Department: domain.MemberDepartment(department), IsActive: true},
			}, 1, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	result, err := service.GetAboutPage(context.Background(), 8, searchQuery)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(result.Departments) != 4 {
		t.Errorf("Expected 4 departments, got: %d", len(result.Departments))
	}
}

// Test: GetAboutPage error dari member repository - should continue gracefully
func TestPublicAboutGetAboutPage_MemberRepoError(t *testing.T) {
	mockAboutRepo := &MockAboutRepository{
		GetFunc: func() (*domain.About, error) {
			return &domain.About{ID: 1}, nil
		},
	}

	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			return nil, 0, errors.New("database error")
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	result, err := service.GetAboutPage(context.Background(), 8, "")

	// Should NOT error - continues gracefully with empty departments
	if err != nil {
		t.Errorf("Expected no error (graceful degradation), got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	// Departments should be empty due to errors
	if len(result.Departments) != 0 {
		t.Errorf("Expected 0 departments (all errored), got: %d", len(result.Departments))
	}
}

// ==================== GET MEMBERS BY DEPARTMENT TESTS ====================

// Test: GetMembersByDepartment berhasil
func TestPublicAboutGetMembersByDepartment_Success(t *testing.T) {
	photoURI := "photo.jpg"
	dept := domain.DepartmentPengurusHarian

	mockAboutRepo := &MockAboutRepository{}
	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			if department != string(dept) {
				t.Errorf("Expected department '%s', got: '%s'", dept, department)
			}
			return []domain.Member{
				{ID: 1, FullName: "Member 1", Position: "Ketua", Department: dept, PhotoURI: &photoURI, IsActive: true},
				{ID: 2, FullName: "Member 2", Position: "Sekretaris", Department: dept, IsActive: true},
			}, 2, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return "https://cloudinary.com/" + folder + "/" + filename
		},
	}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	members, page, lastPage, total, err := service.GetMembersByDepartment(context.Background(), string(dept), 1, 8, "")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("Expected 2 members, got: %d", len(members))
	}
	if page != 1 {
		t.Errorf("Expected page 1, got: %d", page)
	}
	if lastPage != 1 {
		t.Errorf("Expected lastPage 1, got: %d", lastPage)
	}
	if total != 2 {
		t.Errorf("Expected total 2, got: %d", total)
	}
	if members[0].Photo != "https://cloudinary.com/members/photo.jpg" {
		t.Errorf("Expected photo URL, got: %s", members[0].Photo)
	}
}

// Test: GetMembersByDepartment dengan search filter
func TestPublicAboutGetMembersByDepartment_WithSearch(t *testing.T) {
	dept := domain.DepartmentKabid

	mockAboutRepo := &MockAboutRepository{}
	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			if search != "Ketua" {
				t.Errorf("Expected search 'Ketua', got: '%s'", search)
			}
			return []domain.Member{
				{ID: 1, FullName: "John", Position: "Ketua Bidang", Department: dept, IsActive: true},
			}, 1, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	members, _, _, _, err := service.GetMembersByDepartment(context.Background(), string(dept), 1, 8, "Ketua")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("Expected 1 member, got: %d", len(members))
	}
}

// Test: GetMembersByDepartment default pagination
func TestPublicAboutGetMembersByDepartment_DefaultPagination(t *testing.T) {
	dept := domain.DepartmentWasekbid

	mockAboutRepo := &MockAboutRepository{}
	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			if page != 1 {
				t.Errorf("Expected default page 1, got: %d", page)
			}
			if limit != 8 {
				t.Errorf("Expected default limit 8, got: %d", limit)
			}
			return []domain.Member{}, 0, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	_, page, lastPage, _, err := service.GetMembersByDepartment(context.Background(), string(dept), 0, 0, "") // Invalid values should use defaults

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if page != 1 {
		t.Errorf("Expected page 1, got: %d", page)
	}
	if lastPage != 1 {
		t.Errorf("Expected lastPage 1 for empty, got: %d", lastPage)
	}
}

// Test: GetMembersByDepartment error
func TestPublicAboutGetMembersByDepartment_Error(t *testing.T) {
	dept := domain.DepartmentWakilBendahara

	mockAboutRepo := &MockAboutRepository{}
	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			return nil, 0, errors.New("database error")
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	_, _, _, _, err := service.GetMembersByDepartment(context.Background(), string(dept), 1, 8, "")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test: GetMembersByDepartment pagination calculation dengan banyak data
func TestPublicAboutGetMembersByDepartment_PaginationCalculation(t *testing.T) {
	dept := domain.DepartmentPengurusHarian

	mockAboutRepo := &MockAboutRepository{}
	mockMemberRepo := &MockPublicMemberRepository{
		FindActiveByDepartmentFunc: func(department string, page, limit int, search string) ([]domain.Member, int64, error) {
			// Simulate 25 total members, requesting page 3 with limit 8
			return []domain.Member{
				{ID: 17, FullName: "Member 17", Position: "Position", Department: dept, IsActive: true},
			}, 25, nil
		},
	}

	mockCloudinary := &MockAboutCloudinaryService{
		GetImageURLFunc: func(folder string, filename string) string {
			return ""
		},
	}

	service := NewPublicAboutService(mockAboutRepo, mockMemberRepo, &MockContactRepository{}, mockCloudinary)
	_, page, lastPage, total, err := service.GetMembersByDepartment(context.Background(), string(dept), 3, 8, "")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if page != 3 {
		t.Errorf("Expected page 3, got: %d", page)
	}
	if lastPage != 4 { // 25/8 = 3.125, rounded up = 4
		t.Errorf("Expected lastPage 4, got: %d", lastPage)
	}
	if total != 25 {
		t.Errorf("Expected total 25, got: %d", total)
	}
}
