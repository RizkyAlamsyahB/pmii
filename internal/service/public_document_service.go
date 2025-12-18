package service

import (
	"context"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// PublicDocumentService interface untuk public document API
type PublicDocumentService interface {
	GetAllPublic(ctx context.Context) ([]responses.PublicDocumentGroupResponse, error)
	GetByTypePublic(ctx context.Context, fileType string) (*responses.PublicDocumentGroupResponse, error)
}

type publicDocumentService struct {
	documentRepo      repository.DocumentRepository
	cloudinaryService CloudinaryService
}

// NewPublicDocumentService constructor untuk PublicDocumentService
func NewPublicDocumentService(
	documentRepo repository.DocumentRepository,
	cloudinaryService CloudinaryService,
) PublicDocumentService {
	return &publicDocumentService{
		documentRepo:      documentRepo,
		cloudinaryService: cloudinaryService,
	}
}

// GetAllPublic mengambil semua document aktif grouped by type (untuk public)
func (s *publicDocumentService) GetAllPublic(ctx context.Context) ([]responses.PublicDocumentGroupResponse, error) {
	result := make([]responses.PublicDocumentGroupResponse, 0, 3)

	// Loop through all document types
	for _, docType := range domain.ValidDocumentTypes() {
		documents, err := s.documentRepo.FindActiveByType(string(docType))
		if err != nil || len(documents) == 0 {
			continue // Skip if error or no documents
		}

		// Convert to public response
		publicDocs := make([]responses.PublicDocumentResponse, len(documents))
		for i, d := range documents {
			publicDocs[i] = s.toPublicResponseDTO(&d)
		}

		result = append(result, responses.PublicDocumentGroupResponse{
			FileType:      string(docType),
			FileTypeLabel: docType.GetLabel(),
			Documents:     publicDocs,
		})
	}

	return result, nil
}

// GetByTypePublic mengambil document aktif berdasarkan type (untuk public)
func (s *publicDocumentService) GetByTypePublic(ctx context.Context, fileType string) (*responses.PublicDocumentGroupResponse, error) {
	// Validate document type
	docType := domain.DocumentType(fileType)
	if !docType.IsValid() {
		return nil, nil // Return nil if invalid type
	}

	// Get documents by type
	documents, err := s.documentRepo.FindActiveByType(fileType)
	if err != nil {
		return nil, nil
	}

	if len(documents) == 0 {
		return nil, nil
	}

	// Convert to public response
	publicDocs := make([]responses.PublicDocumentResponse, len(documents))
	for i, d := range documents {
		publicDocs[i] = s.toPublicResponseDTO(&d)
	}

	result := &responses.PublicDocumentGroupResponse{
		FileType:      string(docType),
		FileTypeLabel: docType.GetLabel(),
		Documents:     publicDocs,
	}

	return result, nil
}

// toPublicResponseDTO converts domain.Document to responses.PublicDocumentResponse
// Uses GetDownloadURL to force download instead of preview
func (s *publicDocumentService) toPublicResponseDTO(d *domain.Document) responses.PublicDocumentResponse {
	folder := d.FileType.GetCloudinaryFolder()
	fileURL := s.cloudinaryService.GetDownloadURL(folder, d.FileURI)

	return responses.PublicDocumentResponse{
		ID:      d.ID,
		Name:    d.Name,
		FileURL: fileURL,
	}
}
