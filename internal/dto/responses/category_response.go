package responses

import "github.com/garuda-labs-1/pmii-be/internal/domain"

// CategoryResponse adalah struct JSON output untuk category
type CategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"` // Opsional
}

// Mapper: Domain -> Response
func FromDomainToCategoryResponse(category domain.Category) CategoryResponse {

	desc := ""

	if category.Description != nil {
		desc = *category.Description
	}

	return CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: desc,
	}
}

// Mapper: List Domain -> List Response
func FromDomainListToCategoryResponse(categories []domain.Category) []CategoryResponse {
	var responses []CategoryResponse
	for _, cat := range categories {
		responses = append(responses, FromDomainToCategoryResponse(cat))
	}
	return responses
}
