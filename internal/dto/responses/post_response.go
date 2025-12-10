package responses

import (
	"time"

	// Ganti dengan nama module Anda
	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// PostResponse adalah bentuk JSON yang akan dikirim ke client
type PostResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Excerpt     string    `json:"excerpt"`           // Mapping dari Description
	Content     string    `json:"content,omitempty"` // Omitempty agar di list tidak berat
	ImageUrl    string    `json:"imageUrl"`          // Sesuai request: imageUrl
	AuthorId    int       `json:"authorId"`
	CategoryId  int       `json:"categoryId"`
	Tags        string    `json:"tags"`
	PublishedAt time.Time `json:"publishedAt"`
	Views       int       `json:"views"`
}

// Function Helper untuk convert dari Domain (Database) ke Response (JSON)
func FromDomainToPostResponse(post domain.Post) PostResponse {
	return PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Excerpt:     post.Description, // Map description DB ke excerpt JSON
		Content:     post.Content,
		ImageUrl:    post.Image,
		AuthorId:    post.UserID,
		CategoryId:  post.CategoryID,
		Tags:        post.Tags,
		PublishedAt: post.Date,
		Views:       post.Views,
	}
}

// Helper untuk convert List (Array)
func FromDomainListToPostResponse(posts []domain.Post) []PostResponse {
	var responses []PostResponse
	for _, post := range posts {
		// Untuk list, kita kosongkan content agar payload ringan
		dto := FromDomainToPostResponse(post)
		dto.Content = ""
		responses = append(responses, dto)
	}
	return responses
}
