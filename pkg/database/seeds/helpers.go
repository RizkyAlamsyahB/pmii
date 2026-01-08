package seeds

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/garuda-labs-1/pmii-be/pkg/cloudinary"
)

// uploadFile uploads file to Cloudinary from local path with overwrite support
// Uses the original filename (without ext) as public_id to prevent duplicates
// If database is reset and seeder runs again, Cloudinary files will be overwritten, not duplicated
func uploadFile(uploader *cloudinary.Service, filePath, folder string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filePath)
	}

	// Extract original filename without extension as public_id
	// e.g., "1766022783.pdf" -> "1766022783"
	originalFilename := filepath.Base(filePath)
	publicID := strings.TrimSuffix(originalFilename, filepath.Ext(originalFilename))

	// Upload with overwrite enabled using fixed public_id
	url, err := uploader.UploadFromPathWithOverwrite(context.Background(), folder, filePath, publicID)
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return url, nil
}

// getFilePath constructs full file path from seeds directory
func getFilePath(seedsPath string, parts ...string) string {
	allParts := append([]string{seedsPath}, parts...)
	return filepath.Join(allParts...)
}

// logSeederStart logs the start of a seeder
func logSeederStart(name string) {
	log.Printf("📝 Seeding %s...", name)
}

// logSeederSkip logs when seeder is skipped
func logSeederSkip(name string) {
	log.Printf("ℹ️  %s already exists, skipping...", name)
}

// logSeederProgress logs seeding progress
func logSeederProgress(current, total int, item string) {
	log.Printf("   [%d/%d] Created: %s", current, total, item)
}

// logSeederResult logs final seeding result
func logSeederResult(name string, success, total int) {
	log.Printf("✅ %s seeded: %d/%d successfully", name, success, total)
}
