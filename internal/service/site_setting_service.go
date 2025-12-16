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

// SiteSettingService interface untuk business logic site settings
type SiteSettingService interface {
	Get(ctx context.Context) (*responses.SiteSettingResponse, error)
	Update(ctx context.Context, req requests.UpdateSiteSettingRequest, favicon, logoHeader, logoBig *multipart.FileHeader) (*responses.SiteSettingResponse, error)
}

type siteSettingService struct {
	siteSettingRepo   repository.SiteSettingRepository
	cloudinaryService CloudinaryService
}

// NewSiteSettingService constructor untuk SiteSettingService
func NewSiteSettingService(siteSettingRepo repository.SiteSettingRepository, cloudinaryService CloudinaryService) SiteSettingService {
	return &siteSettingService{
		siteSettingRepo:   siteSettingRepo,
		cloudinaryService: cloudinaryService,
	}
}

// Get mengambil site settings
func (s *siteSettingService) Get(ctx context.Context) (*responses.SiteSettingResponse, error) {
	setting, err := s.siteSettingRepo.Get()
	if err != nil {
		// Return empty response if not found
		return &responses.SiteSettingResponse{}, nil
	}

	return s.toResponseDTO(setting), nil
}

// Update mengupdate site settings
func (s *siteSettingService) Update(ctx context.Context, req requests.UpdateSiteSettingRequest, favicon, logoHeader, logoBig *multipart.FileHeader) (*responses.SiteSettingResponse, error) {
	// Get existing settings
	setting, _ := s.siteSettingRepo.Get()
	if setting == nil {
		setting = &domain.SiteSetting{}
	}

	// Track old images for cleanup
	oldFavicon := setting.Favicon
	oldLogoHeader := setting.LogoHeader
	oldLogoBig := setting.LogoBig

	// Track new uploads for rollback
	var newFavicon, newLogoHeader, newLogoBig *string

	// Upload favicon if provided
	if favicon != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "settings", favicon)
		if err != nil {
			return nil, errors.New("gagal mengupload favicon")
		}
		newFavicon = &filename
		setting.Favicon = &filename
	}

	// Upload logo header if provided
	if logoHeader != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "settings", logoHeader)
		if err != nil {
			// Rollback favicon if uploaded
			if newFavicon != nil {
				_ = s.cloudinaryService.DeleteImage(ctx, "settings", *newFavicon)
			}
			return nil, errors.New("gagal mengupload logo header")
		}
		newLogoHeader = &filename
		setting.LogoHeader = &filename
	}

	// Upload logo big if provided
	if logoBig != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "settings", logoBig)
		if err != nil {
			// Rollback previous uploads
			if newFavicon != nil {
				_ = s.cloudinaryService.DeleteImage(ctx, "settings", *newFavicon)
			}
			if newLogoHeader != nil {
				_ = s.cloudinaryService.DeleteImage(ctx, "settings", *newLogoHeader)
			}
			return nil, errors.New("gagal mengupload logo big")
		}
		newLogoBig = &filename
		setting.LogoBig = &filename
	}

	// Update text fields
	if req.SiteName != nil {
		setting.SiteName = req.SiteName
	}
	if req.SiteTitle != nil {
		setting.SiteTitle = req.SiteTitle
	}
	if req.SiteDescription != nil {
		setting.SiteDescription = req.SiteDescription
	}
	if req.FacebookURL != nil {
		setting.FacebookURL = req.FacebookURL
	}
	if req.TwitterURL != nil {
		setting.TwitterURL = req.TwitterURL
	}
	if req.LinkedinURL != nil {
		setting.LinkedinURL = req.LinkedinURL
	}
	if req.InstagramURL != nil {
		setting.InstagramURL = req.InstagramURL
	}
	if req.YoutubeURL != nil {
		setting.YoutubeURL = req.YoutubeURL
	}
	if req.GithubURL != nil {
		setting.GithubURL = req.GithubURL
	}

	// Save to database
	if err := s.siteSettingRepo.Update(setting); err != nil {
		// Rollback all new uploads
		if newFavicon != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "settings", *newFavicon)
		}
		if newLogoHeader != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "settings", *newLogoHeader)
		}
		if newLogoBig != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "settings", *newLogoBig)
		}
		return nil, errors.New("gagal menyimpan pengaturan situs")
	}

	// Delete old images AFTER successful save
	if newFavicon != nil && oldFavicon != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "settings", *oldFavicon)
	}
	if newLogoHeader != nil && oldLogoHeader != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "settings", *oldLogoHeader)
	}
	if newLogoBig != nil && oldLogoBig != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "settings", *oldLogoBig)
	}

	return s.toResponseDTO(setting), nil
}

// toResponseDTO converts domain.SiteSetting to responses.SiteSettingResponse
func (s *siteSettingService) toResponseDTO(setting *domain.SiteSetting) *responses.SiteSettingResponse {
	var faviconURL, logoHeaderURL, logoBigURL string

	if setting.Favicon != nil {
		faviconURL = s.cloudinaryService.GetImageURL("settings", *setting.Favicon)
	}
	if setting.LogoHeader != nil {
		logoHeaderURL = s.cloudinaryService.GetImageURL("settings", *setting.LogoHeader)
	}
	if setting.LogoBig != nil {
		logoBigURL = s.cloudinaryService.GetImageURL("settings", *setting.LogoBig)
	}

	return &responses.SiteSettingResponse{
		SiteName:        setting.SiteName,
		SiteTitle:       setting.SiteTitle,
		SiteDescription: setting.SiteDescription,
		Favicon:         faviconURL,
		LogoHeader:      logoHeaderURL,
		LogoBig:         logoBigURL,
		FacebookURL:     setting.FacebookURL,
		TwitterURL:      setting.TwitterURL,
		LinkedinURL:     setting.LinkedinURL,
		InstagramURL:    setting.InstagramURL,
		YoutubeURL:      setting.YoutubeURL,
		GithubURL:       setting.GithubURL,
		UpdatedAt:       setting.UpdatedAt,
	}
}
