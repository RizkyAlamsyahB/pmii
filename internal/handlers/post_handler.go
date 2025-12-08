package handlers

import (
	"net/http"
	"time"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- DEFINISI MODEL (Karena tidak ada postModel.go) ---
// Jika nanti ingin dirapikan, pindahkan struct ini ke folder internal/domain/post.go
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Image     string         `json:"image"` // URL dari Cloudinary
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- HANDLER FUNCTIONS ---

// 1. CREATE POST (Buat Berita Baru dengan Gambar)
func CreatePost(c *gin.Context) {
	// A. Validasi Input Text
	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul dan Konten wajib diisi"})
		return
	}

	// B. Proses Upload Gambar
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar wajib diupload"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file gambar"})
		return
	}
	defer file.Close()

	filename := "post-" + time.Now().Format("20060102-150405")
	imageUrl, err := utils.UploadToCloudinary(file, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload ke Cloudinary: " + err.Error()})
		return
	}

	// C. Simpan ke Database
	post := Post{
		Title:   title,
		Content: content,
		Image:   imageUrl,
	}

	// Jika error "DB not found", pastikan package config sudah di-import dan variabel DB public
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berita berhasil dibuat",
		"data":    post,
	})
}

// 2. GET ALL POSTS (List Berita)
func GetPosts(c *gin.Context) {
	var posts []Post

	// Mengambil semua data (bisa ditambahkan pagination nanti)
	if err := config.DB.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}

// 3. GET POST BY ID (Detail Berita)
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post Post

	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

// 4. UPDATE POST
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post Post

	// Cek apakah data ada
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// Update Text (jika ada input baru)
	if title := c.PostForm("title"); title != "" {
		post.Title = title
	}
	if content := c.PostForm("content"); content != "" {
		post.Content = content
	}

	// Update Gambar (Opsional, hanya jika user upload file baru)
	fileHeader, err := c.FormFile("image")
	if err == nil {
		// User upload gambar baru
		file, _ := fileHeader.Open()
		defer file.Close()

		filename := "post-update-" + time.Now().Format("20060102-150405")
		newImageUrl, errUpload := utils.UploadToCloudinary(file, filename)
		if errUpload != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar baru"})
			return
		}
		post.Image = newImageUrl
	}

	// Simpan perubahan
	config.DB.Save(&post)

	c.JSON(http.StatusOK, gin.H{
		"message": "Berita berhasil diupdate",
		"data":    post,
	})
}

// 5. DELETE POST
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post Post

	// Cek keberadaan data
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// Hapus (Soft Delete karena pakai gorm.DeletedAt)
	config.DB.Delete(&post)

	c.JSON(http.StatusOK, gin.H{
		"message": "Berita berhasil dihapus",
	})
}
