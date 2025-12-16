package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
)

// MockMemberRepository adalah mock untuk MemberRepository
type MockMemberRepository struct {
	CreateFunc                   func(member *domain.Member) error
	FindAllFunc                  func(page, limit int) ([]domain.Member, int64, error)
	FindByIDFunc                 func(id int) (*domain.Member, error)
	UpdateFunc                   func(member *domain.Member) error
	DeleteFunc                   func(id int) error
	FindActiveWithPaginationFunc func(page, limit int, search string) ([]domain.Member, int64, error)
	FindActiveByDepartmentFunc   func(department string, page, limit int, search string) ([]domain.Member, int64, error)
}

func (m *MockMemberRepository) Create(member *domain.Member) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(member)
	}
	return errors.New("mock not configured")
}

func (m *MockMemberRepository) FindAll(page, limit int) ([]domain.Member, int64, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(page, limit)
	}
	return nil, 0, errors.New("mock not configured")
}

func (m *MockMemberRepository) FindByID(id int) (*domain.Member, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockMemberRepository) Update(member *domain.Member) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(member)
	}
	return errors.New("mock not configured")
}

func (m *MockMemberRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return errors.New("mock not configured")
}

func (m *MockMemberRepository) FindActiveWithPagination(page, limit int, search string) ([]domain.Member, int64, error) {
	if m.FindActiveWithPaginationFunc != nil {
		return m.FindActiveWithPaginationFunc(page, limit, search)
	}
	return nil, 0, errors.New("mock not configured")
}

func (m *MockMemberRepository) FindActiveByDepartment(department string, page, limit int, search string) ([]domain.Member, int64, error) {
	if m.FindActiveByDepartmentFunc != nil {
		return m.FindActiveByDepartmentFunc(department, page, limit, search)
	}
	return nil, 0, errors.New("mock not configured")
}

// ==================== CREATE TESTS ====================

// Test: Upload error harus return error
func TestMemberCreate_ErrorUploadFailed(t *testing.T) {
	mockRepo := &MockMemberRepository{}
	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	service := NewMemberService(mockRepo, mockCloudinary)
	req := requests.CreateMemberRequest{FullName: "Test", Position: "Developer"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Create(context.Background(), req, mockFile)

	if err == nil || err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected 'gagal mengupload foto' error, got: %v", err)
	}
}

// Test: Database error harus rollback (hapus foto yang sudah diupload)
func TestMemberCreate_ErrorDatabaseWithRollback(t *testing.T) {
	deleteCalled := false

	mockRepo := &MockMemberRepository{
		CreateFunc: func(member *domain.Member) error {
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

	service := NewMemberService(mockRepo, mockCloudinary)
	req := requests.CreateMemberRequest{FullName: "Test", Position: "Developer"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Create(context.Background(), req, mockFile)

	if err == nil || err.Error() != "gagal menyimpan member" {
		t.Errorf("Expected 'gagal menyimpan member' error, got: %v", err)
	}
	if !deleteCalled {
		t.Error("Expected rollback: foto harus dihapus ketika database error")
	}
}

// ==================== GET ALL TESTS ====================

// Test: Pagination default values harus diset (page=1, limit=10)
func TestMemberGetAll_ValidationDefaultPagination(t *testing.T) {
	mockRepo := &MockMemberRepository{
		FindAllFunc: func(page, limit int) ([]domain.Member, int64, error) {
			if page != 1 || limit != 10 {
				t.Errorf("Expected default page=1, limit=10, got page=%d, limit=%d", page, limit)
			}
			return []domain.Member{}, 0, nil
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewMemberService(mockRepo, mockCloudinary)

	// Test dengan invalid values
	_, _, _, _, err := service.GetAll(context.Background(), 0, -5)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test: Database error harus return error
func TestMemberGetAll_ErrorDatabaseFailed(t *testing.T) {
	mockRepo := &MockMemberRepository{
		FindAllFunc: func(page, limit int) ([]domain.Member, int64, error) {
			return nil, 0, errors.New("database error")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewMemberService(mockRepo, mockCloudinary)

	_, _, _, _, err := service.GetAll(context.Background(), 1, 10)

	if err == nil || err.Error() != "gagal mengambil data member" {
		t.Errorf("Expected 'gagal mengambil data member' error, got: %v", err)
	}
}

// Test: Kalkulasi lastPage harus benar
func TestMemberGetAll_LogicLastPageCalculation(t *testing.T) {
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
		mockRepo := &MockMemberRepository{
			FindAllFunc: func(page, limit int) ([]domain.Member, int64, error) {
				return []domain.Member{}, tt.total, nil
			},
		}

		mockCloudinary := &MockCloudinaryService{}
		service := NewMemberService(mockRepo, mockCloudinary)

		_, _, lastPage, _, err := service.GetAll(context.Background(), 1, tt.limit)

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
func TestMemberGetByID_ErrorNotFound(t *testing.T) {
	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewMemberService(mockRepo, mockCloudinary)

	_, err := service.GetByID(context.Background(), 999)

	if err == nil || err.Error() != "member tidak ditemukan" {
		t.Errorf("Expected 'member tidak ditemukan' error, got: %v", err)
	}
}

// ==================== UPDATE TESTS ====================

// Test: Update dengan foto baru harus hapus foto lama SETELAH database berhasil
func TestMemberUpdate_LogicDeleteOldPhotoAfterSuccess(t *testing.T) {
	oldPhoto := "old.jpg"
	deleteOldPhotoCalled := false
	updateCalled := false

	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return &domain.Member{
				ID:       1,
				FullName: "Test",
				Position: "Dev",
				PhotoURI: &oldPhoto,
			}, nil
		},
		UpdateFunc: func(member *domain.Member) error {
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

	service := NewMemberService(mockRepo, mockCloudinary)
	req := requests.UpdateMemberRequest{FullName: "Updated"}
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
func TestMemberUpdate_ErrorUploadFailed(t *testing.T) {
	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return &domain.Member{ID: 1, FullName: "Test", Position: "Dev"}, nil
		},
	}

	mockCloudinary := &MockCloudinaryService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	service := NewMemberService(mockRepo, mockCloudinary)
	req := requests.UpdateMemberRequest{FullName: "Updated"}
	mockFile := &multipart.FileHeader{Filename: "photo.jpg"}

	_, err := service.Update(context.Background(), 1, req, mockFile)

	if err == nil || err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected 'gagal mengupload foto' error, got: %v", err)
	}
}

// Test: Database error harus rollback (hapus foto baru, JANGAN hapus foto lama)
func TestMemberUpdate_ErrorDatabaseWithRollback(t *testing.T) {
	oldPhoto := "old.jpg"
	deleteNewPhotoCalled := false
	deleteOldPhotoCalled := false

	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return &domain.Member{ID: 1, FullName: "Test", Position: "Dev", PhotoURI: &oldPhoto}, nil
		},
		UpdateFunc: func(member *domain.Member) error {
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

	service := NewMemberService(mockRepo, mockCloudinary)
	req := requests.UpdateMemberRequest{FullName: "Updated"}
	mockFile := &multipart.FileHeader{Filename: "new.jpg"}

	_, err := service.Update(context.Background(), 1, req, mockFile)

	if err == nil || err.Error() != "gagal mengupdate member" {
		t.Errorf("Expected 'gagal mengupdate member' error, got: %v", err)
	}
	if !deleteNewPhotoCalled {
		t.Error("Expected rollback: foto baru harus dihapus")
	}
	if deleteOldPhotoCalled {
		t.Error("BUG: Foto lama TIDAK BOLEH dihapus ketika database error!")
	}
}

// Test: Not found harus return error
func TestMemberUpdate_ErrorNotFound(t *testing.T) {
	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewMemberService(mockRepo, mockCloudinary)

	_, err := service.Update(context.Background(), 999, requests.UpdateMemberRequest{}, nil)

	if err == nil || err.Error() != "member tidak ditemukan" {
		t.Errorf("Expected 'member tidak ditemukan' error, got: %v", err)
	}
}

// ==================== DELETE TESTS ====================

// Test: Delete harus hapus record database dulu, BARU hapus foto
func TestMemberDelete_LogicDeleteDatabaseFirst(t *testing.T) {
	photo := "photo.jpg"
	databaseDeleteCalled := false

	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return &domain.Member{ID: 1, FullName: "Test", Position: "Dev", PhotoURI: &photo}, nil
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

	service := NewMemberService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test: Database error harus return error (foto TIDAK dihapus)
func TestMemberDelete_ErrorDatabaseFailed(t *testing.T) {
	photo := "photo.jpg"
	photoDeleteCalled := false

	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return &domain.Member{ID: 1, FullName: "Test", Position: "Dev", PhotoURI: &photo}, nil
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

	service := NewMemberService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 1)

	if err == nil || err.Error() != "gagal menghapus member" {
		t.Errorf("Expected 'gagal menghapus member' error, got: %v", err)
	}
	if photoDeleteCalled {
		t.Error("BUG: Foto TIDAK BOLEH dihapus ketika database error!")
	}
}

// Test: Not found harus return error
func TestMemberDelete_ErrorNotFound(t *testing.T) {
	mockRepo := &MockMemberRepository{
		FindByIDFunc: func(id int) (*domain.Member, error) {
			return nil, errors.New("record not found")
		},
	}

	mockCloudinary := &MockCloudinaryService{}
	service := NewMemberService(mockRepo, mockCloudinary)

	err := service.Delete(context.Background(), 999)

	if err == nil || err.Error() != "member tidak ditemukan" {
		t.Errorf("Expected 'member tidak ditemukan' error, got: %v", err)
	}
}
