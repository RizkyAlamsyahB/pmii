package responses

import (
	"strings"
	"time"

	// Ganti dengan nama module Anda
	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// PostResponse adalah bentuk JSON yang akan dikirim ke client
type PostResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Excerpt     string    `json:"excerpt"`
	Content     string    `json:"content,omitempty"`
	ImageUrl    string    `json:"imageUrl"`
	PublishedAt time.Time `json:"publishedAt"`
	Views       int       `json:"views"`

	CategoryId int    `json:"categoryId"`
	AuthorId   int    `json:"authorId"`
	Tags       string `json:"tags"`
}

func FromDomainToPostResponse(post domain.Post) PostResponse {
	// Convert Array Struct Tag -> Comma Separated String
	var tagNames []string
	for _, tag := range post.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	tagsString := strings.Join(tagNames, ",")

	return PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Excerpt:     post.Excerpt,
		Content:     post.Content,
		ImageUrl:    post.FeaturedImage,
		PublishedAt: post.PublishedAt,
		Views:       post.Views,
		CategoryId:  post.CategoryID,
		AuthorId:    post.UserID,
		Tags:        tagsString,
	}
}

// Helper untuk convert List (Array)
func FromDomainListToPostResponse(posts []domain.Post) []PostResponse {
	var responses []PostResponse
	for _, post := range posts {
		dto := FromDomainToPostResponse(post)
		dto.Content = ""
		responses = append(responses, dto)
	}
	return responses
}
