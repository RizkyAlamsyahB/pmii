package responses

// DocumentResponse adalah DTO untuk response document list (admin)
type DocumentResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FileTypeLabel string `json:"fileTypeLabel"`
	FileURL       string `json:"fileUrl"`
}

// DocumentDetailResponse adalah DTO untuk response document detail (admin - get by ID)
type DocumentDetailResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FileType      string `json:"fileType"`
	FileTypeLabel string `json:"fileTypeLabel"`
	FileURL       string `json:"fileUrl"`
}

// PublicDocumentResponse adalah response untuk public API
type PublicDocumentResponse struct {
	Name          string `json:"name"`
	FileTypeLabel string `json:"fileTypeLabel"`
	FileURL       string `json:"fileUrl"`
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
