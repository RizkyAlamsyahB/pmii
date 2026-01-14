package responses

// TestimonialResponse adalah DTO untuk response testimonial list (admin table)
type TestimonialResponse struct {
	ID       int     `json:"id"`
	ImageUrl string  `json:"imageUrl,omitempty"`
	Name     string  `json:"name"`
	Position *string `json:"position,omitempty"`
	Content  string  `json:"content"`
}

// TestimonialDetailResponse adalah DTO untuk response testimonial detail (admin edit form)
type TestimonialDetailResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Organization *string `json:"organization,omitempty"`
	Position     *string `json:"position,omitempty"`
	Content      string  `json:"content"`
	ImageUrl     string  `json:"imageUrl,omitempty"`
	IsActive     bool    `json:"isActive"`
}
