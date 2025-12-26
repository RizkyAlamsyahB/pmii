package requests

import "strings"

type CategoryRequest struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
}

func (r *CategoryRequest) GetSlug() string {
	return strings.ToLower(strings.ReplaceAll(r.Name, " ", "-"))
}

type TagRequest struct {
	Name string `form:"name" binding:"required"`
}

func (r *TagRequest) GetSlug() string {
	return strings.ToLower(strings.ReplaceAll(r.Name, " ", "-"))
}
