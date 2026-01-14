package service

import (
	"context"
	"math"
	"strings"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

type TagService interface {
	GetAll(page, limit int, search string) ([]responses.TagResponse, int, int64, error)
	Create(ctx context.Context, req requests.TagRequest) (responses.TagResponse, error)
	Update(ctx context.Context, id string, req requests.TagRequest) (responses.TagResponse, error)
	Delete(ctx context.Context, id string) error
}

type tagService struct {
	repo            repository.TagRepository
	activityLogRepo repository.ActivityLogRepository
}

func NewTagService(repo repository.TagRepository, activityLogRepo repository.ActivityLogRepository) TagService {
	return &tagService{
		repo:            repo,
		activityLogRepo: activityLogRepo,
	}
}

func (s *tagService) GetAll(page, limit int, search string) ([]responses.TagResponse, int, int64, error) {
	offset := (page - 1) * limit
	tags, total, err := s.repo.FindAll(offset, limit, search)
	if err != nil {
		return nil, 0, 0, err
	}

	// Hitung lastPage
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage == 0 {
		lastPage = 1
	}

	return responses.FromDomainListToTagResponse(tags), lastPage, total, nil
}

func (s *tagService) Create(ctx context.Context, req requests.TagRequest) (responses.TagResponse, error) {
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	tag := domain.Tag{
		Name: req.Name,
		Slug: slug,
	}

	if err := s.repo.Create(&tag); err != nil {
		return responses.TagResponse{}, err
	}

	// Log activity - Create Tag
	s.logActivity(ctx, domain.ActionCreate, domain.ModuleTags, "Membuat tag baru: "+tag.Name, nil, map[string]any{
		"id":   tag.ID,
		"name": tag.Name,
		"slug": tag.Slug,
	}, &tag.ID)

	return responses.FromDomainToTagResponse(tag), nil
}

func (s *tagService) Update(ctx context.Context, id string, req requests.TagRequest) (responses.TagResponse, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return responses.TagResponse{}, err
	}

	// Store old values for audit
	oldValues := map[string]any{
		"name": tag.Name,
		"slug": tag.Slug,
	}

	tag.Name = req.Name
	tag.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	if err := s.repo.Update(&tag); err != nil {
		return responses.TagResponse{}, err
	}

	// Log activity - Update Tag
	s.logActivity(ctx, domain.ActionUpdate, domain.ModuleTags, "Mengupdate tag: "+tag.Name, oldValues, map[string]any{
		"id":   tag.ID,
		"name": tag.Name,
		"slug": tag.Slug,
	}, &tag.ID)

	return responses.FromDomainToTagResponse(tag), nil
}

func (s *tagService) Delete(ctx context.Context, id string) error {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Log activity sebelum delete
	s.logActivity(ctx, domain.ActionDelete, domain.ModuleTags, "Menghapus tag: "+tag.Name, map[string]any{
		"id":   tag.ID,
		"name": tag.Name,
		"slug": tag.Slug,
	}, nil, &tag.ID)

	return s.repo.Delete(&tag)
}

// logActivity helper untuk mencatat activity log
func (s *tagService) logActivity(ctx context.Context, actionType domain.ActivityActionType, module domain.ActivityModuleType, description string, oldValue, newValue map[string]any, targetID *int) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
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
	_ = s.activityLogRepo.Create(log)
}
