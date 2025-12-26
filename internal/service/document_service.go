package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// DocumentService interface untuk business logic document (admin)
type DocumentService interface {
	Create(ctx context.Context, req requests.CreateDocumentRequest, file *multipart.FileHeader) (*responses.DocumentResponse, error)
	GetAll(ctx context.Context, page, limit int, fileType, search string) ([]responses.DocumentResponse, int, int, int64, error)
	GetByID(ctx context.Context, id int) (*responses.DocumentResponse, error)
	Update(ctx context.Context, id int, req requests.UpdateDocumentRequest, file *multipart.FileHeader) (*responses.DocumentResponse, error)
	Delete(ctx context.Context, id int) error
	GetDocumentTypes() []responses.DocumentTypeInfo
}

type documentService struct {
	documentRepo      repository.DocumentRepository
	cloudinaryService CloudinaryService
}

// NewDocumentService constructor untuk DocumentService
func NewDocumentService(documentRepo repository.DocumentRepository, cloudinaryService CloudinaryService) DocumentService {
	return &documentService{
		documentRepo:      documentRepo,
		cloudinaryService: cloudinaryService,
	}
}

// Create membuat document baru dengan upload file ke Cloudinary
func (s *documentService) Create(ctx context.Context, req requests.CreateDocumentRequest, file *multipart.FileHeader) (*responses.DocumentResponse, error) {
	// Validate file type
	docType := domain.DocumentType(req.FileType)
	if !docType.IsValid() {
		return nil, errors.New("jenis file tidak valid. Pilih: produk_hukum, lagu_organisasi, atau logo_organisasi")
	}

	// File is required for create
	if file == nil {
		return nil, errors.New("file wajib diupload")
	}

	// Get Cloudinary folder based on document type
	folder := docType.GetCloudinaryFolder()

	// Upload file ke Cloudinary (using UploadFile for correct resource type)
	filename, err := s.cloudinaryService.UploadFile(ctx, folder, file)
	if err != nil {
		return nil, errors.New("gagal mengupload file")
	}

	// Prepare domain entity
	document := &domain.Document{
		Name:     req.Name,
		FileType: docType,
		FileURI:  filename,
	}

	// Save ke database
	if err := s.documentRepo.Create(document); err != nil {
		// Rollback: hapus file dari Cloudinary jika save gagal
		_ = s.cloudinaryService.DeleteFile(ctx, folder, filename)
		return nil, errors.New("gagal menyimpan dokumen")
	}

	return s.toResponseDTO(document), nil
}

// GetAll mengambil semua document dengan pagination dan search
func (s *documentService) GetAll(ctx context.Context, page, limit int, fileType, search string) ([]responses.DocumentResponse, int, int, int64, error) {
	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	documents, total, err := s.documentRepo.FindAll(page, limit, fileType, search)
	if err != nil {
		return nil, 0, 0, 0, errors.New("gagal mengambil data dokumen")
	}

	// Calculate last page
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}
	if lastPage < 1 && total > 0 {
		lastPage = 1
	}

	// Auto-clamp: jika page melebihi lastPage dan ada data, clamp ke lastPage
	if page > lastPage && lastPage > 0 {
		page = lastPage
		// Re-fetch dengan page yang sudah di-clamp
		documents, _, err = s.documentRepo.FindAll(page, limit, fileType, search)
		if err != nil {
			return nil, 0, 0, 0, errors.New("gagal mengambil data dokumen")
		}
	}

	// Convert to response DTOs
	result := make([]responses.DocumentResponse, len(documents))
	for i, d := range documents {
		result[i] = *s.toResponseDTO(&d)
	}

	return result, page, lastPage, total, nil
}

// GetByID mengambil document berdasarkan ID
func (s *documentService) GetByID(ctx context.Context, id int) (*responses.DocumentResponse, error) {
	document, err := s.documentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("dokumen tidak ditemukan")
	}

	return s.toResponseDTO(document), nil
}

// Update mengupdate document dengan optional upload file baru
func (s *documentService) Update(ctx context.Context, id int, req requests.UpdateDocumentRequest, file *multipart.FileHeader) (*responses.DocumentResponse, error) {
	// Ambil document existing
	document, err := s.documentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("dokumen tidak ditemukan")
	}

	// Simpan info lama untuk rollback/cleanup
	oldFileURI := document.FileURI
	oldFileType := document.FileType

	// Update file type if provided
	if req.FileType != "" {
		docType := domain.DocumentType(req.FileType)
		if !docType.IsValid() {
			return nil, errors.New("jenis file tidak valid. Pilih: produk_hukum, lagu_organisasi, atau logo_organisasi")
		}
		document.FileType = docType
	}

	// Upload file baru ke Cloudinary (jika ada)
	var newFilename *string
	if file != nil {
		folder := document.FileType.GetCloudinaryFolder()
		filename, err := s.cloudinaryService.UploadFile(ctx, folder, file)
		if err != nil {
			return nil, errors.New("gagal mengupload file")
		}
		newFilename = &filename
		document.FileURI = filename
	}

	// Update name if provided
	if req.Name != "" {
		document.Name = req.Name
	}

	// Save ke database
	if err := s.documentRepo.Update(document); err != nil {
		// Rollback: hapus file baru jika update gagal
		if newFilename != nil {
			_ = s.cloudinaryService.DeleteFile(ctx, document.FileType.GetCloudinaryFolder(), *newFilename)
		}
		return nil, errors.New("gagal mengupdate dokumen")
	}

	// Hapus file lama SETELAH database update berhasil
	if newFilename != nil {
		_ = s.cloudinaryService.DeleteFile(ctx, oldFileType.GetCloudinaryFolder(), oldFileURI)
	}

	return s.toResponseDTO(document), nil
}

// Delete menghapus document (soft delete) dan file dari Cloudinary
func (s *documentService) Delete(ctx context.Context, id int) error {
	// Ambil document untuk mendapatkan info file
	document, err := s.documentRepo.FindByID(id)
	if err != nil {
		return errors.New("dokumen tidak ditemukan")
	}

	// Soft delete dari database
	if err := s.documentRepo.Delete(id); err != nil {
		return errors.New("gagal menghapus dokumen")
	}

	// Hapus file dari Cloudinary SETELAH database delete berhasil
	folder := document.FileType.GetCloudinaryFolder()
	_ = s.cloudinaryService.DeleteFile(ctx, folder, document.FileURI)

	return nil
}

// GetDocumentTypes returns all available document types
func (s *documentService) GetDocumentTypes() []responses.DocumentTypeInfo {
	types := domain.ValidDocumentTypes()
	result := make([]responses.DocumentTypeInfo, len(types))
	for i, t := range types {
		result[i] = responses.DocumentTypeInfo{
			Value: string(t),
			Label: t.GetLabel(),
		}
	}
	return result
}

// toResponseDTO converts domain.Document to responses.DocumentResponse
// Uses GetDownloadURL for direct download
func (s *documentService) toResponseDTO(d *domain.Document) *responses.DocumentResponse {
	folder := d.FileType.GetCloudinaryFolder()
	fileURL := s.cloudinaryService.GetDownloadURL(folder, d.FileURI)

	return &responses.DocumentResponse{
		ID:            d.ID,
		Name:          d.Name,
		FileType:      string(d.FileType),
		FileTypeLabel: d.FileType.GetLabel(),
		FileURL:       fileURL,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
