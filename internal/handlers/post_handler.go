package handlers

import (
	"math"
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
	// Ambil Query Param
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var posts []domain.Post
	var total int64

	// Preload Relasi
	query := config.DB.Model(&domain.Post{}).Preload("Tags").Preload("Category")

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchKeyword, searchKeyword)
	}

	// Hitung Total Data
	query.Count(&total)

	// Ambil Data
	result := query.Limit(limit).Offset(offset).Order("published_at DESC").Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data berita"))
		return
	}

	data := responses.FromDomainListToPostResponse(posts)

	// --- PERBAIKAN UTAMA DISINI ---
	// 1. Hitung Last Page secara manual (karena base_response membutuhkannya)
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage < 1 {
		lastPage = 1
	}

	// 2. Panggil fungsi dengan urutan argumen yang BENAR sesuai base_response.go:
	// func SuccessResponseWithPagination(code int, message string, data interface{}, page, limit int, total int64, lastPage int)
	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(200, "List of posts", data, page, limit, total, lastPage))
}

// 2. CREATE POST
func CreatePost(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryID, _ := strconv.Atoi(c.PostForm("category_id"))

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Judul dan Konten wajib diisi"))
		return
	}

	imageName := ""
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, _ := fileHeader.Open()
		defer file.Close()
		filename := "post-" + time.Now().Format("20060102-150405")

		// Pastikan utils mengembalikan URL string
		url, errUpload := utils.UploadToCloudinary(file, filename)
		if errUpload != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal upload gambar"))
			return
		}
		imageName = url // Simpan URL lengkap
	}

	var excerptPtr *string = nil
	var featuredImagePtr *string = nil

	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))

	excerpt := content
	if len(content) > 150 {
		excerpt = content[:150] + "..."
	}

	if len(excerpt) > 0 {
		excerptPtr = &excerpt
	}

	if imageName != "" {
		featuredImagePtr = &imageName
	}

	publishedTime := time.Now()
	publishedAtPtr := &publishedTime

	// Logic Tags
	var tagEntities []domain.Tag
	tagsInput := c.PostForm("tags")
	if tagsInput != "" {
		tagNameList := strings.Split(tagsInput, ",")
		for _, tagName := range tagNameList {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}

			tagSlug := strings.ToLower(strings.ReplaceAll(tagName, " ", "-"))

			var tag domain.Tag
			err := config.DB.Where(domain.Tag{Slug: tagSlug}).Attrs(domain.Tag{Name: tagName}).FirstOrCreate(&tag).Error
			if err == nil {
				tagEntities = append(tagEntities, tag)
			}
		}
	}

	post := domain.Post{
		Title:         title,
		Slug:          slug,
		Content:       content,
		Excerpt:       excerptPtr,
		FeaturedImage: featuredImagePtr,
		CategoryID:    categoryID,
		UserID:        1,
		PublishedAt:   publishedAtPtr,
		Tags:          tagEntities,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menyimpan berita"))
		return
	}

	// Reload data agar ID tags & category muncul di response
	config.DB.Preload("Category").Preload("Tags").First(&post, post.ID)

	responseDTO := responses.FromDomainToPostResponse(post)

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Berita berhasil dibuat", responseDTO))
}

// 3. GET DETAIL POST
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post domain.Post

	query := config.DB.Preload("Category").Preload("Tags")
	if _, err := strconv.Atoi(id); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("slug = ?", id)
	}

	if err := query.First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}

	responseDTO := responses.FromDomainToPostResponse(post)

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Detail berita ditemukan", responseDTO))
}

// 4. DELETE POST
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post domain.Post

	if err := config.DB.Unscoped().Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}

	config.DB.Delete(&post)

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berita berhasil dihapus", nil))
}

// 5. UPDATE POST
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post domain.Post

	// Preload Tags eksisting agar bisa di-replace
	if err := config.DB.Preload("Tags").Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}

	if val := c.PostForm("title"); val != "" {
		post.Title = val
		post.Slug = strings.ToLower(strings.ReplaceAll(val, " ", "-"))
	}

	if val := c.PostForm("content"); val != "" {
		post.Content = val
		excerpt := val
		if len(val) > 150 {
			excerpt = val[:150] + "..."
		}
		if len(excerpt) > 0 {
			post.Excerpt = &excerpt
		}
	}

	if val := c.PostForm("category_id"); val != "" {
		if catID, err := strconv.Atoi(val); err == nil {
			post.CategoryID = catID
		}
	}

	// Logic Update Tags
	tagString := c.PostForm("tags")
	// Jika key "tags" dikirim (meskipun string kosong), kita proses
	if c.Request.PostForm.Has("tags") {
		var tagEntities []domain.Tag
		if tagString != "" {
			tagNameList := strings.Split(tagString, ",")
			for _, tagName := range tagNameList {
				tagName = strings.TrimSpace(tagName)
				if tagName == "" {
					continue
				}

				tagSlug := strings.ToLower(strings.ReplaceAll(tagName, " ", "-"))
				var tag domain.Tag
				err := config.DB.Where(domain.Tag{Slug: tagSlug}).Attrs(domain.Tag{Name: tagName}).FirstOrCreate(&tag).Error
				if err == nil {
					tagEntities = append(tagEntities, tag)
				}
			}
		}
		// Replace tags lama dengan yang baru (atau hapus semua jika kosong)
		config.DB.Model(&post).Association("Tags").Replace(tagEntities)
	}

	// Update Gambar
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, _ := fileHeader.Open()
		defer file.Close()
		filename := "post-update-" + time.Now().Format("20060102-150405")

		url, _ := utils.UploadToCloudinary(file, filename)
		if url != "" {
			post.FeaturedImage = &url
		}
	}

	config.DB.Save(&post)

	config.DB.Preload("Category").Preload("Tags").First(&post, post.ID)
	responseDTO := responses.FromDomainToPostResponse(post)
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berita berhasil diupdate", responseDTO))
}
