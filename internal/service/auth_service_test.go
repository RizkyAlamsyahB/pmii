package service

import (
	"context"
	"errors"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	FindByEmailFunc func(email string) (*domain.User, error)
	FindByIDFunc    func(id int) (*domain.User, error)
	UpdateFunc      func(user *domain.User) error
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockUserRepository) FindByID(id int) (*domain.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockUserRepository) Update(user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

// Stub methods (interface requirement)
func (m *MockUserRepository) Create(user *domain.User) error { return nil }
func (m *MockUserRepository) Delete(id int) error            { return nil }
func (m *MockUserRepository) FindAll(page, limit int) ([]domain.User, int64, error) {
	return nil, 0, nil
}

// MockActivityLogRepoForAuth adalah mock untuk ActivityLogRepository
type MockActivityLogRepoForAuth struct {
	CreateFunc func(log *domain.ActivityLog) error
}

func (m *MockActivityLogRepoForAuth) Create(log *domain.ActivityLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(log)
	}
	return nil
}

func (m *MockActivityLogRepoForAuth) GetActivityLogs(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
	return nil, 0, nil
}

// TestLogin_UserNotFound menguji login dengan email yang tidak terdaftar
func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	user, token, err := authService.Login(context.Background(), "notfound@example.com", "password123")

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
	if token != "" {
		t.Errorf("Expected empty token, got %s", token)
	}
}

// TestLogin_WrongPassword menguji login dengan password yang salah
func TestLogin_WrongPassword(t *testing.T) {
	// Hash untuk "admin123"
	hashedPassword := "$2a$10$CaI1bA6w2H0LKVGF./iweOteBj/rAkkpx3QUO/dDK5.dRP6BDeB8a"

	mockRepo := &MockUserRepository{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hashedPassword,
				Role:         1,
				IsActive:     true,
			}, nil
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	user, token, err := authService.Login(context.Background(), "test@example.com", "wrongpassword")

	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
	if token != "" {
		t.Errorf("Expected empty token, got %s", token)
	}
}

// TestLogin_InactiveUser menguji login dengan user yang tidak aktif
func TestLogin_InactiveUser(t *testing.T) {
	hashedPassword := "$2a$10$CaI1bA6w2H0LKVGF./iweOteBj/rAkkpx3QUO/dDK5.dRP6BDeB8a"

	mockRepo := &MockUserRepository{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hashedPassword,
				Role:         1,
				IsActive:     false,
			}, nil
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	user, token, err := authService.Login(context.Background(), "test@example.com", "admin123")

	if err == nil {
		t.Error("Expected error for inactive user, got nil")
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
	if token != "" {
		t.Errorf("Expected empty token, got %s", token)
	}
}

// TestChangePassword_UserNotFound menguji change password dengan user yang tidak ditemukan
func TestChangePassword_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	req := requests.ChangePasswordRequest{
		OldPassword: "oldpass123!",
		NewPassword: "newpass123!",
	}

	err := authService.ChangePassword(999, req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "user tidak ditemukan" {
		t.Errorf("Expected 'user tidak ditemukan', got %s", err.Error())
	}
}

// TestChangePassword_InactiveUser menguji change password dengan user yang tidak aktif
func TestChangePassword_InactiveUser(t *testing.T) {
	hashedPassword := "$2a$10$CaI1bA6w2H0LKVGF./iweOteBj/rAkkpx3QUO/dDK5.dRP6BDeB8a"

	mockRepo := &MockUserRepository{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hashedPassword,
				Role:         1,
				IsActive:     false,
			}, nil
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	req := requests.ChangePasswordRequest{
		OldPassword: "admin123",
		NewPassword: "newpass123!",
	}

	err := authService.ChangePassword(1, req)

	if err == nil {
		t.Error("Expected error for inactive user, got nil")
	}
	if err.Error() != "user tidak aktif" {
		t.Errorf("Expected 'user tidak aktif', got %s", err.Error())
	}
}

// TestChangePassword_WrongOldPassword menguji change password dengan password lama yang salah
func TestChangePassword_WrongOldPassword(t *testing.T) {
	hashedPassword := "$2a$10$CaI1bA6w2H0LKVGF./iweOteBj/rAkkpx3QUO/dDK5.dRP6BDeB8a"

	mockRepo := &MockUserRepository{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hashedPassword,
				Role:         1,
				IsActive:     true,
			}, nil
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	req := requests.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpass123!",
	}

	err := authService.ChangePassword(1, req)

	if err == nil {
		t.Error("Expected error for wrong old password, got nil")
	}
	if err.Error() != "password lama salah" {
		t.Errorf("Expected 'password lama salah', got %s", err.Error())
	}
}

// TestChangePassword_SamePassword menguji change password dengan password baru yang sama dengan password lama
func TestChangePassword_SamePassword(t *testing.T) {
	hashedPassword := "$2a$10$CaI1bA6w2H0LKVGF./iweOteBj/rAkkpx3QUO/dDK5.dRP6BDeB8a"

	mockRepo := &MockUserRepository{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hashedPassword,
				Role:         1,
				IsActive:     true,
			}, nil
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	req := requests.ChangePasswordRequest{
		OldPassword: "admin123",
		NewPassword: "admin123",
	}

	err := authService.ChangePassword(1, req)

	if err == nil {
		t.Error("Expected error for same password, got nil")
	}
	if err.Error() != "password baru tidak boleh sama dengan password lama" {
		t.Errorf("Expected 'password baru tidak boleh sama dengan password lama', got %s", err.Error())
	}
}

// TestChangePassword_Success menguji change password yang berhasil
func TestChangePassword_Success(t *testing.T) {
	hashedPassword := "$2a$10$CaI1bA6w2H0LKVGF./iweOteBj/rAkkpx3QUO/dDK5.dRP6BDeB8a"
	updateCalled := false

	mockRepo := &MockUserRepository{
		FindByIDFunc: func(id int) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hashedPassword,
				Role:         1,
				IsActive:     true,
			}, nil
		},
		UpdateFunc: func(user *domain.User) error {
			updateCalled = true
			if user.PasswordHash == hashedPassword {
				t.Error("Expected password hash to be changed")
			}
			return nil
		},
	}

	authService := NewAuthService(mockRepo, &MockActivityLogRepoForAuth{})
	req := requests.ChangePasswordRequest{
		OldPassword: "admin123",
		NewPassword: "newpass123!",
	}

	err := authService.ChangePassword(1, req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !updateCalled {
		t.Error("Expected Update to be called")
	}
}
