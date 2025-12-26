package service

import (
	"math"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

type CategoryService interface {
	GetAll(page, limit int, search string) ([]responses.CategoryResponse, int, int64, error)
	Create(req requests.CategoryRequest) (responses.CategoryResponse, error)
	Update(id string, req requests.CategoryRequest) (responses.CategoryResponse, error)
	Delete(id string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
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

func (s *categoryService) Create(req requests.CategoryRequest) (responses.CategoryResponse, error) {
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
	return responses.FromDomainToCategoryResponse(category), nil
}

func (s *categoryService) Update(id string, req requests.CategoryRequest) (responses.CategoryResponse, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return responses.CategoryResponse{}, err
	}

	category.Name = req.Name
	category.Slug = req.GetSlug()
	if req.Description != "" {
		category.Description = &req.Description
	}

	if err := s.repo.Update(&category); err != nil {
		return responses.CategoryResponse{}, err
	}
	return responses.FromDomainToCategoryResponse(category), nil
}

func (s *categoryService) Delete(id string) error {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(&category)
}
