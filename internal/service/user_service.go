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
	GetAllUsers(ctx context.Context, page, limit int) ([]domain.User, int, int, int64, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	CreateUser(ctx context.Context, req *requests.CreateUserRequest, photoFile *multipart.FileHeader) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, req *requests.UpdateUserRequest, photoFile *multipart.FileHeader) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	userRepo          repository.UserRepository
	cloudinaryService CloudinaryService
}

// NewUserService constructor untuk UserService
func NewUserService(userRepo repository.UserRepository, cloudinaryService CloudinaryService) UserService {
	return &userService{userRepo: userRepo, cloudinaryService: cloudinaryService}
}

// resolvePhotoURL converts photo filename to full Cloudinary URL
func (s *userService) resolvePhotoURL(user *domain.User) {
	if user.PhotoURI != nil && *user.PhotoURI != "" {
		fullURL := s.cloudinaryService.GetImageURL("users/avatars", *user.PhotoURI)
		user.PhotoURI = &fullURL
	}
}

// GetAllUsers mengambil semua user dari database dengan pagination
func (s *userService) GetAllUsers(ctx context.Context, page, limit int) ([]domain.User, int, int, int64, error) {
	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	users, total, err := s.userRepo.FindAll(page, limit)
	if err != nil {
		return nil, 0, 0, 0, ErrUserFetchFailed
	}

	// Calculate last page
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	// Resolve photo URLs untuk semua user
	for i := range users {
		s.resolvePhotoURL(&users[i])
	}

	return users, page, lastPage, total, nil
}

// GetUserByID mengambil user berdasarkan ID
func (s *userService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	s.resolvePhotoURL(user)
	return user, nil
}

// CreateUser membuat user baru
func (s *userService) CreateUser(ctx context.Context, req *requests.CreateUserRequest, photoFile *multipart.FileHeader) (*domain.User, error) {
	// Validasi password harus kombinasi huruf dan angka
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(req.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(req.Password)
	if !hasLetter || !hasNumber {
		return nil, ErrInvalidPassword
	}

	// Cek apakah email sudah terdaftar
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, ErrPasswordProcessing
	}

	// Upload photo ke cloudinary (jika ada)
	var photoFileName *string
	if photoFile != nil {
		fileName, err := s.cloudinaryService.UploadImage(ctx, "users/avatars", photoFile)
		if err != nil {
			return nil, ErrPhotoUploadFailed
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
		// Rollback: hapus foto dari Cloudinary jika save gagal
		if photoFileName != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "users/avatars", *photoFileName)
		}
		return nil, ErrUserCreateFailed
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
		return nil, ErrUserNotFound
	}

	oldPhoto := user.PhotoURI

	// Upload photo baru ke cloudinary (jika ada)
	var newPhotoFileName *string
	if photoFile != nil {
		fileName, err := s.cloudinaryService.UploadImage(ctx, "users/avatars", photoFile)
		if err != nil {
			return nil, ErrPhotoUploadFailed
		}
		newPhotoFileName = &fileName
	}

	// Cek apakah email sudah dipakai user lain (hanya jika email diubah)
	if req.Email != nil && *req.Email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(*req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailAlreadyUsed
		}
	}

	// Jika password diisi, validasi dan hash
	if req.Password != nil && *req.Password != "" {
		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(*req.Password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(*req.Password)
		if !hasLetter || !hasNumber {
			return nil, ErrInvalidPassword
		}

		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, ErrPasswordProcessing
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
			_ = s.cloudinaryService.DeleteImage(ctx, "users/avatars", *newPhotoFileName)
		}
		return nil, ErrUserUpdateFailed
	}

	// hapus photo lama (hanya jika ada foto baru yang diupload)
	if newPhotoFileName != nil && oldPhoto != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "users/avatars", *oldPhoto)
	}

	// Resolve photo URL untuk response
	if newPhotoFileName != nil {
		fullImageURL := s.cloudinaryService.GetImageURL("users/avatars", *newPhotoFileName)
		user.PhotoURI = &fullImageURL
	} else {
		// Resolve existing photo URL jika tidak ada foto baru
		s.resolvePhotoURL(user)
	}

	return user, nil
}

// DeleteUser menghapus user berdasarkan ID (soft delete)
func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// Cek apakah user ada dan ambil info foto
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	// Hapus user (soft delete via GORM)
	if err := s.userRepo.Delete(id); err != nil {
		return ErrUserDeleteFailed
	}

	// Foto tidak dihapus dari cloudinary karena soft delete

	return nil
}

// User service errors
var (
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
	ErrEmailAlreadyUsed   = errors.New("email sudah digunakan user lain")
	ErrInvalidPassword    = errors.New("password harus kombinasi huruf dan angka")
	ErrPasswordProcessing = errors.New("gagal memproses password")
	ErrPhotoUploadFailed  = errors.New("gagal mengupload foto")
	ErrUserCreateFailed   = errors.New("gagal membuat user")
	ErrUserUpdateFailed   = errors.New("gagal mengupdate user")
	ErrUserDeleteFailed   = errors.New("gagal menghapus user")
	ErrUserFetchFailed    = errors.New("gagal mengambil data user")
)
