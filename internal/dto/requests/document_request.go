package requests

// CreateDocumentRequest adalah DTO untuk create document
type CreateDocumentRequest struct {
	Name     string `form:"name" binding:"required"`
	FileType string `form:"file_type" binding:"required"`
	// File handled separately via multipart
}

// UpdateDocumentRequest adalah DTO untuk update document
type UpdateDocumentRequest struct {
	Name     string `form:"name"`
	FileType string `form:"file_type"`
	// File handled separately via multipart (optional on update)
}
