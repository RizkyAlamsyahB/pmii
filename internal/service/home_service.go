package service

import (
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/logger"
)

type HomeService interface {
	GetHeroSection() ([]responses.HeroSectionResponse, error)
	GetLatestNewsSection() ([]responses.LatestNewsSectionResponse, error)
	GetAboutUsSection() (*responses.AboutUsSectionResponse, error)
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

	for i := range heroSection {
		if heroSection[i].FeaturedImage != "" {
			heroSection[i].FeaturedImage = s.cloudinaryService.GetImageURL("posts/images", heroSection[i].FeaturedImage)
		}
	}

	return heroSection, nil
}

func (s *homeService) GetLatestNewsSection() ([]responses.LatestNewsSectionResponse, error) {
	latestNewsSection, err := s.homeRepository.GetLatestNewsSection()
	if err != nil {
		logger.Error.Printf("Failed to get latest news section from repository: %v", err)
		return nil, err
	}

	for i := range latestNewsSection {
		if latestNewsSection[i].FeaturedImage != "" {
			latestNewsSection[i].FeaturedImage = s.cloudinaryService.GetImageURL("posts/images", latestNewsSection[i].FeaturedImage)
		}
	}

	return latestNewsSection, nil
}

func (s *homeService) GetAboutUsSection() (*responses.AboutUsSectionResponse, error) {
	aboutUsSection, err := s.homeRepository.GetAboutUsSection()
	if err != nil {
		logger.Error.Printf("Failed to get about us section from repository: %v", err)
		return nil, err
	}

	if aboutUsSection.ImageURI != "" {
		aboutUsSection.ImageURI = s.cloudinaryService.GetImageURL("home/images", aboutUsSection.ImageURI)
	}

	return aboutUsSection, nil
}
