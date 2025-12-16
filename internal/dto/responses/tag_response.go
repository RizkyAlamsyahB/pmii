package responses

import "github.com/garuda-labs-1/pmii-be/internal/domain"

type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func FromDomainToTagResponse(tag domain.Tag) TagResponse {
	return TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}
}

func FromDomainListToTagResponse(tags []domain.Tag) []TagResponse {
	var responses []TagResponse
	for _, tag := range tags {
		responses = append(responses, FromDomainToTagResponse(tag))
	}
	return responses
}
