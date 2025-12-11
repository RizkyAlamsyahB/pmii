package requests

// CreateTestimonialRequest adalah DTO untuk membuat testimonial baru
type CreateTestimonialRequest struct {
	Name         string `form:"name" binding:"required"`
	Organization string `form:"organization"`
	Position     string `form:"position"`
	Content      string `form:"content" binding:"required"`
	// Photo akan dihandle terpisah menggunakan c.FormFile("photo")
}

// UpdateTestimonialRequest adalah DTO untuk update testimonial
type UpdateTestimonialRequest struct {
	Name         string `form:"name"`
	Organization string `form:"organization"`
	Position     string `form:"position"`
	Content      string `form:"content"`
	IsActive     *bool  `form:"is_active"`
	// Photo akan dihandle terpisah menggunakan c.FormFile("photo")
}
