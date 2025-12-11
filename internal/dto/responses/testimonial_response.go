package responses

import "time"

// TestimonialResponse adalah DTO untuk response testimonial
type TestimonialResponse struct {
	TestimonialID        int       `json:"testimonialId"`
	TestimonialName      string    `json:"testimonialName"`
	TestimonialOrg       *string   `json:"testimonialOrg"`
	TestimonialPosition  *string   `json:"testimonialPosition"`
	TestimonialContent   string    `json:"testimonialContent"`
	TestimonialImage     string    `json:"testimonialImage"`
	TestimonialIsActive  bool      `json:"testimonialIsActive,omitempty"`
	TestimonialCreatedAt time.Time `json:"testimonialCreatedAt"`
}
