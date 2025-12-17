package service

import (
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/logger"
)

type HomeService interface {
	GetHeroSection() ([]responses.HeroSectionResponse, error)
}

type homeService struct {
	homeRepository    repository.HomeRepository
	cloudinaryService CloudinaryService
}

func NewPublicHomeService(homeRepository repository.HomeRepository, cloudinaryService CloudinaryService) HomeService {
	return &homeService{
		homeRepository:    homeRepository,
		cloudinaryService: cloudinaryService,
	}
}

func (s *homeService) GetHeroSection() ([]responses.HeroSectionResponse, error) {
	heroSection, err := s.homeRepository.GetHeroSection()
	if err != nil {
		logger.Error.Printf("Failed to get hero section from repository: %v", err)
		return nil, err
	}

	return heroSection, nil
}
