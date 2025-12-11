package responses

import "time"

// TestimonialResponse adalah DTO untuk response testimonial (mengikuti standar API)
type TestimonialResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Organization *string   `json:"organization,omitempty"`
	Position     *string   `json:"position,omitempty"`
	Content      string    `json:"content"`
	ImageUrl     string    `json:"imageUrl,omitempty"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
}
