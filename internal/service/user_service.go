package service

import (
	"errors"
	"regexp"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// UserService interface untuk business logic user operations
type UserService interface {
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	CreateUser(req *requests.CreateUserRequest) (*domain.User, error)
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
func (s *userService) GetUserByID(id int) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

// CreateUser membuat user baru
func (s *userService) CreateUser(req *requests.CreateUserRequest) (*domain.User, error) {
	// Validasi password harus kombinasi huruf dan angka
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(req.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(req.Password)
	if !hasLetter || !hasNumber {
		return nil, errors.New("password harus kombinasi huruf dan angka")
	}

	// Cek apakah email sudah terdaftar
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	// Buat user baru dengan role default Author (2)
	var photoURI *string
	if req.PhotoURI != "" {
		photoURI = &req.PhotoURI
	}

	user := &domain.User{
		Role:         2, // Default: Author
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		PhotoURI:     photoURI,
		IsActive:     true,
	}

	// Simpan ke database
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("gagal membuat user")
	}

	return user, nil
}
