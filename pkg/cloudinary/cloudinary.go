package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
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
	ResourceTypeRaw   ResourceType = "raw"   // ZIP, TXT, etc.
)

// GetResourceTypeFromExt determines the correct Cloudinary resource type from file extension
// Note: PDF uses "image" resource type for preview support in browser
func GetResourceTypeFromExt(filename string) ResourceType {
	ext := strings.ToLower(filepath.Ext(filename))

	// Image extensions (including PDF for preview support)
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".bmp": true, ".ico": true, ".svg": true,
		".pdf": true, // PDF uses image resource type for browser preview
	}

	// Video & Audio extensions (Cloudinary uses "video" for both)
	// Note: WMA is included as video - will be transformed to MP3 on delivery
	videoExts := map[string]bool{
		".mp4": true, ".webm": true, ".mov": true, ".avi": true, ".mkv": true,
		".mp3": true, ".wav": true, ".ogg": true, ".flac": true, ".aac": true,
		".m4a": true, ".opus": true, ".wma": true, // Audio formats (WMA will be converted to MP3)
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
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Upload to Cloudinary - public ID is just timestamp, folder for organization only
	_, err = s.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		AssetFolder:    fmt.Sprintf("uploads/%s", folder),
		PublicID:       timestamp,
		UseFilename:    boolPtr(false),
		UniqueFilename: boolPtr(false),
		ResourceType:   "image",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return filename only (e.g., "1234567890.jpg")
	filename := fmt.Sprintf("%s%s", timestamp, ext)
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
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Determine resource type from file extension
	resourceType := GetResourceTypeFromExt(file.Filename)

	// For raw files, include extension in public ID
	var publicID string
	if resourceType == ResourceTypeRaw {
		publicID = fmt.Sprintf("%s%s", timestamp, ext) // e.g., "1234567890.zip"
	} else {
		publicID = timestamp // e.g., "1234567890"
	}

	// Upload to Cloudinary - public ID is just timestamp, folder for organization only
	_, err = s.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		AssetFolder:    fmt.Sprintf("uploads/%s", folder),
		PublicID:       publicID,
		UseFilename:    boolPtr(false),
		UniqueFilename: boolPtr(false),
		ResourceType:   string(resourceType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return filename with extension (e.g., "1234567890.pdf")
	filename := fmt.Sprintf("%s%s", timestamp, ext)
	return filename, nil
}

// UploadFromPath uploads file from local path to Cloudinary
// folder: target folder in Cloudinary (e.g., "members", "documents/produk_hukum")
// filePath: local file path to upload
// Returns: filename with extension (e.g., "1234567890.jpg"), error
func (s *Service) UploadFromPath(ctx context.Context, folder string, filePath string) (string, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get filename and extension
	originalFilename := filepath.Base(filePath)
	ext := filepath.Ext(originalFilename)
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1000000) // milliseconds for uniqueness

	// Determine resource type from file extension
	resourceType := GetResourceTypeFromExt(originalFilename)

	// For raw files, include extension in public ID
	var publicID string
	if resourceType == ResourceTypeRaw {
		publicID = fmt.Sprintf("%s%s", timestamp, ext)
	} else {
		publicID = timestamp
	}

	// Upload to Cloudinary - public ID is just timestamp, folder for organization only
	_, err = s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		AssetFolder:    fmt.Sprintf("uploads/%s", folder),
		PublicID:       publicID,
		UseFilename:    boolPtr(false),
		UniqueFilename: boolPtr(false),
		ResourceType:   string(resourceType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return filename with extension
	filename := fmt.Sprintf("%s%s", timestamp, ext)
	return filename, nil
}

// UploadFromPathWithOverwrite uploads file from local path to Cloudinary with fixed public_id
// This is used by seeders to prevent duplicate uploads when database is reset
// If file with same public_id exists, it will be overwritten (not duplicated)
// folder: target folder in Cloudinary (e.g., "seeds/members")
// filePath: local file path to upload
// customPublicID: fixed identifier for the file (e.g., original filename without ext)
// Returns: filename with extension, error
func (s *Service) UploadFromPathWithOverwrite(ctx context.Context, folder string, filePath string, customPublicID string) (string, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get filename and extension
	originalFilename := filepath.Base(filePath)
	ext := filepath.Ext(originalFilename)

	// Determine resource type from file extension
	resourceType := GetResourceTypeFromExt(originalFilename)

	// For raw files, include extension in public ID
	var publicID string
	if resourceType == ResourceTypeRaw {
		publicID = customPublicID + ext
	} else {
		publicID = customPublicID
	}

	// Upload to Cloudinary - public ID is just customPublicID, folder for organization only
	_, err = s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		AssetFolder:    fmt.Sprintf("uploads/%s", folder),
		PublicID:       publicID,
		UseFilename:    boolPtr(false),
		UniqueFilename: boolPtr(false),
		ResourceType:   string(resourceType),
		Overwrite:      boolPtr(true), // Overwrite existing file with same public_id
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// Return filename with extension
	filename := fmt.Sprintf("%s%s", customPublicID, ext)
	return filename, nil
}

// boolPtr returns a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}

// DeleteImage deletes image from Cloudinary
// folder: not used anymore, kept for backward compatibility
// filename: filename to delete (e.g., "abc123.jpg")
func (s *Service) DeleteImage(ctx context.Context, folder string, filename string) error {
	// Remove extension from filename to get public ID
	ext := filepath.Ext(filename)
	publicID := filename[:len(filename)-len(ext)]

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
// folder: not used anymore, kept for backward compatibility
// filename: filename to delete (e.g., "abc123.pdf")
func (s *Service) DeleteFile(ctx context.Context, folder string, filename string) error {
	// Determine resource type from extension
	resourceType := GetResourceTypeFromExt(filename)

	// Construct public ID based on resource type
	var publicID string
	if resourceType == ResourceTypeRaw {
		// For raw files, public ID includes extension
		publicID = filename
	} else {
		// For image/video, public ID excludes extension
		ext := filepath.Ext(filename)
		publicID = filename[:len(filename)-len(ext)]
	}

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
// folder: not used anymore, kept for backward compatibility
// filename: filename (e.g., "abc123.jpg")
func (s *Service) GetImageURL(folder string, filename string) string {
	if filename == "" {
		return ""
	}

	ext := filepath.Ext(filename)
	publicID := filename[:len(filename)-len(ext)]

	asset, err := s.cld.Image(publicID)
	if err != nil {
		return ""
	}
	url, _ := asset.String()
	return url
}

// GetFileURL returns full Cloudinary URL for any file type with correct resource type
// folder: not used anymore, kept for backward compatibility
// filename: filename with extension (e.g., "abc123.pdf", "abc123.mp3")
func (s *Service) GetFileURL(folder string, filename string) string {
	if filename == "" {
		return ""
	}

	ext := strings.ToLower(filepath.Ext(filename))
	filenameWithoutExt := filename[:len(filename)-len(ext)]
	cloudName := s.cld.Config.Cloud.CloudName

	// Determine resource type from extension
	resourceType := GetResourceTypeFromExt(filename)

	// Build URL based on resource type
	var url string
	switch resourceType {
	case ResourceTypeImage:
		// For PDF, construct URL manually to include .pdf extension
		if ext == ".pdf" {
			url = fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/v1/%s%s", cloudName, filenameWithoutExt, ext)
		} else {
			asset, err := s.cld.Image(filenameWithoutExt)
			if err != nil {
				return ""
			}
			url, _ = asset.String()
		}
	case ResourceTypeVideo:
		asset, err := s.cld.Video(filenameWithoutExt)
		if err != nil {
			return ""
		}
		url, _ = asset.String()
	case ResourceTypeRaw:
		// For raw files, public ID includes extension (e.g., "1234567890.zip")
		url = fmt.Sprintf("https://res.cloudinary.com/%s/raw/upload/v1/%s", cloudName, filename)
	}

	return url
}

// GetDownloadURL returns Cloudinary URL with fl_attachment for forced download
// folder: not used anymore, kept for backward compatibility
// filename: filename with extension (e.g., "abc123.pdf")
func (s *Service) GetDownloadURL(folder string, filename string) string {
	if filename == "" {
		return ""
	}

	ext := strings.ToLower(filepath.Ext(filename))
	filenameWithoutExt := filename[:len(filename)-len(ext)]
	resourceType := GetResourceTypeFromExt(filename)
	cloudName := s.cld.Config.Cloud.CloudName

	// Special handling for WMA files - always use video endpoint and convert to mp3
	if ext == ".wma" {
		// Use f_mp3 to convert WMA to MP3 on-the-fly, fl_attachment for download
		return fmt.Sprintf("https://res.cloudinary.com/%s/video/upload/f_mp3,fl_attachment/v1/%s",
			cloudName, filenameWithoutExt)
	}

	// Build download URL with fl_attachment transformation
	if resourceType == ResourceTypeRaw {
		// For raw files, public ID includes extension
		// Note: fl_attachment is NOT supported for raw files, use direct URL
		return fmt.Sprintf("https://res.cloudinary.com/%s/raw/upload/v1/%s", cloudName, filename)
	}

	// For image/video, public ID excludes extension
	return fmt.Sprintf("https://res.cloudinary.com/%s/%s/upload/fl_attachment/v1/%s%s",
		cloudName, resourceType, filenameWithoutExt, ext)
}
