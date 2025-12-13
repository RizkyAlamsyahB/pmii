package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
)

// MockUserRepository adalah mock untuk UserRepository (untuk testing UserService)
type MockUserRepositoryForUserService struct {
	FindAllFunc     func() ([]domain.User, error)
	FindByIDFunc    func(id int) (*domain.User, error)
	FindByEmailFunc func(email string) (*domain.User, error)
	CreateFunc      func(user *domain.User) error
	UpdateFunc      func(user *domain.User) error
	DeleteFunc      func(id int) error
}

func (m *MockUserRepositoryForUserService) FindAll() ([]domain.User, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, errors.New("mock not configured")
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

// Update method with mock function support
func (m *MockUserRepositoryForUserService) Update(user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

// Delete method with mock function support
func (m *MockUserRepositoryForUserService) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockCloudinaryService adalah mock untuk CloudinaryService (untuk testing UserService)
type MockCloudinaryServiceForUserService struct {
	UploadImageFunc  func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error)
	GetImageURLFunc  func(folder string, fileName string) string
	DeleteImageFunc  func(ctx context.Context, folder string, fileName string) error
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

// ============================================================
// Test Cases untuk GetAllUsers
// ============================================================

// TestGetAllUsers_Success menguji pengambilan semua users berhasil
func TestGetAllUsers_Success(t *testing.T) {
	// Mock data
	mockUsers := []domain.User{
		{
			ID:       1,
			FullName: "Admin User",
			Email:    "admin@example.com",
			Role:     1,
			IsActive: true,
		},
		{
			ID:       2,
			FullName: "Regular User",
			Email:    "user@example.com",
			Role:     2,
			IsActive: true,
		},
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindAllFunc: func() ([]domain.User, error) {
			return mockUsers, nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	users, err := userService.GetAllUsers()

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if users == nil {
		t.Error("Expected users list, got nil")
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	if users[0].ID != 1 {
		t.Errorf("Expected first user ID to be 1, got %d", users[0].ID)
	}
	if users[1].Email != "user@example.com" {
		t.Errorf("Expected second user email to be user@example.com, got %s", users[1].Email)
	}
}

// TestGetAllUsers_EmptyList menguji ketika tidak ada users di database
func TestGetAllUsers_EmptyList(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindAllFunc: func() ([]domain.User, error) {
			return []domain.User{}, nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	users, err := userService.GetAllUsers()

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if users == nil {
		t.Error("Expected empty slice, got nil")
	}
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}
}

// TestGetAllUsers_RepositoryError menguji ketika repository gagal
func TestGetAllUsers_RepositoryError(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindAllFunc: func() ([]domain.User, error) {
			return nil, errors.New("database connection failed")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	users, err := userService.GetAllUsers()

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "gagal mengambil data user" {
		t.Errorf("Expected error message 'gagal mengambil data user', got '%s'", err.Error())
	}
	if users != nil {
		t.Errorf("Expected nil users, got %v", users)
	}
}

// ============================================================
// Test Cases untuk GetUserByID
// ============================================================

// TestGetUserByID_Success menguji pengambilan user by ID berhasil
func TestGetUserByID_Success(t *testing.T) {
	// Mock data
	mockUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			if id == 1 {
				return mockUser, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	user, err := userService.GetUserByID(1)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.ID != 1 {
		t.Errorf("Expected user ID to be 1, got %d", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}
	if user.FullName != "Test User" {
		t.Errorf("Expected full name 'Test User', got %s", user.FullName)
	}
}

// TestGetUserByID_NotFound menguji ketika user tidak ditemukan
func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return nil, errors.New("record not found")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	user, err := userService.GetUserByID(999)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "user tidak ditemukan" {
		t.Errorf("Expected error message 'user tidak ditemukan', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestGetUserByID_RepositoryError menguji ketika repository gagal
func TestGetUserByID_RepositoryError(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return nil, errors.New("database connection failed")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	user, err := userService.GetUserByID(1)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "user tidak ditemukan" {
		t.Errorf("Expected error message 'user tidak ditemukan', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// ============================================================
// Test Cases untuk CreateUser
// ============================================================

// TestCreateUser_Success menguji pembuatan user berhasil tanpa foto
func TestCreateUser_Success(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("not found") // Email belum terdaftar
		},
		CreateFunc: func(user *domain.User) error {
			user.ID = 1 // Simulate auto increment
			return nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123", // Kombinasi huruf dan angka
	}

	ctx := context.Background()
	user, err := userService.CreateUser(ctx, req, nil) // tanpa foto

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.FullName != "Test User" {
		t.Errorf("Expected full name 'Test User', got %s", user.FullName)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", user.Email)
	}
	if user.Role != 2 {
		t.Errorf("Expected role 2 (Author), got %d", user.Role)
	}
	if !user.IsActive {
		t.Error("Expected user to be active")
	}
}

// TestCreateUser_PasswordOnlyLetters menguji password hanya huruf (tanpa angka)
func TestCreateUser_PasswordOnlyLetters(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "passwordonly", // Hanya huruf, tanpa angka
	}

	ctx := context.Background()
	user, err := userService.CreateUser(ctx, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "password harus kombinasi huruf dan angka" {
		t.Errorf("Expected error message 'password harus kombinasi huruf dan angka', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestCreateUser_PasswordOnlyNumbers menguji password hanya angka (tanpa huruf)
func TestCreateUser_PasswordOnlyNumbers(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "12345678", // Hanya angka, tanpa huruf
	}

	ctx := context.Background()
	user, err := userService.CreateUser(ctx, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "password harus kombinasi huruf dan angka" {
		t.Errorf("Expected error message 'password harus kombinasi huruf dan angka', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestCreateUser_EmailAlreadyExists menguji email sudah terdaftar
func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Existing User",
		Email:    "existing@example.com",
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return existingUser, nil // Email sudah ada
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	ctx := context.Background()
	user, err := userService.CreateUser(ctx, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "email sudah terdaftar" {
		t.Errorf("Expected error message 'email sudah terdaftar', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestCreateUser_RepositoryError menguji ketika repository gagal create
func TestCreateUser_RepositoryError(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		CreateFunc: func(user *domain.User) error {
			return errors.New("database connection failed")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	ctx := context.Background()
	user, err := userService.CreateUser(ctx, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "gagal membuat user" {
		t.Errorf("Expected error message 'gagal membuat user', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestCreateUser_WithPhotoFile menguji pembuatan user dengan photo file (cloudinary upload)
func TestCreateUser_WithPhotoFile(t *testing.T) {
	var uploadCalled bool

	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		CreateFunc: func(user *domain.User) error {
			user.ID = 1
			return nil
		},
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

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	ctx := context.Background()
	// Simulasikan file header (dalam real test perlu mock file header)
	photoFile := &multipart.FileHeader{Filename: "test.jpg"}
	user, err := userService.CreateUser(ctx, req, photoFile)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if !uploadCalled {
		t.Error("Expected cloudinary upload to be called")
	}
	if user.PhotoURI == nil {
		t.Error("Expected photo URI, got nil")
	}
	expectedURL := "https://res.cloudinary.com/test/users/avatars/uploaded-photo.jpg"
	if *user.PhotoURI != expectedURL {
		t.Errorf("Expected photo URI '%s', got '%s'", expectedURL, *user.PhotoURI)
	}
}

// TestCreateUser_UploadPhotoError menguji ketika upload foto gagal
func TestCreateUser_UploadPhotoError(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	ctx := context.Background()
	photoFile := &multipart.FileHeader{Filename: "test.jpg"}
	user, err := userService.CreateUser(ctx, req, photoFile)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected error message 'gagal mengupload foto', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// ============================================================
// Test Cases untuk UpdateUser
// ============================================================

// Helper functions untuk membuat pointer dari value
func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func boolPtr(b bool) *bool    { return &b }

// TestUpdateUser_Success menguji update user berhasil tanpa foto baru
func TestUpdateUser_Success(t *testing.T) {
	existingUser := &domain.User{
		ID:           1,
		FullName:     "Old Name",
		Email:        "old@example.com",
		PasswordHash: "oldhash",
		Role:         2,
		IsActive:     true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			if id == 1 {
				return existingUser, nil
			}
			return nil, errors.New("not found")
		},
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("not found") // Email baru belum dipakai
		},
		UpdateFunc: func(user *domain.User) error {
			return nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("New Name"),
		Email:    strPtr("new@example.com"),
		Role:     intPtr(1),
		IsActive: boolPtr(false),
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil) // tanpa foto baru

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.FullName != "New Name" {
		t.Errorf("Expected full name 'New Name', got %s", user.FullName)
	}
	if user.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got %s", user.Email)
	}
	if user.Role != 1 {
		t.Errorf("Expected role 1, got %d", user.Role)
	}
	if user.IsActive != false {
		t.Error("Expected user to be inactive")
	}
}

// TestUpdateUser_NotFound menguji update user yang tidak ditemukan
func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return nil, errors.New("record not found")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("New Name"),
		Email:    strPtr("new@example.com"),
		Role:     intPtr(1),
		IsActive: boolPtr(true),
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 999, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "user tidak ditemukan" {
		t.Errorf("Expected error message 'user tidak ditemukan', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestUpdateUser_EmailAlreadyUsedByOther menguji email sudah dipakai user lain
func TestUpdateUser_EmailAlreadyUsedByOther(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "User One",
		Email:    "user1@example.com",
		Role:     2,
		IsActive: true,
	}

	otherUser := &domain.User{
		ID:       2,
		FullName: "User Two",
		Email:    "user2@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			if id == 1 {
				return existingUser, nil
			}
			return nil, errors.New("not found")
		},
		FindByEmailFunc: func(email string) (*domain.User, error) {
			if email == "user2@example.com" {
				return otherUser, nil // Email sudah dipakai user lain
			}
			return nil, errors.New("not found")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("User One Updated"),
		Email:    strPtr("user2@example.com"), // Mencoba pakai email user lain
		Role:     intPtr(2),
		IsActive: boolPtr(true),
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "email sudah digunakan user lain" {
		t.Errorf("Expected error message 'email sudah digunakan user lain', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestUpdateUser_SameEmailAllowed menguji email tidak berubah (diperbolehkan)
func TestUpdateUser_SameEmailAllowed(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "User One",
		Email:    "user1@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			if id == 1 {
				return existingUser, nil
			}
			return nil, errors.New("not found")
		},
		UpdateFunc: func(user *domain.User) error {
			return nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("User One Updated"),
		Email:    strPtr("user1@example.com"), // Email sama, tidak berubah
		Role:     intPtr(1),
		IsActive: boolPtr(true),
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.FullName != "User One Updated" {
		t.Errorf("Expected full name 'User One Updated', got %s", user.FullName)
	}
}

// TestUpdateUser_PasswordOnlyLetters menguji password hanya huruf
func TestUpdateUser_PasswordOnlyLetters(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("Test User"),
		Email:    strPtr("test@example.com"),
		Role:     intPtr(2),
		IsActive: boolPtr(true),
		Password: strPtr("onlyletters"), // Hanya huruf
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "password harus kombinasi huruf dan angka" {
		t.Errorf("Expected error message 'password harus kombinasi huruf dan angka', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestUpdateUser_PasswordOnlyNumbers menguji password hanya angka
func TestUpdateUser_PasswordOnlyNumbers(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("Test User"),
		Email:    strPtr("test@example.com"),
		Role:     intPtr(2),
		IsActive: boolPtr(true),
		Password: strPtr("12345678"), // Hanya angka
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "password harus kombinasi huruf dan angka" {
		t.Errorf("Expected error message 'password harus kombinasi huruf dan angka', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestUpdateUser_WithValidPassword menguji update dengan password valid
func TestUpdateUser_WithValidPassword(t *testing.T) {
	existingUser := &domain.User{
		ID:           1,
		FullName:     "Test User",
		Email:        "test@example.com",
		PasswordHash: "oldhash",
		Role:         2,
		IsActive:     true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
		UpdateFunc: func(user *domain.User) error {
			return nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("Test User"),
		Email:    strPtr("test@example.com"),
		Role:     intPtr(2),
		IsActive: boolPtr(true),
		Password: strPtr("newpassword123"), // Password valid
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	// Password harus berubah (di-hash)
	if user.PasswordHash == "oldhash" {
		t.Error("Expected password hash to be updated")
	}
}

// TestUpdateUser_WithoutPassword menguji update tanpa mengubah password
func TestUpdateUser_WithoutPassword(t *testing.T) {
	existingUser := &domain.User{
		ID:           1,
		FullName:     "Test User",
		Email:        "test@example.com",
		PasswordHash: "existinghash",
		Role:         2,
		IsActive:     true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
		UpdateFunc: func(user *domain.User) error {
			return nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("Updated Name"),
		Email:    strPtr("test@example.com"),
		Role:     intPtr(2),
		IsActive: boolPtr(true),
		// Password nil, tidak diubah
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	// Password tidak boleh berubah
	if user.PasswordHash != "existinghash" {
		t.Errorf("Expected password hash to remain 'existinghash', got '%s'", user.PasswordHash)
	}
	if user.FullName != "Updated Name" {
		t.Errorf("Expected full name 'Updated Name', got %s", user.FullName)
	}
}

// TestUpdateUser_RepositoryError menguji ketika repository gagal update
func TestUpdateUser_RepositoryError(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
		UpdateFunc: func(user *domain.User) error {
			return errors.New("database connection failed")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("Updated Name"),
		Email:    strPtr("test@example.com"),
		Role:     intPtr(2),
		IsActive: boolPtr(true),
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "gagal mengupdate user" {
		t.Errorf("Expected error message 'gagal mengupdate user', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestUpdateUser_WithNewPhotoFile menguji update user dengan foto baru (cloudinary upload)
func TestUpdateUser_WithNewPhotoFile(t *testing.T) {
	oldPhoto := "old-photo.jpg"
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		PhotoURI: &oldPhoto,
		Role:     2,
		IsActive: true,
	}

	var uploadCalled, deleteOldCalled bool

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
		UpdateFunc: func(user *domain.User) error {
			return nil
		},
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

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("Test User"),
		Email:    strPtr("test@example.com"),
		Role:     intPtr(2),
		IsActive: boolPtr(true),
	}

	ctx := context.Background()
	photoFile := &multipart.FileHeader{Filename: "new-photo.jpg"}
	user, err := userService.UpdateUser(ctx, 1, req, photoFile)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if !uploadCalled {
		t.Error("Expected cloudinary upload to be called")
	}
	if !deleteOldCalled {
		t.Error("Expected old photo to be deleted from cloudinary")
	}
	expectedURL := "https://res.cloudinary.com/test/users/avatars/new-uploaded-photo.jpg"
	if user.PhotoURI == nil || *user.PhotoURI != expectedURL {
		t.Errorf("Expected photo URI '%s', got '%v'", expectedURL, user.PhotoURI)
	}
}

// TestUpdateUser_UploadPhotoError menguji ketika upload foto gagal saat update
func TestUpdateUser_UploadPhotoError(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{
		UploadImageFunc: func(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
			return "", errors.New("cloudinary error")
		},
	}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("Test User"),
		Email:    strPtr("test@example.com"),
		Role:     intPtr(2),
		IsActive: boolPtr(true),
	}

	ctx := context.Background()
	photoFile := &multipart.FileHeader{Filename: "new-photo.jpg"}
	user, err := userService.UpdateUser(ctx, 1, req, photoFile)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "gagal mengupload foto" {
		t.Errorf("Expected error message 'gagal mengupload foto', got '%s'", err.Error())
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

// TestUpdateUser_PartialUpdate_OnlyFullName menguji partial update hanya FullName
func TestUpdateUser_PartialUpdate_OnlyFullName(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Old Name",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
		UpdateFunc: func(user *domain.User) error {
			return nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		FullName: strPtr("New Name"), // Hanya update FullName
		// Field lain nil
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	// FullName harus berubah
	if user.FullName != "New Name" {
		t.Errorf("Expected full name 'New Name', got %s", user.FullName)
	}
	// Field lain harus tetap sama
	if user.Email != "test@example.com" {
		t.Errorf("Expected email to remain 'test@example.com', got %s", user.Email)
	}
	if user.Role != 2 {
		t.Errorf("Expected role to remain 2, got %d", user.Role)
	}
	if user.IsActive != true {
		t.Error("Expected user to remain active")
	}
}

// TestUpdateUser_PartialUpdate_OnlyRole menguji partial update hanya Role
func TestUpdateUser_PartialUpdate_OnlyRole(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
		UpdateFunc: func(user *domain.User) error {
			return nil
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	req := &requests.UpdateUserRequest{
		Role: intPtr(1), // Hanya update Role
		// Field lain nil
	}

	ctx := context.Background()
	user, err := userService.UpdateUser(ctx, 1, req, nil)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	// Role harus berubah
	if user.Role != 1 {
		t.Errorf("Expected role 1, got %d", user.Role)
	}
	// Field lain harus tetap sama
	if user.FullName != "Test User" {
		t.Errorf("Expected full name to remain 'Test User', got %s", user.FullName)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email to remain 'test@example.com', got %s", user.Email)
	}
	if user.IsActive != true {
		t.Error("Expected user to remain active")
	}
}

// ============================================================
// Test Cases untuk DeleteUser
// ============================================================

// TestDeleteUser_Success menguji penghapusan user berhasil
func TestDeleteUser_Success(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			if id == 1 {
				return existingUser, nil
			}
			return nil, errors.New("not found")
		},
		// Delete stub method will return nil by default
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	err := userService.DeleteUser(1)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestDeleteUser_NotFound menguji penghapusan user yang tidak ditemukan
func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return nil, errors.New("record not found")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	err := userService.DeleteUser(999)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "user tidak ditemukan" {
		t.Errorf("Expected error message 'user tidak ditemukan', got '%s'", err.Error())
	}
}

// TestDeleteUser_RepositoryError menguji ketika repository gagal delete
func TestDeleteUser_RepositoryError(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	var deleteCalled bool

	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return existingUser, nil
		},
		DeleteFunc: func(id int) error {
			deleteCalled = true
			return errors.New("database connection failed")
		},
	}
	mockCloudinary := &MockCloudinaryServiceForUserService{}

	userService := NewUserService(mockRepo, mockCloudinary)
	err := userService.DeleteUser(1)

	// Validasi
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "gagal menghapus user" {
		t.Errorf("Expected error message 'gagal menghapus user', got '%s'", err.Error())
	}
	if !deleteCalled {
		t.Error("Expected Delete method to be called")
	}
}
