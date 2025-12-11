package service

import (
	"errors"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// AuthService interface untuk business logic authentication
type AuthService interface {
	Login(email, password string) (*domain.User, string, error)
	Logout(token string) error
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService constructor untuk AuthService
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Login melakukan proses login user
func (s *authService) Login(email, password string) (*domain.User, string, error) {
	// 1. Cari user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// 2. Verify password dengan bcrypt
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	// 3. Cek status user aktif
	if !user.IsActive {
		return nil, "", errors.New("user account is inactive")
	}

	// 4. Generate JWT token (convert role int to string)
	token, err := utils.GenerateJWT(user.ID, strconv.Itoa(user.Role))
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// Logout melakukan proses logout user dengan blacklist token
func (s *authService) Logout(token string) error {
	// Validate token untuk get expiry time
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return errors.New("invalid token")
	}

	// Add token to blacklist dengan expiry time
	utils.AddToBlacklist(token, claims.ExpiresAt.Time)

	return nil
}
