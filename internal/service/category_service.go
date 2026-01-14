package service

import (
	"context"
	"math"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

type CategoryService interface {
	GetAll(page, limit int, search string) ([]responses.CategoryResponse, int, int64, error)
	Create(ctx context.Context, req requests.CategoryRequest) (responses.CategoryResponse, error)
	Update(ctx context.Context, id string, req requests.CategoryRequest) (responses.CategoryResponse, error)
	Delete(ctx context.Context, id string) error
}

type categoryService struct {
	repo            repository.CategoryRepository
	activityLogRepo repository.ActivityLogRepository
}

func NewCategoryService(repo repository.CategoryRepository, activityLogRepo repository.ActivityLogRepository) CategoryService {
	return &categoryService{
		repo:            repo,
		activityLogRepo: activityLogRepo,
	}
}

func (s *categoryService) GetAll(page, limit int, search string) ([]responses.CategoryResponse, int, int64, error) {
	offset := (page - 1) * limit
	categories, total, err := s.repo.FindAll(offset, limit, search)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage == 0 {
		lastPage = 1
	}

	return responses.FromDomainListToCategoryResponse(categories), lastPage, total, nil
}

func (s *categoryService) Create(ctx context.Context, req requests.CategoryRequest) (responses.CategoryResponse, error) {
	var descPtr *string
	if req.Description != "" {
		descPtr = &req.Description
	}

	category := domain.Category{
		Name:        req.Name,
		Slug:        req.GetSlug(),
		Description: descPtr,
	}

	if err := s.repo.Create(&category); err != nil {
		return responses.CategoryResponse{}, err
	}

	// Log activity - Create Category
	s.logActivity(ctx, domain.ActionCreate, domain.ModuleCategory, "Membuat kategori baru: "+category.Name, nil, map[string]any{
		"id":   category.ID,
		"name": category.Name,
		"slug": category.Slug,
	}, &category.ID)

	return responses.FromDomainToCategoryResponse(category), nil
}

func (s *categoryService) Update(ctx context.Context, id string, req requests.CategoryRequest) (responses.CategoryResponse, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return responses.CategoryResponse{}, err
	}

	// Store old values for audit
	oldValues := map[string]any{
		"name": category.Name,
		"slug": category.Slug,
	}

	category.Name = req.Name
	category.Slug = req.GetSlug()
	if req.Description != "" {
		category.Description = &req.Description
	}

	if err := s.repo.Update(&category); err != nil {
		return responses.CategoryResponse{}, err
	}

	// Log activity - Update Category
	s.logActivity(ctx, domain.ActionUpdate, domain.ModuleCategory, "Mengupdate kategori: "+category.Name, oldValues, map[string]any{
		"id":   category.ID,
		"name": category.Name,
		"slug": category.Slug,
	}, &category.ID)

	return responses.FromDomainToCategoryResponse(category), nil
}

func (s *categoryService) Delete(ctx context.Context, id string) error {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Log activity sebelum delete
	s.logActivity(ctx, domain.ActionDelete, domain.ModuleCategory, "Menghapus kategori: "+category.Name, map[string]any{
		"id":   category.ID,
		"name": category.Name,
		"slug": category.Slug,
	}, nil, &category.ID)

	return s.repo.Delete(&category)
}

// logActivity helper untuk mencatat activity log
func (s *categoryService) logActivity(ctx context.Context, actionType domain.ActivityActionType, module domain.ActivityModuleType, description string, oldValue, newValue map[string]any, targetID *int) {
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
