package requests

import "mime/multipart"

// UpdateAdRequest untuk update single ad slot (image only)
type UpdateAdRequest struct {
	Image *multipart.FileHeader `form:"image"`
}

// UpdateAdImageRequest untuk update image via form-data
type UpdateAdImageRequest struct {
	Image *multipart.FileHeader `form:"image" binding:"required"`
}
