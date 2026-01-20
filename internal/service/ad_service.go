package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// AdService interface untuk business logic ads
type AdService interface {
	GetAllAds(ctx context.Context) ([]responses.AdResponse, error)
	GetAdByID(ctx context.Context, id int) (*responses.AdResponse, error)
	GetAdsByPage(ctx context.Context, page string) ([]responses.AdResponse, error)
	UpdateAd(ctx context.Context, id int, image *multipart.FileHeader) (*responses.AdResponse, error)
	DeleteAdImage(ctx context.Context, id int) (*responses.AdResponse, error)
}

type adService struct {
	adRepo          repository.AdRepository
	cloudinary      CloudinaryService
	activityLogRepo repository.ActivityLogRepository
}

// NewAdService constructor untuk AdService
func NewAdService(adRepo repository.AdRepository, cloudinaryService CloudinaryService, activityLogRepo repository.ActivityLogRepository) AdService {
	return &adService{
		adRepo:          adRepo,
		cloudinary:      cloudinaryService,
		activityLogRepo: activityLogRepo,
	}
}

// GetAllAds mengambil semua ads sebagai flat array
func (s *adService) GetAllAds(ctx context.Context) ([]responses.AdResponse, error) {
	ads, err := s.adRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Convert to response with full image URLs
	return s.toAdResponseListWithFullURL(ads), nil
}

// GetAdByID mengambil ad berdasarkan ID
func (s *adService) GetAdByID(ctx context.Context, id int) (*responses.AdResponse, error) {
	ad, err := s.adRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ad tidak ditemukan")
	}

	response := s.toAdResponseWithFullURL(ad)
	return &response, nil
}

// GetAdsByPage mengambil ads berdasarkan page
func (s *adService) GetAdsByPage(ctx context.Context, page string) ([]responses.AdResponse, error) {
	if !domain.IsValidAdPage(page) {
		return nil, errors.New("halaman tidak valid")
	}

	ads, err := s.adRepo.FindByPage(domain.AdPage(page))
	if err != nil {
		return nil, err
	}

	return s.toAdResponseListWithFullURL(ads), nil
}

// UpdateAd mengupdate ad image
func (s *adService) UpdateAd(ctx context.Context, id int, image *multipart.FileHeader) (*responses.AdResponse, error) {
	// Get existing ad
	ad, err := s.adRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ad tidak ditemukan")
	}

	oldData := map[string]any{
		"image_url": ad.ImageURL,
	}

	// Upload new image if provided
	if image != nil {
		// Validate image aspect ratio (10% tolerance)
		if err := utils.ValidateImageAspectRatio(image, ad.Resolution, 0.1); err != nil {
			return nil, err
		}

		// Delete old image from cloudinary if exists
		if ad.ImageURL != nil && *ad.ImageURL != "" {
			_ = s.cloudinary.DeleteImage(ctx, "ads", *ad.ImageURL)
		}

		// Upload new image
		imageFilename, err := s.cloudinary.UploadImage(ctx, "ads", image)
		if err != nil {
			return nil, errors.New("gagal mengupload gambar: " + err.Error())
		}
		ad.ImageURL = &imageFilename
	}

	ad.UpdatedAt = time.Now()

	// Save to database
	if err := s.adRepo.Update(ad); err != nil {
		return nil, errors.New("gagal mengupdate ad")
	}

	// Log activity
	newData := map[string]any{
		"image_url": ad.ImageURL,
	}
	s.logActivity(ctx, domain.ActionUpdate, domain.ModuleAds, "Update "+ad.GetSlotName(), oldData, newData, &ad.ID)

	response := s.toAdResponseWithFullURL(ad)
	return &response, nil
}

// DeleteAdImage menghapus image ad (set to null)
func (s *adService) DeleteAdImage(ctx context.Context, id int) (*responses.AdResponse, error) {
	// Get existing ad
	ad, err := s.adRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ad tidak ditemukan")
	}

	oldImageURL := ad.ImageURL

	// Delete image from cloudinary if exists
	if ad.ImageURL != nil && *ad.ImageURL != "" {
		_ = s.cloudinary.DeleteImage(ctx, "ads", *ad.ImageURL)
	}

	// Set image to nil
	ad.ImageURL = nil
	ad.UpdatedAt = time.Now()

	// Save to database
	if err := s.adRepo.Update(ad); err != nil {
		return nil, errors.New("gagal menghapus gambar ad")
	}

	// Log activity
	s.logActivity(ctx, domain.ActionDelete, domain.ModuleAds, "Delete image "+ad.GetSlotName(), map[string]any{
		"deleted_image_url": oldImageURL,
	}, nil, &ad.ID)

	response := s.toAdResponseWithFullURL(ad)
	return &response, nil
}

// logActivity helper untuk mencatat activity log
func (s *adService) logActivity(ctx context.Context, actionType domain.ActivityActionType, module domain.ActivityModuleType, description string, oldValue, newValue map[string]any, targetID *int) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		// Log warning for debugging
		fmt.Println("[WARN] logActivity: user_id not found in context, skipping activity log")
		return // Skip if no user in context
	}

	ipAddress := utils.GetIPAddress(ctx)
	userAgent := utils.GetUserAgent(ctx)

	var ipPtr, uaPtr *string
	if ipAddress != "" {
		ipPtr = &ipAddress
	}
	if userAgent != "" {
		uaPtr = &userAgent
	}

	log := &domain.ActivityLog{
		UserID:      userID,
		ActionType:  actionType,
		Module:      module,
		Description: &description,
		TargetID:    targetID,
		OldValue:    oldValue,
		NewValue:    newValue,
		IPAddress:   ipPtr,
		UserAgent:   uaPtr,
	}

	// Ignore error - logging should not affect main operation
	if err := s.activityLogRepo.Create(log); err != nil {
		fmt.Printf("[WARN] logActivity: failed to create activity log: %v\n", err)
	}
}

// toAdResponseWithFullURL converts domain.Ad to AdResponse with full cloudinary URL
func (s *adService) toAdResponseWithFullURL(ad *domain.Ad) responses.AdResponse {
	var imageURL *string
	if ad.ImageURL != nil && *ad.ImageURL != "" {
		fullURL := s.cloudinary.GetImageURL("ads", *ad.ImageURL)
		imageURL = &fullURL
	}

	return responses.AdResponse{
		SlotName: ad.GetSlotName(),
		ImageURL: imageURL,
	}
}

// toAdResponseListWithFullURL converts slice of domain.Ad to slice of AdResponse with full URLs
func (s *adService) toAdResponseListWithFullURL(ads []domain.Ad) []responses.AdResponse {
	result := make([]responses.AdResponse, len(ads))
	for i, ad := range ads {
		result[i] = s.toAdResponseWithFullURL(&ad)
	}
	return result
}
