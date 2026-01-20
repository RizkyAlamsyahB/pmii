package service

import (
	"context"
	"errors"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// PublicAdService interface untuk public ads (tanpa auth)
type PublicAdService interface {
	GetAdsByPage(ctx context.Context, page string) ([]responses.PublicAdResponse, error)
}

type publicAdService struct {
	adRepo     repository.AdRepository
	cloudinary CloudinaryService
}

// NewPublicAdService constructor untuk PublicAdService
func NewPublicAdService(adRepo repository.AdRepository, cloudinaryService CloudinaryService) PublicAdService {
	return &publicAdService{
		adRepo:     adRepo,
		cloudinary: cloudinaryService,
	}
}

// GetAdsByPage mengambil ads berdasarkan page untuk publik
func (s *publicAdService) GetAdsByPage(ctx context.Context, page string) ([]responses.PublicAdResponse, error) {
	if !domain.IsValidAdPage(page) {
		return nil, errors.New("halaman tidak valid")
	}

	ads, err := s.adRepo.FindByPage(domain.AdPage(page))
	if err != nil {
		return nil, err
	}

	return s.toPublicAdResponseListWithFullURL(ads), nil
}

// toPublicAdResponseWithFullURL converts domain.Ad to PublicAdResponse with full cloudinary URL
func (s *publicAdService) toPublicAdResponseWithFullURL(ad *domain.Ad) responses.PublicAdResponse {
	var imageURL *string
	if ad.ImageURL != nil && *ad.ImageURL != "" {
		fullURL := s.cloudinary.GetImageURL("ads", *ad.ImageURL)
		imageURL = &fullURL
	}

	return responses.PublicAdResponse{
		Slot:     ad.Slot,
		ImageURL: imageURL,
	}
}

// toPublicAdResponseListWithFullURL converts slice of domain.Ad to slice of PublicAdResponse with full URLs
func (s *publicAdService) toPublicAdResponseListWithFullURL(ads []domain.Ad) []responses.PublicAdResponse {
	result := make([]responses.PublicAdResponse, len(ads))
	for i, ad := range ads {
		result[i] = s.toPublicAdResponseWithFullURL(&ad)
	}
	return result
}
