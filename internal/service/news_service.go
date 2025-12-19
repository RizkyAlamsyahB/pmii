package service

import (
	"math"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

type NewsService interface {
	FetchPublicNews(page, limit int, search string) ([]responses.PostResponse, int, int64, error)
	FetchNewsDetail(slug string) (responses.PostResponse, error)
	//metod untuk mendapatkan berita berdasarkan kategori
	FetchNewsByCategory(categorySlug string, page, limit int) ([]responses.PostResponse, int, int64, error)
}

type newsService struct {
	repo repository.NewsRepository
}

func NewNewsService(repo repository.NewsRepository) NewsService {
	return &newsService{repo: repo.(repository.NewsRepository)}
}

func (s *newsService) FetchPublicNews(page, limit int, search string) ([]responses.PostResponse, int, int64, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.GetPublishedNews(offset, limit, search)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	data := responses.FromDomainListToPostResponse(posts)

	return data, lastPage, total, nil
}

func (s *newsService) FetchNewsDetail(slug string) (responses.PostResponse, error) {
	post, err := s.repo.GetNewsBySlug(slug)
	if err != nil {
		return responses.PostResponse{}, err
	}
	return responses.FromDomainToPostResponse(post), nil
}

// metod untuk mendapatkan berita berdasarkan kategori
func (s *newsService) FetchNewsByCategory(categorySlug string, page, limit int) ([]responses.PostResponse, int, int64, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.GetNewsByCategorySlug(categorySlug, offset, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	data := responses.FromDomainListToPostResponse(posts)

	return data, lastPage, total, nil
}
