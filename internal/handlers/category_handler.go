package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
)

// 1. GET ALL CATEGORIES
func GetCategories(c *gin.Context) {
	var categories []domain.Category

	// Ambil semua data, urutkan dari yang terbaru
	if err := config.DB.Order("created_at DESC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data kategori"))
		return
	}

	data := responses.FromDomainListToCategoryResponse(categories)
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "List of categories", data))
}

// 2. CREATE CATEGORY
func CreateCategory(c *gin.Context) {
	// Menerima input x-www-form-urlencoded
	name := c.PostForm("name")

	// Validasi sederhana
	if name == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Nama kategori wajib diisi"))
		return
	}

	// Generate Slug
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	var descriptionPtr *string // Default nil
	descInput := c.PostForm("description")

	// Jika input tidak kosong, ambil address-nya
	if descInput != "" {
		descriptionPtr = &descInput
	}

	category := domain.Category{
		Name:        name,
		Slug:        slug,
		Description: descriptionPtr,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menyimpan kategori"))
		return
	}

	responseDTO := responses.FromDomainToCategoryResponse(category)
	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Kategori berhasil dibuat", responseDTO))
}

// 3. DELETE CATEGORY
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category domain.Category

	// Cek apakah data ada
	if err := config.DB.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Kategori tidak ditemukan"))
		return
	}

	// Hard Delete (atau Soft Delete tergantung model Anda)
	if err := config.DB.Delete(&category).Error; err != nil {
		// Handle error Foreign Key constraint (jika kategori dipakai di berita)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menghapus kategori. Pastikan tidak ada berita yang menggunakan kategori ini."))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Kategori berhasil dihapus", nil))
}

// 4. UPDATE CATEGORY (Opsional - Bonus)
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category domain.Category

	if err := config.DB.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Kategori tidak ditemukan"))
		return
	}

	newName := c.PostForm("name")
	if newName != "" {
		category.Name = newName
		category.Slug = strings.ToLower(strings.ReplaceAll(newName, " ", "-"))
	}

	if desc := c.PostForm("description"); desc != "" {
		category.Description = &desc
	}

	config.DB.Save(&category)

	responseDTO := responses.FromDomainToCategoryResponse(category)
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Kategori berhasil diupdate", responseDTO))
}
