package service

import (
	"context"
	"errors"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// AboutService interface untuk business logic about page
type AboutService interface {
	Get(ctx context.Context) (*responses.AboutResponse, error)
	Update(ctx context.Context, req requests.UpdateAboutRequest) (*responses.AboutResponse, error)
}

type aboutService struct {
	aboutRepo repository.AboutRepository
}

// NewAboutService constructor untuk AboutService
func NewAboutService(aboutRepo repository.AboutRepository) AboutService {
	return &aboutService{
		aboutRepo: aboutRepo,
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

// Update mengupdate about page
func (s *aboutService) Update(ctx context.Context, req requests.UpdateAboutRequest) (*responses.AboutResponse, error) {
	// Ambil about existing (jika ada)
	about, err := s.aboutRepo.Get()
	if err != nil {
		// Belum ada data, buat baru
		about = &domain.About{}
	}

	// Update fields yang dikirim
	if req.Title != "" {
		about.Title = &req.Title
	}
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
		return nil, errors.New("gagal menyimpan about")
	}

	return s.toResponseDTO(about), nil
}

// toResponseDTO converts domain.About to responses.AboutResponse
func (s *aboutService) toResponseDTO(a *domain.About) *responses.AboutResponse {
	return &responses.AboutResponse{
		ID:       a.ID,
		Title:    a.Title,
		History:  a.History,
		Vision:   a.Vision,
		Mission:  a.Mission,
		VideoURL: a.VideoURL,
	}
}
