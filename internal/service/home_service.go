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
	GetWhySection() (*responses.WhySectionResponse, error)
	GetTestimonialSection() ([]responses.TestimonialSectionResponse, error)
	GetFaqSection() (*responses.FaqSectionResponse, error)
	GetCtaSection() (*responses.CtaSectionResponse, error)
}

type homeService struct {
	homeRepository        repository.HomeRepository
	testimonialRepository repository.TestimonialRepository
	cloudinaryService     CloudinaryService
}

func NewPublicHomeService(homeRepository repository.HomeRepository, testimonialRepository repository.TestimonialRepository, cloudinaryService CloudinaryService) HomeService {
	return &homeService{
		homeRepository:        homeRepository,
		testimonialRepository: testimonialRepository,
		cloudinaryService:     cloudinaryService,
	}
}

func (s *homeService) GetHeroSection() ([]responses.HeroSectionResponse, error) {
	heroSection, err := s.homeRepository.GetHeroSection()
	if err != nil {
		logger.Error.Printf("Failed to get hero section from repository: %v", err)
		return nil, err
	}

	for i := range heroSection {
		if heroSection[i].ImageURL != "" {
			heroSection[i].ImageURL = s.cloudinaryService.GetImageURL("posts/images", heroSection[i].ImageURL)
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

func (s *homeService) GetWhySection() (*responses.WhySectionResponse, error) {
	whySection, err := s.homeRepository.GetWhySection()
	if err != nil {
		logger.Error.Printf("Failed to get why section from repository: %v", err)
		return nil, err
	}

	for i := range whySection.Data {
		if whySection.Data[i].IconURI != "" {
			whySection.Data[i].IconURI = s.cloudinaryService.GetImageURL("why/images", whySection.Data[i].IconURI)
		}
	}

	return whySection, nil
}

func (s *homeService) GetTestimonialSection() ([]responses.TestimonialSectionResponse, error) {
	testimonials, _, err := s.testimonialRepository.FindAll(1, 7)
	if err != nil {
		logger.Error.Printf("Failed to get testimonial section from repository: %v", err)
		return nil, err
	}

	var result []responses.TestimonialSectionResponse
	for _, testimonial := range testimonials {
		item := responses.TestimonialSectionResponse{
			Testimoni: testimonial.Content,
			Name:      testimonial.Name,
		}

		// Nil check for Position field
		if testimonial.Position != nil {
			item.Status = *testimonial.Position
		} else {
			item.Status = ""
		}

		if testimonial.Organization != nil {
			item.Career = *testimonial.Organization
		} else {
			item.Career = ""
		}

		// Nil check for PhotoURI field
		if testimonial.PhotoURI != nil {
			item.ImageURI = s.cloudinaryService.GetImageURL("home/images", *testimonial.PhotoURI)
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *homeService) GetFaqSection() (*responses.FaqSectionResponse, error) {
	faqSection, err := s.homeRepository.GetFaqSection()
	if err != nil {
		logger.Error.Printf("Failed to get faq section from repository: %v", err)
		return nil, err
	}

	return faqSection, nil
}

func (s *homeService) GetCtaSection() (*responses.CtaSectionResponse, error) {
	ctaSection, err := s.homeRepository.GetCtaSection()
	if err != nil {
		logger.Error.Printf("Failed to get cta section from repository: %v", err)
		return nil, err
	}

	return ctaSection, nil
}
