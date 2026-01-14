package responses

// AboutResponse adalah DTO untuk response about page
type AboutResponse struct {
	ID       int     `json:"id"`
	History  *string `json:"history,omitempty"`
	Vision   *string `json:"vision,omitempty"`
	Mission  *string `json:"mission,omitempty"`
	ImageUrl string  `json:"imageUrl,omitempty"`
	VideoURL *string `json:"videoUrl,omitempty"`
}
