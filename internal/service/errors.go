package service

import "errors"

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
