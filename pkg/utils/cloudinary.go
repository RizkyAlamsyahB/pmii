package utils

import (
	"context"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file multipart.File, filename string) (string, error) {
	// 1. Ambil config dari Env
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")

	// 2. Inisialisasi Cloudinary
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", err
	}

	// 3. Setup Context (Time out 10 detik agar tidak hanging jika internet lemot)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Proses Upload
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   folder,   // Folder tujuan di Cloudinary
		PublicID: filename, // Nama file (opsional, bisa diganti auto)
	})

	if err != nil {
		return "", err
	}

	// 5. Kembalikan URL Secure (https)
	return uploadResult.SecureURL, nil
}
