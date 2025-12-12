package service

import (
	"errors"
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

// Stub methods (interface requirement)
func (m *MockUserRepositoryForUserService) Delete(id int) error { return nil }

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

	userService := NewUserService(mockRepo)
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

	userService := NewUserService(mockRepo)
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

	userService := NewUserService(mockRepo)
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

	userService := NewUserService(mockRepo)
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

	userService := NewUserService(mockRepo)
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

	userService := NewUserService(mockRepo)
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

// TestCreateUser_Success menguji pembuatan user berhasil
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

	userService := NewUserService(mockRepo)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123", // Kombinasi huruf dan angka
		PhotoURI: "",
	}

	user, err := userService.CreateUser(req)

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

	userService := NewUserService(mockRepo)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "passwordonly", // Hanya huruf, tanpa angka
	}

	user, err := userService.CreateUser(req)

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

	userService := NewUserService(mockRepo)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "12345678", // Hanya angka, tanpa huruf
	}

	user, err := userService.CreateUser(req)

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

	userService := NewUserService(mockRepo)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	user, err := userService.CreateUser(req)

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

	userService := NewUserService(mockRepo)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, err := userService.CreateUser(req)

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

// TestCreateUser_WithPhotoURI menguji pembuatan user dengan photo URI
func TestCreateUser_WithPhotoURI(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		CreateFunc: func(user *domain.User) error {
			user.ID = 1
			return nil
		},
	}

	userService := NewUserService(mockRepo)
	req := &requests.CreateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
		PhotoURI: "https://example.com/photo.jpg",
	}

	user, err := userService.CreateUser(req)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.PhotoURI == nil {
		t.Error("Expected photo URI, got nil")
	}
	if *user.PhotoURI != "https://example.com/photo.jpg" {
		t.Errorf("Expected photo URI 'https://example.com/photo.jpg', got '%s'", *user.PhotoURI)
	}
}

// ============================================================
// Test Cases untuk UpdateUser
// ============================================================

// TestUpdateUser_Success menguji update user berhasil
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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "New Name",
		Email:    "new@example.com",
		Role:     1,
		PhotoURI: "https://example.com/new-photo.jpg",
		IsActive: false,
	}

	user, err := userService.UpdateUser(1, req)

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
	if user.PhotoURI == nil || *user.PhotoURI != "https://example.com/new-photo.jpg" {
		t.Error("Expected photo URI to be updated")
	}
}

// TestUpdateUser_NotFound menguji update user yang tidak ditemukan
func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := &MockUserRepositoryForUserService{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return nil, errors.New("record not found")
		},
	}

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "New Name",
		Email:    "new@example.com",
		Role:     1,
		IsActive: true,
	}

	user, err := userService.UpdateUser(999, req)

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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "User One Updated",
		Email:    "user2@example.com", // Mencoba pakai email user lain
		Role:     2,
		IsActive: true,
	}

	user, err := userService.UpdateUser(1, req)

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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "User One Updated",
		Email:    "user1@example.com", // Email sama, tidak berubah
		Role:     1,
		IsActive: true,
	}

	user, err := userService.UpdateUser(1, req)

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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
		Password: "onlyletters", // Hanya huruf
	}

	user, err := userService.UpdateUser(1, req)

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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
		Password: "12345678", // Hanya angka
	}

	user, err := userService.UpdateUser(1, req)

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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
		Password: "newpassword123", // Password valid
	}

	user, err := userService.UpdateUser(1, req)

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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "Updated Name",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
		Password: "", // Password kosong, tidak diubah
	}

	user, err := userService.UpdateUser(1, req)

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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "Updated Name",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
	}

	user, err := userService.UpdateUser(1, req)

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

// TestUpdateUser_ClearPhotoURI menguji update dengan menghapus photo URI
func TestUpdateUser_ClearPhotoURI(t *testing.T) {
	photoURI := "https://example.com/old-photo.jpg"
	existingUser := &domain.User{
		ID:       1,
		FullName: "Test User",
		Email:    "test@example.com",
		PhotoURI: &photoURI,
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

	userService := NewUserService(mockRepo)
	req := &requests.UpdateUserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     2,
		IsActive: true,
		PhotoURI: "", // Kosongkan photo URI
	}

	user, err := userService.UpdateUser(1, req)

	// Validasi
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.PhotoURI != nil {
		t.Errorf("Expected photo URI to be nil, got %v", user.PhotoURI)
	}
}
