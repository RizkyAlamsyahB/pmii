package responses

import "time"

// DocumentResponse adalah DTO untuk response document (admin)
type DocumentResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	FileType      string    `json:"fileType"`
	FileTypeLabel string    `json:"fileTypeLabel"`
	FileURL       string    `json:"fileUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PublicDocumentResponse adalah response untuk public API
type PublicDocumentResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	FileURL string `json:"fileUrl"`
}

// DocumentTypeInfo informasi tipe dokumen
type DocumentTypeInfo struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// PublicDocumentGroupResponse adalah response untuk public API dengan grouping
type PublicDocumentGroupResponse struct {
	FileType      string                   `json:"fileType"`
	FileTypeLabel string                   `json:"fileTypeLabel"`
	Documents     []PublicDocumentResponse `json:"documents"`
}
