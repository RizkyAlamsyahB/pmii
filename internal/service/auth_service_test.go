package service

import (
	"errors"
	"testing"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	FindByEmailFunc func(email string) (*domain.User, error)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, errors.New("mock not configured")
}

// Stub methods (interface requirement)
func (m *MockUserRepository) Create(user *domain.User) error        { return nil }
func (m *MockUserRepository) Update(user *domain.User) error        { return nil }
func (m *MockUserRepository) Delete(id int) error                   { return nil }
func (m *MockUserRepository) FindByID(id int) (*domain.User, error) { return nil, nil }
func (m *MockUserRepository) FindAll() ([]domain.User, error)       { return nil, nil }

// TestLogin_UserNotFound menguji login dengan email yang tidak terdaftar
func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}

	authService := NewAuthService(mockRepo)
	user, token, err := authService.Login("notfound@example.com", "password123")

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

	authService := NewAuthService(mockRepo)
	user, token, err := authService.Login("test@example.com", "wrongpassword")

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

	authService := NewAuthService(mockRepo)
	user, token, err := authService.Login("test@example.com", "admin123")

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
