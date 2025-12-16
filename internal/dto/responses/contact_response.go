package responses

import "time"

// ContactResponse adalah DTO untuk response contact info
type ContactResponse struct {
	Address       *string   `json:"address,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Phone         *string   `json:"phone,omitempty"`
	GoogleMapsURL *string   `json:"googleMapsUrl,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PublicContactResponse adalah response untuk public API
type PublicContactResponse struct {
	Address       *string `json:"address,omitempty"`
	Email         *string `json:"email,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	GoogleMapsURL *string `json:"googleMapsUrl,omitempty"`
}
