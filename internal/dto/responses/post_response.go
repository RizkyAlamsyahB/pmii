package responses

import (
	"strings"
	"time"

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
	// Jika ingin menampilkan data user dan kategori lengkap, bisa
	// mengganti CategoryId dan AuthorId dengan struct CategoryResponse dan UserResponse
}

func FromDomainToPostResponse(post domain.Post) PostResponse {
	// 1. Convert Array Struct Tag -> Comma Separated String
	var tagNames []string
	for _, tag := range post.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	tagsString := strings.Join(tagNames, ",")

	// 2. Handling Pointers (untuk Excerpt, FeaturedImage, PublishedAt)

	// Excerpt
	excerpt := ""
	if post.Excerpt != nil {
		excerpt = *post.Excerpt
	}

	// ImageUrl (FeaturedImage)
	imageUrl := ""
	if post.FeaturedImage != nil {
		imageUrl = *post.FeaturedImage
	}

	// PublishedAt
	var publishedAt time.Time
	if post.PublishedAt != nil {
		publishedAt = *post.PublishedAt
	} else {
		// Jika PublishedAt null/nil (misalnya status masih 'draft'), gunakan waktu sekarang atau CreatedAt
		// Saya gunakan CreatedAt karena itu pasti terisi
		publishedAt = post.CreatedAt
	}

	// 3. Return Response
	return PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Excerpt:     excerpt,
		Content:     post.Content,
		ImageUrl:    imageUrl,
		PublishedAt: publishedAt,
		// Views:       post.Views,
		CategoryId: post.CategoryID,
		AuthorId:   post.UserID,
		Tags:       tagsString,
	}
}

// Helper untuk convert List (Array)
func FromDomainListToPostResponse(posts []domain.Post) []PostResponse {
	var responses []PostResponse
	for _, post := range posts {
		dto := FromDomainToPostResponse(post)
		// Menghapus Content dari list response agar payload lebih kecil (jika itu intent Anda)
		dto.Content = ""
		responses = append(responses, dto)
	}
	return responses
}
