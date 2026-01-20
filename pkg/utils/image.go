package utils

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
)

// ImageDimension represents width and height of an image
type ImageDimension struct {
	Width  int
	Height int
}

// ParseResolution parses resolution string like "728x90" or "16x9" to width and height ratio
func ParseResolution(resolution string) (*ImageDimension, error) {
	parts := strings.Split(strings.ToLower(resolution), "x")
	if len(parts) != 2 {
		return nil, errors.New("format resolusi tidak valid, gunakan format WxH (contoh: 728x90)")
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.New("lebar resolusi tidak valid")
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errors.New("tinggi resolusi tidak valid")
	}

	return &ImageDimension{Width: width, Height: height}, nil
}

// GetImageDimension gets the dimension of an uploaded image file
func GetImageDimension(file *multipart.FileHeader) (*ImageDimension, error) {
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("gagal membuka file gambar")
	}
	defer src.Close()

	img, _, err := image.DecodeConfig(src)
	if err != nil {
		return nil, errors.New("gagal membaca dimensi gambar, pastikan file adalah gambar yang valid (jpg, png, gif)")
	}

	return &ImageDimension{Width: img.Width, Height: img.Height}, nil
}

// ValidateImageAspectRatio validates if the uploaded image matches the required aspect ratio
// tolerance is the allowed deviation percentage (e.g., 0.1 = 10%)
func ValidateImageAspectRatio(file *multipart.FileHeader, requiredResolution string, tolerance float64) error {
	// Parse required resolution
	required, err := ParseResolution(requiredResolution)
	if err != nil {
		return err
	}

	// Get actual image dimension
	actual, err := GetImageDimension(file)
	if err != nil {
		return err
	}

	// Calculate aspect ratios
	requiredRatio := float64(required.Width) / float64(required.Height)
	actualRatio := float64(actual.Width) / float64(actual.Height)

	// Check if the aspect ratios match within tolerance
	diff := math.Abs(requiredRatio-actualRatio) / requiredRatio
	if diff > tolerance {
		return errors.New("aspek rasio gambar tidak sesuai. Diperlukan rasio " + requiredResolution + " (rasio " + formatRatio(requiredRatio) + "), gambar Anda memiliki rasio " + formatRatio(actualRatio))
	}

	return nil
}

// formatRatio formats a ratio to a readable string
func formatRatio(ratio float64) string {
	return strconv.FormatFloat(ratio, 'f', 2, 64) + ":1"
}

// ValidateImageDimension validates if the uploaded image matches the exact required dimensions
func ValidateImageDimension(file *multipart.FileHeader, requiredResolution string) error {
	// Parse required resolution
	required, err := ParseResolution(requiredResolution)
	if err != nil {
		return err
	}

	// Get actual image dimension
	actual, err := GetImageDimension(file)
	if err != nil {
		return err
	}

	// Check exact match
	if actual.Width != required.Width || actual.Height != required.Height {
		return errors.New("dimensi gambar tidak sesuai. Diperlukan " + requiredResolution + ", gambar Anda " + strconv.Itoa(actual.Width) + "x" + strconv.Itoa(actual.Height))
	}

	return nil
}
