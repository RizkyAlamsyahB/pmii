package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// AboutService interface untuk business logic about page
type AboutService interface {
	Get(ctx context.Context) (*responses.AboutResponse, error)
	Update(ctx context.Context, req requests.UpdateAboutRequest, imageFile *multipart.FileHeader) (*responses.AboutResponse, error)
}

type aboutService struct {
	aboutRepo         repository.AboutRepository
	cloudinaryService CloudinaryService
}

// NewAboutService constructor untuk AboutService
func NewAboutService(aboutRepo repository.AboutRepository, cloudinaryService CloudinaryService) AboutService {
	return &aboutService{
		aboutRepo:         aboutRepo,
		cloudinaryService: cloudinaryService,
	}
}

// Get mengambil data about page
func (s *aboutService) Get(ctx context.Context) (*responses.AboutResponse, error) {
	about, err := s.aboutRepo.Get()
	if err != nil {
		// Jika belum ada data, return empty response
		return &responses.AboutResponse{}, nil
	}

	return s.toResponseDTO(about), nil
}

// Update mengupdate about page dengan optional upload image baru
func (s *aboutService) Update(ctx context.Context, req requests.UpdateAboutRequest, imageFile *multipart.FileHeader) (*responses.AboutResponse, error) {
	// Ambil about existing (jika ada)
	about, err := s.aboutRepo.Get()
	if err != nil {
		// Belum ada data, buat baru
		about = &domain.About{}
	}

	// Simpan image lama untuk rollback
	oldImageURI := about.ImageURI

	// Upload image baru ke Cloudinary (jika ada)
	var newImageFilename *string
	if imageFile != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "about", imageFile)
		if err != nil {
			return nil, errors.New("gagal mengupload gambar")
		}
		newImageFilename = &filename
		about.ImageURI = &filename
	}

	// Update fields yang dikirim
	if req.History != "" {
		about.History = &req.History
	}
	if req.Vision != "" {
		about.Vision = &req.Vision
	}
	if req.Mission != "" {
		about.Mission = &req.Mission
	}
	if req.VideoURL != "" {
		about.VideoURL = &req.VideoURL
	}

	// Save ke database (upsert)
	if err := s.aboutRepo.Upsert(about); err != nil {
		// Rollback: hapus image baru jika update gagal
		if newImageFilename != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "about", *newImageFilename)
		}
		return nil, errors.New("gagal menyimpan about")
	}

	// Hapus image lama SETELAH database update berhasil
	if newImageFilename != nil && oldImageURI != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "about", *oldImageURI)
	}

	return s.toResponseDTO(about), nil
}

// toResponseDTO converts domain.About to responses.AboutResponse
func (s *aboutService) toResponseDTO(a *domain.About) *responses.AboutResponse {
	var imageURL string
	if a.ImageURI != nil {
		imageURL = s.cloudinaryService.GetImageURL("about", *a.ImageURI)
	}

	return &responses.AboutResponse{
		ID:        a.ID,
		History:   a.History,
		Vision:    a.Vision,
		Mission:   a.Mission,
		ImageUrl:  imageURL,
		VideoURL:  a.VideoURL,
		UpdatedAt: a.UpdatedAt,
	}
}
