package requests

// UpdateAboutRequest adalah DTO untuk update about page
type UpdateAboutRequest struct {
	Title    string `form:"title"`
	History  string `form:"history"`
	Vision   string `form:"vision"`
	Mission  string `form:"mission"`
	VideoURL string `form:"video_url"`
}
