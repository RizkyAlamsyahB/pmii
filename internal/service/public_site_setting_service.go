package service

import (
	"context"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// PublicSiteSettingService interface untuk public site settings
type PublicSiteSettingService interface {
	Get(ctx context.Context) (*responses.PublicSiteSettingResponse, error)
}

type publicSiteSettingService struct {
	siteSettingRepo   repository.SiteSettingRepository
	cloudinaryService CloudinaryService
}

// NewPublicSiteSettingService constructor untuk PublicSiteSettingService
func NewPublicSiteSettingService(
	siteSettingRepo repository.SiteSettingRepository,
	cloudinaryService CloudinaryService,
) PublicSiteSettingService {
	return &publicSiteSettingService{
		siteSettingRepo:   siteSettingRepo,
		cloudinaryService: cloudinaryService,
	}
}

// Get mengambil site settings untuk public
func (s *publicSiteSettingService) Get(ctx context.Context) (*responses.PublicSiteSettingResponse, error) {
	setting, err := s.siteSettingRepo.Get()
	if err != nil {
		// Return empty response if not found
		return &responses.PublicSiteSettingResponse{}, nil
	}

	return s.toResponseDTO(setting), nil
}

// toResponseDTO converts domain.SiteSetting to responses.PublicSiteSettingResponse
func (s *publicSiteSettingService) toResponseDTO(setting *domain.SiteSetting) *responses.PublicSiteSettingResponse {
	response := &responses.PublicSiteSettingResponse{
		SiteName:        setting.SiteName,
		SiteTitle:       setting.SiteTitle,
		SiteDescription: setting.SiteDescription,
		FacebookURL:     setting.FacebookURL,
		TwitterURL:      setting.TwitterURL,
		LinkedinURL:     setting.LinkedinURL,
		InstagramURL:    setting.InstagramURL,
		YoutubeURL:      setting.YoutubeURL,
		GithubURL:       setting.GithubURL,
	}

	// Resolve image URLs
	if setting.Favicon != nil && *setting.Favicon != "" {
		response.Favicon = s.cloudinaryService.GetImageURL("settings", *setting.Favicon)
	}
	if setting.LogoHeader != nil && *setting.LogoHeader != "" {
		response.LogoHeader = s.cloudinaryService.GetImageURL("settings", *setting.LogoHeader)
	}
	if setting.LogoBig != nil && *setting.LogoBig != "" {
		response.LogoBig = s.cloudinaryService.GetImageURL("settings", *setting.LogoBig)
	}

	return response
}
