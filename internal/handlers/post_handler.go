package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// 1. GET ALL POSTS (Dengan Pagination & Search)
func GetPosts(c *gin.Context) {
	// Ambil Query Parameter
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	// Hitung Offset
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var posts []domain.Post
	var total int64

	// Build Query
	query := config.DB.Model(&domain.Post{})

	// Fitur Search (Optional)
	if search != "" {
		// Mencari di Title atau Content
		searchKeyword := "%" + search + "%"
		query = query.Where("post_title ILIKE ? OR post_contents ILIKE ?", searchKeyword, searchKeyword)
	}

	// Hitung Total Data (Untuk Pagination)
	query.Count(&total)

	// Ambil Data dengan Limit & Offset
	result := query.Limit(limit).Offset(offset).Order("post_date DESC").Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data berita"))
		return
	}

	// KONVERSI KE DTO (Agar JSON outputnya camelCase: imageUrl, title, dll)
	data := responses.FromDomainListToPostResponse(posts)

	// Return menggunakan Helper Pagination Baru
	c.JSON(http.StatusOK, responses.SuccessPaginationResponse(200, "List of posts", page, limit, total, data))
}

// 2. CREATE POST
func CreatePost(c *gin.Context) {
	// Tangkap Input
	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryID, _ := strconv.Atoi(c.PostForm("category_id"))

	// Validasi Dasar
	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Judul dan Konten wajib diisi"))
		return
	}

	// Upload Gambar
	imageUrl := ""
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, _ := fileHeader.Open()
		defer file.Close()
		filename := "post-" + time.Now().Format("20060102-150405")

		// Upload ke Cloudinary (Pakai Utils)
		uploadedUrl, errUpload := utils.UploadToCloudinary(file, filename)
		if errUpload != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal upload gambar"))
			return
		}
		imageUrl = uploadedUrl
	}

	// Buat Slug Sederhana
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))

	// Buat Excerpt
	excerpt := content
	if len(content) > 150 {
		excerpt = content[:150] + "..."
	}

	// LOGIC TAGS (Complex Part)
	var tagEntities []domain.Tag
	tagsInput := c.PostForm("tags") // input: "news,politik"
	if tagsInput != "" {
		tagNameList := strings.Split(tagsInput, ",")
		for _, tagName := range tagNameList {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}

			tagSlug := strings.ToLower(strings.ReplaceAll(tagName, " ", "-"))

			// Cari Tag di DB, kalau tidak ada -> Buat Baru. Kalau ada -> Pakai yang lama.
			var tag domain.Tag
			err := config.DB.Where(domain.Tag{Slug: tagSlug}).Attrs(domain.Tag{Name: tagName}).FirstOrCreate(&tag).Error
			if err == nil {
				tagEntities = append(tagEntities, tag)
			}
		}
	}

	// Isi Model Domain
	post := domain.Post{
		Title:         title,
		Slug:          slug,
		Content:       content,
		Excerpt:       excerpt,
		FeaturedImage: imageUrl, // Kolom baru
		CategoryID:    categoryID,
		UserID:        1, // Hardcode ID user dulu
		Status:        1,
		Views:         0,
		PublishedAt:   time.Now(),
		Tags:          tagEntities, // GORM akan otomatis isi tabel pivot 'post_tags'
	}

	// Simpan ke DB
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menyimpan berita"))
		return
	}

	// Konversi ke DTO Response
	responseDTO := responses.FromDomainToPostResponse(post)

	// Return Standard Response
	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Berita berhasil dibuat", responseDTO))
}

// 3. GET DETAIL POST
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post domain.Post

	// Cek apakah ID berupa angka (Search by ID) atau String (Search by Slug)
	query := config.DB
	if _, err := strconv.Atoi(id); err == nil {
		// Search by ID
		query = query.Where("post_id = ?", id)
	} else {
		// Search by Slug
		query = query.Where("post_slug = ?", id)
	}

	if err := query.First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}

	// Increment Views
	config.DB.Model(&post).UpdateColumn("post_views", post.Views+1)

	// Konversi ke DTO Response
	responseDTO := responses.FromDomainToPostResponse(post)

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Detail berita ditemukan", responseDTO))
}

// 4. DELETE POST
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post domain.Post

	// Cek Data
	if err := config.DB.Where("post_id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}

	// Delete
	config.DB.Delete(&post)

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berita berhasil dihapus", nil))
}

// 5. UPDATE POST
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post domain.Post

	// 1. Cek Data Eksisting (Preload Tags agar tag lama terbaca)
	if err := config.DB.Preload("Tags").Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}

	// 2. Update Text Fields jika ada input
	if val := c.PostForm("title"); val != "" {
		post.Title = val
		post.Slug = strings.ToLower(strings.ReplaceAll(val, " ", "-"))
	}
	if val := c.PostForm("content"); val != "" {
		post.Content = val
		if len(val) > 150 {
			post.Excerpt = val[:150] + "..."
		} else {
			post.Excerpt = val
		}
	}

	// Update Category ID
	if val := c.PostForm("category_id"); val != "" {
		if catID, err := strconv.Atoi(val); err == nil {
			post.CategoryID = catID
		}
	}

	// 3. LOGIC UPDATE TAGS (PENTING: Ini perbaikan error "val")
	tagString := c.PostForm("tags")
	if tagString != "" {
		var tagEntities []domain.Tag
		tagNameList := strings.Split(tagString, ",")

		for _, tagName := range tagNameList {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}

			tagSlug := strings.ToLower(strings.ReplaceAll(tagName, " ", "-"))

			// Cari atau Buat Tag baru
			var tag domain.Tag
			err := config.DB.Where(domain.Tag{Slug: tagSlug}).Attrs(domain.Tag{Name: tagName}).FirstOrCreate(&tag).Error
			if err == nil {
				tagEntities = append(tagEntities, tag)
			}
		}

		config.DB.Model(&post).Association("Tags").Replace(tagEntities)
	}

	// 4. Cek Image Baru (Update jika ada file)
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, _ := fileHeader.Open()
		defer file.Close()
		filename := "post-update-" + time.Now().Format("20060102-150405")
		newUrl, _ := utils.UploadToCloudinary(file, filename)
		post.FeaturedImage = newUrl
	}

	// 5. Simpan Perubahan ke Tabel Posts
	config.DB.Save(&post)

	config.DB.Preload("Category").Preload("Tags").First(&post, post.ID)
	responseDTO := responses.FromDomainToPostResponse(post)
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berita berhasil diupdate", responseDTO))
}
