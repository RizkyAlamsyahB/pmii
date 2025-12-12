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

// Stub methods (interface requirement)
func (m *MockUserRepositoryForUserService) Update(user *domain.User) error { return nil }
func (m *MockUserRepositoryForUserService) Delete(id int) error            { return nil }

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
