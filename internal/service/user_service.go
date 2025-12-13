package service

import (
	"context"
	"errors"
	"mime/multipart"
	"regexp"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// UserService interface untuk business logic user operations
type UserService interface {
	GetAllUsers(page, limit int) ([]domain.User, int, int, int64, error)
	GetUserByID(id int) (*domain.User, error)
	CreateUser(ctx context.Context, req *requests.CreateUserRequest, photoFile *multipart.FileHeader) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, req *requests.UpdateUserRequest, photoFile *multipart.FileHeader) (*domain.User, error)
	DeleteUser(id int) error
}

type userService struct {
	userRepo          repository.UserRepository
	cloudinaryService CloudinaryService
}

// NewUserService constructor untuk UserService
func NewUserService(userRepo repository.UserRepository, cloudinaryService CloudinaryService) UserService {
	return &userService{userRepo: userRepo, cloudinaryService: cloudinaryService}
}

// GetAllUsers mengambil semua user dari database dengan pagination
func (s *userService) GetAllUsers(page, limit int) ([]domain.User, int, int, int64, error) {
	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	users, total, err := s.userRepo.FindAll(page, limit)
	if err != nil {
		return nil, 0, 0, 0, errors.New("gagal mengambil data user")
	}

	// Calculate last page
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	return users, page, lastPage, total, nil
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
func (s *userService) CreateUser(ctx context.Context, req *requests.CreateUserRequest, photoFile *multipart.FileHeader) (*domain.User, error) {
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

	// Upload photo ke cloudinary (jika ada)
	var photoFileName *string
	if photoFile != nil {
		fileName, err := s.cloudinaryService.UploadImage(ctx, "users/avatars", photoFile)
		if err != nil {
			return nil, errors.New("gagal mengupload foto")
		}
		photoFileName = &fileName
	}

	// Buat user baru dengan role default Author (2)
	user := &domain.User{
		Role:         2, // Default: Author
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		PhotoURI:     photoFileName,
		IsActive:     true,
	}

	// Simpan ke database
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("gagal membuat user")
	}

	// Gunakan full URL untuk photo di response (jika ada foto)
	if photoFileName != nil {
		fullPhotoURL := s.cloudinaryService.GetImageURL("users/avatars", *photoFileName)
		user.PhotoURI = &fullPhotoURL
	}

	return user, nil
}

// UpdateUser mengupdate data user berdasarkan ID (Admin only)
func (s *userService) UpdateUser(ctx context.Context, id int, req *requests.UpdateUserRequest, photoFile *multipart.FileHeader) (*domain.User, error) {
	// Cek apakah user ada
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	oldPhoto := user.PhotoURI

	// Upload photo baru ke cloudinary (jika ada)
	var newPhotoFileName *string
	if photoFile != nil {
		fileName, err := s.cloudinaryService.UploadImage(ctx, "users/avatars", photoFile)
		if err != nil {
			return nil, errors.New("gagal mengupload foto")
		}
		newPhotoFileName = &fileName
	}

	// Cek apakah email sudah dipakai user lain (hanya jika email diubah)
	if req.Email != nil && *req.Email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(*req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email sudah digunakan user lain")
		}
	}

	// Jika password diisi, validasi dan hash
	if req.Password != nil && *req.Password != "" {
		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(*req.Password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(*req.Password)
		if !hasLetter || !hasNumber {
			return nil, errors.New("password harus kombinasi huruf dan angka")
		}

		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, errors.New("gagal memproses password")
		}
		user.PasswordHash = hashedPassword
	}

	// Update fields (hanya field yang ada di request)
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if newPhotoFileName != nil {
		user.PhotoURI = newPhotoFileName
	}
	// Simpan ke database
	if err := s.userRepo.Update(user); err != nil {
		// hapus photo baru jika update gagal (dan ada foto baru)
		if newPhotoFileName != nil {
			s.cloudinaryService.DeleteImage(ctx, "users/avatars", *newPhotoFileName)
		}
		return nil, errors.New("gagal mengupdate user")
	}

	// hapus photo lama (hanya jika ada foto baru yang diupload)
	if newPhotoFileName != nil && oldPhoto != nil {
		s.cloudinaryService.DeleteImage(ctx, "users/avatars", *oldPhoto)
	}

	// return full image URL in response (jika ada foto baru)
	if newPhotoFileName != nil {
		fullImageURL := s.cloudinaryService.GetImageURL("users/avatars", *newPhotoFileName)
		user.PhotoURI = &fullImageURL
	}

	return user, nil
}

// DeleteUser menghapus user berdasarkan ID (soft delete)
func (s *userService) DeleteUser(id int) error {
	// Cek apakah user ada
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Hapus user (soft delete via GORM)
	if err := s.userRepo.Delete(id); err != nil {
		return errors.New("gagal menghapus user")
	}

	return nil
}
