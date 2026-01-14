package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// MockUserRepository adalah mock untuk UserRepository (untuk testing UserService)
type MockUserRepositoryForUserService struct {
	FindAllFunc     func(page, limit int) ([]domain.User, int64, error)
	FindByIDFunc    func(id int) (*domain.User, error)
	FindByEmailFunc func(email string) (*domain.User, error)
	CreateFunc      func(user *domain.User) error
	UpdateFunc      func(user *domain.User) error
	DeleteFunc      func(id int) error
}

func (m *MockUserRepositoryForUserService) FindAll(page, limit int) ([]domain.User, int64, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(page, limit)
	}
	return nil, 0, errors.New("mock not configured")
}

func (m *MockUserRepositoryForUserService) FindByID(id int) (*domain.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockUserRepositoryForUserService) FindByEmail(email string) (*domain.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockUserRepositoryForUserService) Create(user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) Update(user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockCloudinaryService adalah mock untuk CloudinaryService (untuk testing UserService)
type MockCloudinaryServiceForUserService struct {
	UploadImageFunc    func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	GetImageURLFunc    func(folder string, fileName string) string
	DeleteImageFunc    func(ctx context.Context, folder string, fileName string) error
	UploadFileFunc     func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	DeleteFileFunc     func(ctx context.Context, folder string, filename string) error
	GetFileURLFunc     func(folder string, filename string) string
	GetDownloadURLFunc func(folder string, filename string) string
}

func (m *MockCloudinaryServiceForUserService) UploadImage(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadImageFunc != nil {
		return m.UploadImageFunc(ctx, folder, file)
	}
	return "uploaded-file.jpg", nil
}

func (m *MockCloudinaryServiceForUserService) GetImageURL(folder string, fileName string) string {
	if m.GetImageURLFunc != nil {
		return m.GetImageURLFunc(folder, fileName)
	}
	return "https://res.cloudinary.com/test/" + folder + "/" + fileName
}

func (m *MockCloudinaryServiceForUserService) DeleteImage(ctx context.Context, folder string, fileName string) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, folder, fileName)
	}
	return nil
}

func (m *MockCloudinaryServiceForUserService) UploadFile(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(ctx, folder, file)
	}
	return "", nil
}

func (m *MockCloudinaryServiceForUserService) DeleteFile(ctx context.Context, folder string, filename string) error {
	if m.DeleteFileFunc != nil {
		return m.DeleteFileFunc(ctx, folder, filename)
	}
	return nil
}

func (m *MockCloudinaryServiceForUserService) GetFileURL(folder string, filename string) string {
	if m.GetFileURLFunc != nil {
		return m.GetFileURLFunc(folder, filename)
	}
	return ""
}

func (m *MockCloudinaryServiceForUserService) GetDownloadURL(folder string, filename string) string {
	if m.GetDownloadURLFunc != nil {
		return m.GetDownloadURLFunc(folder, filename)
	}
	return ""
}

// MockActivityLogRepoForUser adalah mock untuk ActivityLogRepository
type MockActivityLogRepoForUser struct {
	CreateFunc func(log *domain.ActivityLog) error
}

func (m *MockActivityLogRepoForUser) Create(log *domain.ActivityLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(log)
	}
	return nil
}

func (m *MockActivityLogRepoForUser) GetActivityLogs(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
	return nil, 0, nil
}

// Helper functions untuk membuat pointer dari value
func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func boolPtr(b bool) *bool    { return &b }

// ============================================================
// Test Cases untuk GetAllUsers
// ============================================================

func TestGetAllUsers(t *testing.T) {
	mockUsers := []domain.User{
		{ID: 1, FullName: "Admin User", Email: "admin@example.com", Role: 1, IsActive: true},
		{ID: 2, FullName: "Regular User", Email: "user@example.com", Role: 2, IsActive: true},
	}

	tests := []struct {
		name          string
		page, limit   int
		mockUsers     []domain.User
		mockTotal     int64
		mockErr       error
		expectedLen   int
		expectedPage  int
		expectedLast  int
		expectedTotal int64
		expectedErr   error
	}{
		{
			name:          "success dengan pagination",
			page:          1,
			limit:         20,
			mockUsers:     mockUsers,
			mockTotal:     2,
			expectedLen:   2,
			expectedPage:  1,
			expectedLast:  1,
			expectedTotal: 2,
		},
		{
			name:          "pagination calculation",
			page:          2,
			limit:         20,
			mockUsers:     mockUsers[:1],
			mockTotal:     45,
			expectedLen:   1,
			expectedPage:  2,
			expectedLast:  3, // ceil(45/20)
			expectedTotal: 45,
		},
		{
			name:          "empty list",
			page:          1,
			limit:         20,
			mockUsers:     []domain.User{},
			mockTotal:     0,
			expectedLen:   0,
			expectedPage:  1,
			expectedLast:  0,
			expectedTotal: 0,
		},
		{
			name:        "repository error",
			page:        1,
			limit:       20,
			mockErr:     errors.New("database connection failed"),
			expectedErr: ErrUserFetchFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepositoryForUserService{
				FindAllFunc: func(page, limit int) ([]domain.User, int64, error) {
					return tt.mockUsers, tt.mockTotal, tt.mockErr
				},
			}
			service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
			users, currentPage, lastPage, total, err := service.GetAllUsers(context.Background(), tt.page, tt.limit)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(users) != tt.expectedLen {
				t.Errorf("expected %d users, got %d", tt.expectedLen, len(users))
			}
			if currentPage != tt.expectedPage {
				t.Errorf("expected currentPage %d, got %d", tt.expectedPage, currentPage)
			}
			if lastPage != tt.expectedLast {
				t.Errorf("expected lastPage %d, got %d", tt.expectedLast, lastPage)
			}
			if total != tt.expectedTotal {
				t.Errorf("expected total %d, got %d", tt.expectedTotal, total)
			}
		})
	}
}

// TestGetAllUsers_DefaultValues menguji default values untuk pagination
func TestGetAllUsers_DefaultValues(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindAllFunc: func(page, limit int) ([]domain.User, int64, error) {
			if page != 1 {
				t.Errorf("expected page to be defaulted to 1, got %d", page)
			}
			if limit != 20 {
				t.Errorf("expected limit to be defaulted to 20, got %d", limit)
			}
			return []domain.User{}, 0, nil
		},
	}
	service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
	service.GetAllUsers(context.Background(), 0, 0) // Pass invalid values to test defaults
}

// ============================================================
// Test Cases untuk GetUserByID
// ============================================================

func TestGetUserByID(t *testing.T) {
	mockUser := &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", Role: 2, IsActive: true}

	tests := []struct {
		name        string
		userID      int
		mockUser    *domain.User
		mockErr     error
		expectedErr error
	}{
		{
			name:     "success",
			userID:   1,
			mockUser: mockUser,
		},
		{
			name:        "not found",
			userID:      999,
			mockErr:     errors.New("record not found"),
			expectedErr: ErrUserNotFound,
		},
		{
			name:        "repository error",
			userID:      1,
			mockErr:     errors.New("database connection failed"),
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepositoryForUserService{
				FindByIDFunc: func(id int) (*domain.User, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					if id == tt.userID {
						return tt.mockUser, nil
					}
					return nil, errors.New("not found")
				},
			}
			service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
			user, err := service.GetUserByID(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				if user != nil {
					t.Errorf("expected nil user, got %v", user)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if user.ID != tt.mockUser.ID || user.Email != tt.mockUser.Email {
				t.Errorf("user mismatch: expected %v, got %v", tt.mockUser, user)
			}
		})
	}
}

// ============================================================
// Test Cases untuk CreateUser
// ============================================================

// TestCreateUser_PasswordValidation menguji validasi password
func TestCreateUser_PasswordValidation(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{"only letters", "passwordonly", ErrInvalidPassword},
		{"only numbers", "12345678", ErrInvalidPassword},
		{"valid combo", "password123", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepositoryForUserService{
				FindByEmailFunc: func(email string) (*domain.User, error) { return nil, errors.New("not found") },
				CreateFunc:      func(user *domain.User) error { user.ID = 1; return nil },
			}
			service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
			req := &requests.CreateUserRequest{
				FullName: "Test User",
				Email:    "test@example.com",
				Password: tt.password,
			}
			user, err := service.CreateUser(context.Background(), req, nil)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				if user != nil {
					t.Errorf("expected nil user, got %v", user)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Error("expected user, got nil")
				}
			}
		})
	}
}

// TestCreateUser_ErrorCases menguji berbagai error cases
func TestCreateUser_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService)
		photoFile   *multipart.FileHeader
		expectedErr error
	}{
		{
			name: "email already exists",
			setupMock: func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService) {
				return &MockUserRepositoryForUserService{
					FindByEmailFunc: func(email string) (*domain.User, error) {
						return &domain.User{ID: 1, Email: email}, nil
					},
				}, &MockCloudinaryServiceForUserService{}
			},
			expectedErr: ErrEmailAlreadyExists,
		},
		{
			name: "repository error",
			setupMock: func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService) {
				return &MockUserRepositoryForUserService{
					FindByEmailFunc: func(email string) (*domain.User, error) { return nil, errors.New("not found") },
					CreateFunc:      func(user *domain.User) error { return errors.New("database connection failed") },
				}, &MockCloudinaryServiceForUserService{}
			},
			expectedErr: ErrUserCreateFailed,
		},
		{
			name: "upload photo error",
			setupMock: func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService) {
				return &MockUserRepositoryForUserService{
						FindByEmailFunc: func(email string) (*domain.User, error) { return nil, errors.New("not found") },
					}, &MockCloudinaryServiceForUserService{
						UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
							return "", errors.New("cloudinary error")
						},
					}
			},
			photoFile:   &multipart.FileHeader{Filename: "test.jpg"},
			expectedErr: ErrPhotoUploadFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, mockCloudinary := tt.setupMock()
			service := NewUserService(mockRepo, mockCloudinary, &MockActivityLogRepoForUser{})
			req := &requests.CreateUserRequest{
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "password123",
			}
			user, err := service.CreateUser(context.Background(), req, tt.photoFile)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if user != nil {
				t.Errorf("expected nil user, got %v", user)
			}
		})
	}
}

// TestCreateUser_WithPhotoFile menguji pembuatan user dengan photo file
func TestCreateUser_WithPhotoFile(t *testing.T) {
	var uploadCalled bool
	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) { return nil, errors.New("not found") },
		CreateFunc:      func(user *domain.User) error { user.ID = 1; return nil },
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			uploadCalled = true
			return "uploaded-photo.jpg", nil
		},
		GetImageURLFunc: func(folder string, fileName string) string {
			return "https://res.cloudinary.com/test/" + folder + "/" + fileName
		},
	}

	service := NewUserService(mockRepo, mockCloudinary, &MockActivityLogRepoForUser{})
	req := &requests.CreateUserRequest{FullName: "Test User", Email: "test@example.com", Password: "password123"}
	user, err := service.CreateUser(context.Background(), req, &multipart.FileHeader{Filename: "test.jpg"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !uploadCalled {
		t.Error("expected cloudinary upload to be called")
	}
	expectedURL := "https://res.cloudinary.com/test/users/avatars/uploaded-photo.jpg"
	if user.PhotoURI == nil || *user.PhotoURI != expectedURL {
		t.Errorf("expected photo URI '%s', got '%v'", expectedURL, user.PhotoURI)
	}
}

// TestCreateUser_RollbackOnDBError menguji bahwa foto dihapus jika database create gagal
func TestCreateUser_RollbackOnDBError(t *testing.T) {
	var uploadCalled, deletePhotoCalled bool
	uploadedFileName := "uploaded-photo.jpg"

	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) { return nil, errors.New("not found") },
		CreateFunc:      func(user *domain.User) error { return errors.New("database connection failed") },
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			uploadCalled = true
			return uploadedFileName, nil
		},
		DeleteImageFunc: func(ctx context.Context, folder string, fileName string) error {
			if fileName == uploadedFileName {
				deletePhotoCalled = true
			}
			return nil
		},
	}

	service := NewUserService(mockRepo, mockCloudinary, &MockActivityLogRepoForUser{})
	req := &requests.CreateUserRequest{FullName: "Test User", Email: "test@example.com", Password: "password123"}
	user, err := service.CreateUser(context.Background(), req, &multipart.FileHeader{Filename: "test.jpg"})

	if !errors.Is(err, ErrUserCreateFailed) {
		t.Errorf("expected ErrUserCreateFailed, got %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %v", user)
	}
	if !uploadCalled {
		t.Error("expected cloudinary upload to be called")
	}
	if !deletePhotoCalled {
		t.Error("expected uploaded photo to be deleted after DB error (rollback)")
	}
}

// ============================================================
// Test Cases untuk UpdateUser
// ============================================================

// TestUpdateUser_Success menguji update user berhasil
func TestUpdateUser_Success(t *testing.T) {
	existingUser := &domain.User{ID: 1, FullName: "Old Name", Email: "old@example.com", PasswordHash: "oldhash", Role: 2, IsActive: true}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc:    func(id int) (*domain.User, error) { return existingUser, nil },
		FindByEmailFunc: func(email string) (*domain.User, error) { return nil, errors.New("not found") },
		UpdateFunc:      func(user *domain.User) error { return nil },
	}
	service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
	req := &requests.UpdateUserRequest{
		FullName: strPtr("New Name"),
		Email:    strPtr("new@example.com"),
		Role:     intPtr(1),
		IsActive: boolPtr(false),
	}

	user, err := service.UpdateUser(context.Background(), 1, req, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.FullName != "New Name" || user.Email != "new@example.com" || user.Role != 1 || user.IsActive != false {
		t.Errorf("user not updated correctly: %+v", user)
	}
}

// TestUpdateUser_PasswordValidation menguji validasi password saat update
func TestUpdateUser_PasswordValidation(t *testing.T) {
	existingUser := &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", PasswordHash: "oldhash", Role: 2, IsActive: true}

	tests := []struct {
		name        string
		password    *string
		expectedErr error
		hashChanged bool
	}{
		{"only letters", strPtr("onlyletters"), ErrInvalidPassword, false},
		{"only numbers", strPtr("12345678"), ErrInvalidPassword, false},
		{"valid password", strPtr("newpassword123"), nil, true},
		{"no password change", nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepositoryForUserService{
				FindByIDFunc: func(id int) (*domain.User, error) {
					// Return a copy to avoid mutation across tests
					return &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", PasswordHash: "oldhash", Role: 2, IsActive: true}, nil
				},
				UpdateFunc: func(user *domain.User) error { return nil },
			}
			service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
			req := &requests.UpdateUserRequest{
				FullName: strPtr(existingUser.FullName),
				Email:    strPtr(existingUser.Email),
				Role:     intPtr(existingUser.Role),
				IsActive: boolPtr(existingUser.IsActive),
				Password: tt.password,
			}

			user, err := service.UpdateUser(context.Background(), 1, req, nil)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.hashChanged && user.PasswordHash == "oldhash" {
				t.Error("expected password hash to be updated")
			}
			if !tt.hashChanged && user.PasswordHash != "oldhash" {
				t.Errorf("expected password hash to remain 'oldhash', got '%s'", user.PasswordHash)
			}
		})
	}
}

// TestUpdateUser_ErrorCases menguji berbagai error cases untuk update
func TestUpdateUser_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		setupMock   func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService)
		photoFile   *multipart.FileHeader
		expectedErr error
	}{
		{
			name:   "not found",
			userID: 999,
			setupMock: func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService) {
				return &MockUserRepositoryForUserService{
					FindByIDFunc: func(id int) (*domain.User, error) { return nil, errors.New("record not found") },
				}, &MockCloudinaryServiceForUserService{}
			},
			expectedErr: ErrUserNotFound,
		},
		{
			name:   "email already used by other",
			userID: 1,
			setupMock: func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService) {
				return &MockUserRepositoryForUserService{
					FindByIDFunc: func(id int) (*domain.User, error) {
						return &domain.User{ID: 1, FullName: "User One", Email: "user1@example.com", Role: 2, IsActive: true}, nil
					},
					FindByEmailFunc: func(email string) (*domain.User, error) {
						return &domain.User{ID: 2, Email: "user2@example.com"}, nil // Different user owns this email
					},
				}, &MockCloudinaryServiceForUserService{}
			},
			expectedErr: ErrEmailAlreadyUsed,
		},
		{
			name:   "repository error",
			userID: 1,
			setupMock: func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService) {
				return &MockUserRepositoryForUserService{
					FindByIDFunc: func(id int) (*domain.User, error) {
						return &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", Role: 2, IsActive: true}, nil
					},
					UpdateFunc: func(user *domain.User) error { return errors.New("database connection failed") },
				}, &MockCloudinaryServiceForUserService{}
			},
			expectedErr: ErrUserUpdateFailed,
		},
		{
			name:   "upload photo error",
			userID: 1,
			setupMock: func() (*MockUserRepositoryForUserService, *MockCloudinaryServiceForUserService) {
				return &MockUserRepositoryForUserService{
						FindByIDFunc: func(id int) (*domain.User, error) {
							return &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", Role: 2, IsActive: true}, nil
						},
					}, &MockCloudinaryServiceForUserService{
						UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
							return "", errors.New("cloudinary error")
						},
					}
			},
			photoFile:   &multipart.FileHeader{Filename: "new-photo.jpg"},
			expectedErr: ErrPhotoUploadFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, mockCloudinary := tt.setupMock()
			service := NewUserService(mockRepo, mockCloudinary, &MockActivityLogRepoForUser{})
			req := &requests.UpdateUserRequest{
				FullName: strPtr("Updated Name"),
				Email:    strPtr("user2@example.com"),
				Role:     intPtr(2),
				IsActive: boolPtr(true),
			}

			user, err := service.UpdateUser(context.Background(), tt.userID, req, tt.photoFile)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if user != nil {
				t.Errorf("expected nil user, got %v", user)
			}
		})
	}
}

// TestUpdateUser_SameEmailAllowed menguji email tidak berubah (diperbolehkan)
func TestUpdateUser_SameEmailAllowed(t *testing.T) {
	existingUser := &domain.User{ID: 1, FullName: "User One", Email: "user1@example.com", Role: 2, IsActive: true}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) { return existingUser, nil },
		UpdateFunc:   func(user *domain.User) error { return nil },
	}
	service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
	req := &requests.UpdateUserRequest{
		FullName: strPtr("User One Updated"),
		Email:    strPtr("user1@example.com"), // Email sama
		Role:     intPtr(1),
		IsActive: boolPtr(true),
	}

	user, err := service.UpdateUser(context.Background(), 1, req, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.FullName != "User One Updated" {
		t.Errorf("expected full name 'User One Updated', got %s", user.FullName)
	}
}

// TestUpdateUser_PartialUpdate menguji partial update (hanya beberapa field)
func TestUpdateUser_PartialUpdate(t *testing.T) {
	tests := []struct {
		name     string
		request  *requests.UpdateUserRequest
		validate func(t *testing.T, user *domain.User)
	}{
		{
			name:    "only full name",
			request: &requests.UpdateUserRequest{FullName: strPtr("New Name")},
			validate: func(t *testing.T, user *domain.User) {
				if user.FullName != "New Name" {
					t.Errorf("expected full name 'New Name', got %s", user.FullName)
				}
				if user.Email != "test@example.com" {
					t.Errorf("expected email to remain 'test@example.com', got %s", user.Email)
				}
			},
		},
		{
			name:    "only role",
			request: &requests.UpdateUserRequest{Role: intPtr(1)},
			validate: func(t *testing.T, user *domain.User) {
				if user.Role != 1 {
					t.Errorf("expected role 1, got %d", user.Role)
				}
				if user.FullName != "Test User" {
					t.Errorf("expected full name to remain 'Test User', got %s", user.FullName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepositoryForUserService{
				FindByIDFunc: func(id int) (*domain.User, error) {
					return &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", Role: 2, IsActive: true}, nil
				},
				UpdateFunc: func(user *domain.User) error { return nil },
			}
			service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
			user, err := service.UpdateUser(context.Background(), 1, tt.request, nil)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			tt.validate(t, user)
		})
	}
}

// TestUpdateUser_WithNewPhotoFile menguji update user dengan foto baru
func TestUpdateUser_WithNewPhotoFile(t *testing.T) {
	oldPhoto := "old-photo.jpg"
	existingUser := &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", PhotoURI: &oldPhoto, Role: 2, IsActive: true}
	var uploadCalled, deleteOldCalled bool

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) { return existingUser, nil },
		UpdateFunc:   func(user *domain.User) error { return nil },
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			uploadCalled = true
			return "new-uploaded-photo.jpg", nil
		},
		GetImageURLFunc: func(folder string, fileName string) string {
			return "https://res.cloudinary.com/test/" + folder + "/" + fileName
		},
		DeleteImageFunc: func(ctx context.Context, folder string, fileName string) error {
			if fileName == oldPhoto {
				deleteOldCalled = true
			}
			return nil
		},
	}

	service := NewUserService(mockRepo, mockCloudinary, &MockActivityLogRepoForUser{})
	req := &requests.UpdateUserRequest{FullName: strPtr("Test User"), Email: strPtr("test@example.com"), Role: intPtr(2), IsActive: boolPtr(true)}
	user, err := service.UpdateUser(context.Background(), 1, req, &multipart.FileHeader{Filename: "new-photo.jpg"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !uploadCalled {
		t.Error("expected cloudinary upload to be called")
	}
	if !deleteOldCalled {
		t.Error("expected old photo to be deleted from cloudinary")
	}
	expectedURL := "https://res.cloudinary.com/test/users/avatars/new-uploaded-photo.jpg"
	if user.PhotoURI == nil || *user.PhotoURI != expectedURL {
		t.Errorf("expected photo URI '%s', got '%v'", expectedURL, user.PhotoURI)
	}
}

// ============================================================
// Test Cases untuk DeleteUser
// ============================================================

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		setupMock   func() *MockUserRepositoryForUserService
		expectedErr error
	}{
		{
			name:   "success",
			userID: 1,
			setupMock: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					FindByIDFunc: func(id int) (*domain.User, error) {
						return &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", Role: 2, IsActive: true}, nil
					},
					DeleteFunc: func(id int) error { return nil },
				}
			},
		},
		{
			name:   "not found",
			userID: 999,
			setupMock: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					FindByIDFunc: func(id int) (*domain.User, error) { return nil, errors.New("record not found") },
				}
			},
			expectedErr: ErrUserNotFound,
		},
		{
			name:   "repository error",
			userID: 1,
			setupMock: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					FindByIDFunc: func(id int) (*domain.User, error) {
						return &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", Role: 2, IsActive: true}, nil
					},
					DeleteFunc: func(id int) error { return errors.New("database connection failed") },
				}
			},
			expectedErr: ErrUserDeleteFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			service := NewUserService(mockRepo, &MockCloudinaryServiceForUserService{}, &MockActivityLogRepoForUser{})
			err := service.DeleteUser(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestDeleteUser_NoPhotoCleanupOnSoftDelete menguji bahwa foto TIDAK dihapus saat soft delete
func TestDeleteUser_NoPhotoCleanupOnSoftDelete(t *testing.T) {
	photoURI := "user-photo.jpg"
	existingUser := &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com", PhotoURI: &photoURI, Role: 2, IsActive: true}
	var deletePhotoCalled bool

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) { return existingUser, nil },
		DeleteFunc:   func(id int) error { return nil },
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{
		DeleteImageFunc: func(ctx context.Context, folder string, fileName string) error {
			if fileName == photoURI {
				deletePhotoCalled = true
			}
			return nil
		},
	}

	service := NewUserService(mockRepo, mockCloudinary, &MockActivityLogRepoForUser{})
	err := service.DeleteUser(context.Background(), 1)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if deletePhotoCalled {
		t.Error("photo should NOT be deleted during soft delete")
	}
}
