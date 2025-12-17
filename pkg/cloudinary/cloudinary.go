package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// ResourceType represents Cloudinary resource types
type ResourceType string

const (
	ResourceTypeImage ResourceType = "image"
	ResourceTypeVideo ResourceType = "video" // Includes audio files (MP3, WAV, OGG)
	ResourceTypeRaw   ResourceType = "raw"   // PDF, DOC, ZIP, etc.
)

// GetResourceTypeFromExt determines the correct Cloudinary resource type from file extension
func GetResourceTypeFromExt(filename string) ResourceType {
	ext := strings.ToLower(filepath.Ext(filename))

	// Image extensions
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".bmp": true, ".ico": true, ".svg": true,
	}

	// Video & Audio extensions (Cloudinary uses "video" for both)
	videoExts := map[string]bool{
		".mp4": true, ".webm": true, ".mov": true, ".avi": true, ".mkv": true,
		".mp3": true, ".wav": true, ".ogg": true, ".flac": true, ".aac": true,
	}

	if imageExts[ext] {
		return ResourceTypeImage
	}
	if videoExts[ext] {
		return ResourceTypeVideo
	}
	return ResourceTypeRaw
}

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

	// Upload to Cloudinary (always as image for this method)
	uploadResult, err := s.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder:       fmt.Sprintf("uploads/%s", folder),
		PublicID:     uniqueFilename[:len(uniqueFilename)-len(ext)],
		ResourceType: "image",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return filename only (e.g., "1234567890.jpg")
	filename := fmt.Sprintf("%s%s", uploadResult.PublicID[len(fmt.Sprintf("uploads/%s/", folder)):], ext)
	return filename, nil
}

// UploadFile uploads any file to Cloudinary with correct resource type
// folder: target folder in Cloudinary (e.g., "documents/produk_hukum")
// file: multipart file from request
// Returns: filename with extension (e.g., "abc123.pdf"), error
func (s *Service) UploadFile(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename with timestamp
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	publicID := fmt.Sprintf("%d", timestamp)

	// Determine resource type from file extension
	resourceType := GetResourceTypeFromExt(file.Filename)

	// Upload to Cloudinary with correct resource type
	uploadResult, err := s.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder:       fmt.Sprintf("uploads/%s", folder),
		PublicID:     publicID,
		ResourceType: string(resourceType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return filename with extension (e.g., "1234567890.pdf")
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

	// Delete from Cloudinary (image resource type)
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})
	if err != nil {
		return fmt.Errorf("failed to delete from Cloudinary: %w", err)
	}

	return nil
}

// DeleteFile deletes file from Cloudinary with correct resource type
// folder: target folder in Cloudinary (e.g., "documents/produk_hukum")
// filename: filename to delete (e.g., "abc123.pdf")
func (s *Service) DeleteFile(ctx context.Context, folder string, filename string) error {
	// Remove extension from filename
	ext := filepath.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(ext)]

	// Construct public ID
	publicID := fmt.Sprintf("uploads/%s/%s", folder, filenameWithoutExt)

	// Determine resource type from extension
	resourceType := GetResourceTypeFromExt(filename)

	// Delete from Cloudinary with correct resource type
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: string(resourceType),
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

// GetFileURL returns full Cloudinary URL for any file type with correct resource type
// folder: target folder (e.g., "documents/produk_hukum")
// filename: filename with extension (e.g., "abc123.pdf", "abc123.mp3")
func (s *Service) GetFileURL(folder string, filename string) string {
	if filename == "" {
		return ""
	}

	ext := filepath.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(ext)]
	publicID := fmt.Sprintf("uploads/%s/%s", folder, filenameWithoutExt)

	// Determine resource type from extension
	resourceType := GetResourceTypeFromExt(filename)

	// Build URL based on resource type
	var url string
	switch resourceType {
	case ResourceTypeImage:
		asset, err := s.cld.Image(publicID)
		if err != nil {
			return ""
		}
		url, _ = asset.String()
	case ResourceTypeVideo:
		asset, err := s.cld.Video(publicID)
		if err != nil {
			return ""
		}
		url, _ = asset.String()
	case ResourceTypeRaw:
		// For raw files, construct URL manually with extension
		cloudName := s.cld.Config.Cloud.CloudName
		url = fmt.Sprintf("https://res.cloudinary.com/%s/raw/upload/v1/%s%s", cloudName, publicID, ext)
	}

	return url
}

// GetDownloadURL returns Cloudinary URL with fl_attachment for forced download
// folder: target folder (e.g., "documents/produk_hukum")
// filename: filename with extension (e.g., "abc123.pdf")
func (s *Service) GetDownloadURL(folder string, filename string) string {
	if filename == "" {
		return ""
	}

	ext := filepath.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(ext)]
	publicID := fmt.Sprintf("uploads/%s/%s", folder, filenameWithoutExt)
	resourceType := GetResourceTypeFromExt(filename)
	cloudName := s.cld.Config.Cloud.CloudName

	// Build download URL with fl_attachment transformation
	return fmt.Sprintf("https://res.cloudinary.com/%s/%s/upload/fl_attachment/v1/%s%s",
		cloudName, resourceType, publicID, ext)
}
