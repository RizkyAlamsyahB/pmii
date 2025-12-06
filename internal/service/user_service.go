package service

import (
	"errors"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// UserService interface untuk business logic user operations
type UserService interface {
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService constructor untuk UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// GetAllUsers mengambil semua user dari database
func (s *userService) GetAllUsers() ([]domain.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, errors.New("gagal mengambil data user")
	}
	return users, nil
}

// GetUserByID mengambil user berdasarkan ID
func (s *userService) GetUserByID(id uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}
