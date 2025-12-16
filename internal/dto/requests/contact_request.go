package requests

// UpdateContactRequest adalah DTO untuk update contact info
type UpdateContactRequest struct {
	Address       *string `form:"address" json:"address"`
	Email         *string `form:"email" json:"email"`
	Phone         *string `form:"phone" json:"phone"`
	GoogleMapsURL *string `form:"google_maps_url" json:"google_maps_url"`
}
