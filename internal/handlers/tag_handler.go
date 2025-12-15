package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
)

// 1. GET ALL TAGS
func GetTags(c *gin.Context) {
	var tags []domain.Tag

	// Ambil semua tag (bisa diurutkan by Name ASC agar rapi)
	if err := config.DB.Order("name ASC").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data tags"))
		return
	}

	data := responses.FromDomainListToTagResponse(tags)
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "List of tags", data))
}

// 2. CREATE TAG
func CreateTag(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Nama tag wajib diisi"))
		return
	}

	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	// Cek duplikasi slug agar tidak error database unique constraint
	var existingTag domain.Tag
	if err := config.DB.Where("slug = ?", slug).First(&existingTag).Error; err == nil {
		c.JSON(http.StatusConflict, responses.ErrorResponse(409, "Tag dengan nama tersebut sudah ada"))
		return
	}

	tag := domain.Tag{
		Name: name,
		Slug: slug,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menyimpan tag"))
		return
	}

	responseDTO := responses.FromDomainToTagResponse(tag)
	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Tag berhasil dibuat", responseDTO))
}

// 3. DELETE TAG
func DeleteTag(c *gin.Context) {
	id := c.Param("id")
	var tag domain.Tag

	if err := config.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Tag tidak ditemukan"))
		return
	}

	// Hapus asosiasi di tabel pivot post_tags terlebih dahulu (Optional, GORM biasanya handle ini jika setup benar)
	// Tapi untuk aman, kita delete tag-nya saja, GORM akan mengurus foreign key jika CASCADE diset di DB.
	if err := config.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menghapus tag"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Tag berhasil dihapus", nil))
}

// 4. UPDATE TAG (Opsional, jarang dipakai tapi baik untuk kelengkapan)
func UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag domain.Tag

	if err := config.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Tag tidak ditemukan"))
		return
	}

	newName := c.PostForm("name")
	if newName != "" {
		tag.Name = newName
		tag.Slug = strings.ToLower(strings.ReplaceAll(newName, " ", "-"))

		// Simpan perubahan
		if err := config.DB.Save(&tag).Error; err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengupdate tag"))
			return
		}
	}

	responseDTO := responses.FromDomainToTagResponse(tag)
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Tag berhasil diupdate", responseDTO))
}
