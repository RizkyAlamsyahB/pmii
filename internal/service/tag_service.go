package service

import (
	"math"
	"strings"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

type TagService interface {
	GetAll(page, limit int, search string) ([]responses.TagResponse, int, int64, error)
	Create(req requests.TagRequest) (responses.TagResponse, error)
	Update(id string, req requests.TagRequest) (responses.TagResponse, error)
	Delete(id string) error
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
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

func (s *tagService) Create(req requests.TagRequest) (responses.TagResponse, error) {
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	tag := domain.Tag{
		Name: req.Name,
		Slug: slug,
	}

	if err := s.repo.Create(&tag); err != nil {
		return responses.TagResponse{}, err
	}
	return responses.FromDomainToTagResponse(tag), nil
}

func (s *tagService) Update(id string, req requests.TagRequest) (responses.TagResponse, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return responses.TagResponse{}, err
	}

	tag.Name = req.Name
	tag.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	if err := s.repo.Update(&tag); err != nil {
		return responses.TagResponse{}, err
	}
	return responses.FromDomainToTagResponse(tag), nil
}

func (s *tagService) Delete(id string) error {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(&tag)
}
