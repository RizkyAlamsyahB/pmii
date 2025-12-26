package requests

import (
	"mime/multipart"
	"strings"
)

type PostCreateRequest struct {
	Title      string                `form:"title" binding:"required"`
	Content    string                `form:"content" binding:"required"`
	CategoryID int                   `form:"category_id" binding:"required"`
	Tags       string                `form:"tags"`
	Image      *multipart.FileHeader `form:"image"`
}

type PostUpdateRequest struct {
	Title      string                `form:"title"`
	Content    string                `form:"content"`
	CategoryID int                   `form:"category_id"`
	Tags       string                `form:"tags"`
	Image      *multipart.FileHeader `form:"image"`
}

func (r *PostCreateRequest) GetSlug() string {
	return strings.ToLower(strings.ReplaceAll(r.Title, " ", "-"))
}
