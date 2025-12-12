package requests

// CreateMemberRequest adalah DTO untuk membuat member baru
type CreateMemberRequest struct {
	FullName    string         `form:"full_name" binding:"required"`
	Position    string         `form:"position" binding:"required"`
	SocialLinks map[string]any `form:"social_links"` // Will handle as JSON
	// Photo akan dihandle terpisah menggunakan c.FormFile("photo")
}

// UpdateMemberRequest adalah DTO untuk update member
type UpdateMemberRequest struct {
	FullName    string         `form:"full_name"`
	Position    string         `form:"position"`
	SocialLinks map[string]any `form:"social_links"` // Will handle as JSON
	IsActive    *bool          `form:"is_active"`
	// Photo akan dihandle terpisah menggunakan c.FormFile("photo")
}
