package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Service handles Cloudinary operations
type Service struct {
	cld *cloudinary.Cloudinary
}

// NewService creates a new Cloudinary service instance from URL
// url format: cloudinary://<api_key>:<api_secret>@<cloud_name>
func NewService(url string) (*Service, error) {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &Service{cld: cld}, nil
}

// UploadImage uploads image to Cloudinary and returns the filename only
// folder: target folder in Cloudinary (e.g., "testimonials", "profiles")
// file: multipart file from request
// Returns: filename only (e.g., "abc123.jpg"), error
func (s *Service) UploadImage(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename with timestamp
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

	// Upload to Cloudinary
	uploadResult, err := s.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder:   fmt.Sprintf("uploads/%s", folder),             // uploads/testimonials/
		PublicID: uniqueFilename[:len(uniqueFilename)-len(ext)], // filename without extension
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return filename only (e.g., "1234567890.jpg")
	filename := fmt.Sprintf("%s%s", uploadResult.PublicID[len(fmt.Sprintf("uploads/%s/", folder)):], ext)
	return filename, nil
}

// DeleteImage deletes image from Cloudinary
// folder: target folder in Cloudinary (e.g., "testimonials")
// filename: filename to delete (e.g., "abc123.jpg")
func (s *Service) DeleteImage(ctx context.Context, folder string, filename string) error {
	// Remove extension from filename
	ext := filepath.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(ext)]

	// Construct public ID
	publicID := fmt.Sprintf("uploads/%s/%s", folder, filenameWithoutExt)

	// Delete from Cloudinary
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete from Cloudinary: %w", err)
	}

	return nil
}

// GetImageURL returns full Cloudinary URL for a given filename
// folder: target folder (e.g., "testimonials")
// filename: filename (e.g., "abc123.jpg")
func (s *Service) GetImageURL(folder string, filename string) string {
	if filename == "" {
		return ""
	}

	ext := filepath.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(ext)]

	publicID := fmt.Sprintf("uploads/%s/%s", folder, filenameWithoutExt)
	asset, err := s.cld.Image(publicID)
	if err != nil {
		return ""
	}
	url, _ := asset.String()
	return url
}
