package requests

// UpdateAboutRequest adalah DTO untuk update about page
type UpdateAboutRequest struct {
	History  string `form:"history"`
	Vision   string `form:"vision"`
	Mission  string `form:"mission"`
	VideoURL string `form:"video_url"`
	// Image akan dihandle terpisah menggunakan c.FormFile("image")
}
